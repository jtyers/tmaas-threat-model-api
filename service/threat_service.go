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
	ErrTooManyThreats = errors.New("too many threats")
	ErrNoSuchThreat   = errors.New("no such threat")
	ErrNoDataToUpdate = errors.New("no data to update")
)

type ThreatService interface {
	// Retrieve a threat by threatID.
	GetThreat(ctx context.Context, id m.ThreatID) (*m.Threat, error)

	// Retrieve all threat.
	GetThreats(ctx context.Context) ([]*m.Threat, error)

	// Creates a threat. `threat` should not have ID or ThreatID set.
	CreateThreat(ctx context.Context, threat m.Threat) (*m.Threat, error)

	// Updates a threat
	UpdateThreat(ctx context.Context, threatID m.ThreatID, threat m.Threat) error
}

type DefaultThreatService struct {
	dao dao.ThreatDao

	randomIDProvider id.RandomIDProvider
	validator        validator.StructValidator
}

func NewDefaultThreatService(dao dao.ThreatDao, randomIDProvider id.RandomIDProvider, validator validator.StructValidator) *DefaultThreatService {
	return &DefaultThreatService{dao, randomIDProvider, validator}
}

func (g *DefaultThreatService) GetThreat(ctx context.Context, threatID m.ThreatID) (*m.Threat, error) {
	threat, err := g.dao.Get(ctx, threatID.String())

	if err != nil {
		if err == servicedao.ErrNoSuchDocument {
			return nil, ErrNoSuchThreat
		} else {
			return nil, fmt.Errorf("error retrieving threat: %v", err)
		}
	}

	return threat, nil
}

// CreateThreat Creates a new Threat in Firestore.
//
// The threat supplied should not have its ID or ThreatID fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created threat is returned to the caller, with ID and ThreatID set.
func (g *DefaultThreatService) CreateThreat(ctx context.Context, threat m.Threat) (*m.Threat, error) {
	if threat.ThreatID != "" {
		return nil, fmt.Errorf("cannot create a threat that already has ThreatID set")
	}

	err := g.validator.ValidateForCreate(threat)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(threat)
	if err != nil {
		return nil, err
	}

	threat.ThreatID = m.ThreatID(ThreatIDPrefix + g.randomIDProvider.GenerateID())

	result, err := g.dao.Create(ctx, &threat)
	if err != nil {
		return nil, fmt.Errorf("error creating threat: %v", err)
	}

	return result, nil
}

func (g *DefaultThreatService) UpdateThreat(ctx context.Context, threatID m.ThreatID, threat m.Threat) error {
	err := g.validator.ValidateForUpdate(threat)
	if err != nil {
		return err
	}

	if threat.ThreatID != threatID {
		return fmt.Errorf("given threat IDs do not match")
	}

	queryThreat := m.Threat{ThreatID: threatID}

	_, err = g.dao.UpdateWhereExactSingle(ctx, &queryThreat, &threat)
	if err != nil {
		return fmt.Errorf("error updating threat: %v", err)
	}

	return nil
}

func (g *DefaultThreatService) GetThreats(ctx context.Context) ([]*m.Threat, error) {
	threats, err := g.dao.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("error retrieving threats: %v", err)
	}

	return threats, nil
}
