package service

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mocks.go -package $GOPACKAGE

import (
	"context"
	"errors"
	"fmt"

	dao "github.com/jtyers/tmaas-threat-model-api/dao"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/validator"
	servicedao "github.com/jtyers/tmaas-service-dao"
)

var (
	ErrNoSuchThreatModel = errors.New("no such threat model")
)

// ThreatModelService provides the interface to manage threat models.
type ThreatModelService interface {
	// Retrieve a threatModel by threatModelID.
	GetThreatModel(ctx context.Context, id m.ThreatModelID) (*m.ThreatModel, error)

	// Retrieve all threatModel.
	GetThreatModels(ctx context.Context) ([]*m.ThreatModel, error)

	// Creates a threatModel. `threatModel` should not have ID or threatModelID set.
	CreateThreatModel(ctx context.Context, threatModel m.ThreatModel) (*m.ThreatModel, error)

	// Updates a threatModel
	UpdateThreatModel(ctx context.Context, threatModelID m.ThreatModelID, threatModel m.ThreatModel) error

	// Delete a threatModel by threatModelID.
	DeleteThreatModel(ctx context.Context, id m.ThreatModelID) error
}

type DefaultThreatModelService struct {
	dao       dao.ThreatModelDao
	validator validator.StructValidator
}

func NewDefaultThreatModelService(dao dao.ThreatModelDao, validator validator.StructValidator) *DefaultThreatModelService {
	return &DefaultThreatModelService{dao, validator}
}

func (g *DefaultThreatModelService) GetThreatModel(ctx context.Context, threatModelID m.ThreatModelID) (*m.ThreatModel, error) {
	threatModel, err := g.dao.Get(ctx, threatModelID)

	if err != nil {
		if err == servicedao.ErrNoSuchDocument {
			return nil, ErrNoSuchThreatModel
		}
		return nil, fmt.Errorf("error retrieving threatModel: %v", err)
	}

	return threatModel, nil
}

// CreateThreatModel Creates a new ThreatModel in Firestore.
//
// The threatModel supplied should not have its ID or ThreatModelID fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created threatModel is returned to the caller, with ID and ThreatModelID set.
func (g *DefaultThreatModelService) CreateThreatModel(ctx context.Context, threatModel m.ThreatModel) (*m.ThreatModel, error) {
	if threatModel.ThreatModelID != "" {
		return nil, fmt.Errorf("cannot create a threatModel that already has ThreatModelID set")
	}

	err := g.validator.ValidateForCreate(threatModel)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(threatModel)
	if err != nil {
		return nil, err
	}

	// leave ID blank - the DAO will generate one for us
	result, err := g.dao.Create(ctx, &threatModel)
	if err != nil {
		return nil, fmt.Errorf("error creating threatModel: %v", err)
	}

	return result, nil
}

func (g *DefaultThreatModelService) UpdateThreatModel(ctx context.Context, threatModelID m.ThreatModelID, threatModel m.ThreatModel) error {
	err := g.validator.ValidateForUpdate(threatModel)
	if err != nil {
		return err
	}

	if threatModel.ThreatModelID != threatModelID {
		return fmt.Errorf("given threatModel IDs do not match")
	}

	_, err = g.dao.Update(ctx, threatModelID, &threatModel)
	if err != nil {
		return fmt.Errorf("error updating threatModel: %v", err)
	}

	return nil
}

func (g *DefaultThreatModelService) GetThreatModels(ctx context.Context) ([]*m.ThreatModel, error) {
	threatModels, err := g.dao.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("error retrieving threatModels: %v", err)
	}

	return threatModels, nil
}

func (g *DefaultThreatModelService) DeleteThreatModel(ctx context.Context, id m.ThreatModelID) error {
	err := g.dao.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting threatModel %s: %v", id, err)
	}

	return nil
}
