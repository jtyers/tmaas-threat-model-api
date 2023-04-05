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
	ErrTooManyInherentRisks = errors.New("too many inherentRisks")
	ErrNoSuchInherentRisk   = errors.New("no such inherentRisk")
)

type InherentRiskService interface {
	// Retrieve a inherentRisk by inherentRiskID.
	GetInherentRisk(ctx context.Context, id m.InherentRiskID) (*m.InherentRisk, error)

	// Retrieve all inherentRisk.
	GetInherentRisks(ctx context.Context) ([]*m.InherentRisk, error)

	// Creates a inherentRisk. `inherentRisk` should not have ID or InherentRiskID set.
	CreateInherentRisk(ctx context.Context, inherentRisk m.InherentRisk) (*m.InherentRisk, error)

	// Updates a inherentRisk
	UpdateInherentRisk(ctx context.Context, inherentRiskID m.InherentRiskID, inherentRisk m.InherentRisk) error
}

type DefaultInherentRiskService struct {
	dao dao.InherentRiskDao

	randomIDProvider id.RandomIDProvider
	validator        validator.StructValidator
}

func NewDefaultInherentRiskService(dao dao.InherentRiskDao, randomIDProvider id.RandomIDProvider, validator validator.StructValidator) *DefaultInherentRiskService {
	return &DefaultInherentRiskService{dao, randomIDProvider, validator}
}

func (g *DefaultInherentRiskService) GetInherentRisk(ctx context.Context, inherentRiskID m.InherentRiskID) (*m.InherentRisk, error) {
	inherentRisk, err := g.dao.Get(ctx, inherentRiskID.String())

	if err != nil {
		if err == servicedao.ErrNoSuchDocument {
			return nil, ErrNoSuchInherentRisk
		} else {
			return nil, fmt.Errorf("error retrieving inherentRisk: %v", err)
		}
	}

	return inherentRisk, nil
}

// CreateInherentRisk Creates a new InherentRisk in Firestore.
//
// The inherentRisk supplied should not have its ID or InherentRiskID fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created inherentRisk is returned to the caller, with ID and InherentRiskID set.
func (g *DefaultInherentRiskService) CreateInherentRisk(ctx context.Context, inherentRisk m.InherentRisk) (*m.InherentRisk, error) {
	if inherentRisk.InherentRiskID != "" {
		return nil, fmt.Errorf("cannot create a inherentRisk that already has InherentRiskID set")
	}

	err := g.validator.ValidateForCreate(inherentRisk)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(inherentRisk)
	if err != nil {
		return nil, err
	}

	inherentRisk.InherentRiskID = m.InherentRiskID(InherentRiskIDPrefix + g.randomIDProvider.GenerateID())

	result, err := g.dao.Create(ctx, &inherentRisk)
	if err != nil {
		return nil, fmt.Errorf("error creating inherentRisk: %v", err)
	}

	return result, nil
}

func (g *DefaultInherentRiskService) UpdateInherentRisk(ctx context.Context, inherentRiskID m.InherentRiskID, inherentRisk m.InherentRisk) error {
	err := g.validator.ValidateForUpdate(inherentRisk)
	if err != nil {
		return err
	}

	if inherentRisk.InherentRiskID != inherentRiskID {
		return fmt.Errorf("given inherentRisk IDs do not match")
	}

	queryInherentRisk := m.InherentRisk{InherentRiskID: inherentRiskID}

	_, err = g.dao.UpdateWhereExactSingle(ctx, &queryInherentRisk, &inherentRisk)
	if err != nil {
		return fmt.Errorf("error updating inherentRisk: %v", err)
	}

	return nil
}

func (g *DefaultInherentRiskService) GetInherentRisks(ctx context.Context) ([]*m.InherentRisk, error) {
	inherentRisks, err := g.dao.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("error retrieving inherentRisks: %v", err)
	}

	return inherentRisks, nil
}
