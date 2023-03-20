package service

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mocks.go -package $GOPACKAGE

import (
	"context"
	"errors"
	"fmt"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/validator"
	"github.com/jtyers/tmaas-service-util/id"
	"github.com/jtyers/tmaas-service-util/requestor"
	dao "github.com/jtyers/tmaas-threat-service/dao"
)

var (
	ErrTooManyThreats = errors.New("too many threats")
	ErrNoSuchThreat   = errors.New("no such threat")
	ErrNoDataToUpdate = errors.New("no data to update")
)

type ThreatService interface {
	// Retrieve a threat by threatId.
	GetThreat(ctx context.Context, id m.ThreatId) (*m.Threat, error)

	// Retrieve all threat.
	GetThreats(ctx context.Context) ([]*m.Threat, error)

	// Creates a threat. `threat` should not have Id or ThreatId set.
	CreateThreat(ctx context.Context, threat m.Threat) (*m.Threat, error)

	// Updates a threat
	UpdateThreat(ctx context.Context, threatId m.ThreatId, threat m.Threat) error
}

type DefaultThreatService struct {
	dao dao.ThreatDao

	randomIdProvider id.RandomIdProvider
	validator        validator.StructValidator
	requestor        requestor.Requestor
}

func NewDefaultThreatService(dao dao.ThreatDao, randomIdProvider id.RandomIdProvider, validator validator.StructValidator, requestor requestor.Requestor) *DefaultThreatService {
	return &DefaultThreatService{dao, randomIdProvider, validator, requestor}
}

func (g *DefaultThreatService) GetThreat(ctx context.Context, threatId m.ThreatId) (*m.Threat, error) {
	threat, err := g.dao.Get(ctx, threatId.String())

	if err != nil {
		return nil, fmt.Errorf("error retrieving threat: %v", err)
	}

	return threat, nil
}

// CreateThreat Creates a new Threat in Firestore.
//
// The threat supplied should not have its Id or ThreatId fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created threat is returned to the caller, with Id and ThreatId set.
func (g *DefaultThreatService) CreateThreat(ctx context.Context, threat m.Threat) (*m.Threat, error) {
	if threat.ThreatId != "" {
		return nil, fmt.Errorf("cannot create a threat that already has ThreatId set")
	}

	err := g.validator.ValidateForCreate(threat)
	if err != nil {
		return nil, fmt.Errorf("threat is not valid: %v", err)
	}
	err = g.validator.ValidateForUpdate(threat)
	if err != nil {
		return nil, fmt.Errorf("threat is not valid: %v", err)
	}

	threat.ThreatId = m.ThreatId("th_" + g.randomIdProvider.GenerateId())

	result, err := g.dao.Create(ctx, &threat)
	if err != nil {
		return nil, fmt.Errorf("error creating threat: %v", err)
	}

	return result, nil
}

func (g *DefaultThreatService) UpdateThreat(ctx context.Context, threatId m.ThreatId, threat m.Threat) error {
	err := g.validator.ValidateForUpdate(threat)
	if err != nil {
		return fmt.Errorf("threat is not valid: %v", err)
	}

	threat.ThreatId = m.ThreatId("th_" + g.randomIdProvider.GenerateId())
	queryThreat := m.Threat{ThreatId: threatId}

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
