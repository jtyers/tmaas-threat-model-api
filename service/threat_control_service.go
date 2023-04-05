package service

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mocks.go -package $GOPACKAGE

import (
	"context"
	"errors"
	"fmt"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/validator"
	servicedao "github.com/jtyers/tmaas-service-dao"
	"github.com/jtyers/tmaas-service-util/id"
	dao "github.com/jtyers/tmaas-threat-service/dao"
)

var (
	ErrTooManyThreatControls = errors.New("too many threatControls")
	ErrNoSuchThreatControl   = errors.New("no such threatControl")
)

type ThreatControlService interface {
	// Retrieve a threatControl by threatControlID.
	GetThreatControl(ctx context.Context, id m.ThreatControlID) (*m.ThreatControl, error)

	// Retrieve all threatControl.
	GetThreatControls(ctx context.Context) ([]*m.ThreatControl, error)

	// Creates a threatControl. `threatControl` should not have ID or ThreatControlID set.
	CreateThreatControl(ctx context.Context, threatControl m.ThreatControl) (*m.ThreatControl, error)

	// Updates a threatControl
	UpdateThreatControl(ctx context.Context, threatControlID m.ThreatControlID, threatControl m.ThreatControl) error
}

type DefaultThreatControlService struct {
	dao dao.ThreatControlDao

	randomIDProvider id.RandomIDProvider
	validator        validator.StructValidator
}

func NewDefaultThreatControlService(dao dao.ThreatControlDao, randomIDProvider id.RandomIDProvider, validator validator.StructValidator) *DefaultThreatControlService {
	return &DefaultThreatControlService{dao, randomIDProvider, validator}
}

func (g *DefaultThreatControlService) GetThreatControl(ctx context.Context, threatControlID m.ThreatControlID) (*m.ThreatControl, error) {
	threatControl, err := g.dao.Get(ctx, threatControlID.String())

	if err != nil {
		if err == servicedao.ErrNoSuchDocument {
			return nil, ErrNoSuchThreatControl
		} else {
			return nil, fmt.Errorf("error retrieving threatControl: %v", err)
		}
	}

	return threatControl, nil
}

// CreateThreatControl Creates a new ThreatControl in Firestore.
//
// The threatControl supplied should not have its ID or ThreatControlID fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created threatControl is returned to the caller, with ID and ThreatControlID set.
func (g *DefaultThreatControlService) CreateThreatControl(ctx context.Context, threatControl m.ThreatControl) (*m.ThreatControl, error) {
	if threatControl.ThreatControlID != "" {
		return nil, fmt.Errorf("cannot create a threatControl that already has ThreatControlID set")
	}

	err := g.validator.ValidateForCreate(threatControl)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(threatControl)
	if err != nil {
		return nil, err
	}

	threatControl.ThreatControlID = m.ThreatControlID(ThreatControlIDPrefix + g.randomIDProvider.GenerateID())

	result, err := g.dao.Create(ctx, &threatControl)
	if err != nil {
		return nil, fmt.Errorf("error creating threatControl: %v", err)
	}

	return result, nil
}

func (g *DefaultThreatControlService) UpdateThreatControl(ctx context.Context, threatControlID m.ThreatControlID, threatControl m.ThreatControl) error {
	err := g.validator.ValidateForUpdate(threatControl)
	if err != nil {
		return err
	}

	if threatControl.ThreatControlID != threatControlID {
		return fmt.Errorf("given threatControl IDs do not match")
	}

	queryThreatControl := m.ThreatControl{ThreatControlID: threatControlID}

	_, err = g.dao.UpdateWhereExactSingle(ctx, &queryThreatControl, &threatControl)
	if err != nil {
		return fmt.Errorf("error updating threatControl: %v", err)
	}

	return nil
}

func (g *DefaultThreatControlService) GetThreatControls(ctx context.Context) ([]*m.ThreatControl, error) {
	threatControls, err := g.dao.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("error retrieving threatControls: %v", err)
	}

	return threatControls, nil
}
