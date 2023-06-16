package service

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mocks.go -package $GOPACKAGE

import (
	"context"
	"errors"
	"fmt"

	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/validator"
	servicedao "github.com/jtyers/tmaas-service-dao"
	"github.com/jtyers/tmaas-service-util/idchecker"
	dao "github.com/jtyers/tmaas-threat-model-api/dao"
)

var (
	ErrNoSuchThreatModel = errors.New("no such threat model")
)

// ThreatModelService provides the interface to manage threat models.
type ThreatModelService interface {
	// Retrieve a threatModel by threatModelID.
	Get(ctx context.Context, id m.ThreatModelID) (*m.ThreatModel, error)

	// Retrieve all threatModel.
	GetAll(ctx context.Context) ([]*m.ThreatModel, error)

	Query(ctx context.Context, q *m.ThreatModelQuery) ([]*m.ThreatModel, error)

	QuerySingle(ctx context.Context, q *m.ThreatModelQuery) (*m.ThreatModel, error)

	// Creates a threatModel. `threatModel` should not have ID or threatModelID set.
	Create(ctx context.Context, params m.ThreatModelParams) (*m.ThreatModel, error)

	// Updates a threatModel
	Update(ctx context.Context, id m.ThreatModelID, params m.ThreatModelParams) error

	// Delete a threatModel by threatModelID.
	Delete(ctx context.Context, id m.ThreatModelID) error
}

type DefaultThreatModelService struct {
	dao       dao.ThreatModelDao
	validator validator.StructValidator
	idChecker idchecker.IDChecker
}

var _ ThreatModelService = (*DefaultThreatModelService)(nil)

func NewDefaultThreatModelService(
	dao dao.ThreatModelDao,
	validator validator.StructValidator,
	idChecker idchecker.IDChecker,
) *DefaultThreatModelService {
	return &DefaultThreatModelService{dao, validator, idChecker}
}

func (g *DefaultThreatModelService) Get(ctx context.Context, threatModelID m.ThreatModelID) (*m.ThreatModel, error) {
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
// The threatModel supplied should not have its ID or ThreatModelID
// fields set to anything other than "". An error is emitted if this
// is not the case.
//
// The created threatModel is returned to the caller, with ID and
// ThreatModelID set.
func (g *DefaultThreatModelService) Create(ctx context.Context, params m.ThreatModelParams) (*m.ThreatModel, error) {
	err := g.validator.ValidateForCreate(params)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(params)
	if err != nil {
		return nil, err
	}

	if params.DataFlowDiagramID != nil {
		exists, err := g.idChecker.CheckID(ctx, params.DataFlowDiagramID)
		if err != nil {
			return nil, fmt.Errorf("CheckID failed: %v", err)
		}
		if !exists {
			return nil, fmt.Errorf("threatModel.DataFlowDiagramID %v does not exist",
				params.DataFlowDiagramID)
		}
	}

	// leave ID blank - the DAO will generate one for us
	result, err := g.dao.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating threatModel: %v", err)
	}

	return result, nil
}

func (g *DefaultThreatModelService) Update(ctx context.Context, threatModelID m.ThreatModelID, params m.ThreatModelParams) error {
	err := g.validator.ValidateForUpdate(params)
	if err != nil {
		return err
	}

	if params.DataFlowDiagramID != nil {
		exists, err := g.idChecker.CheckID(ctx, params.DataFlowDiagramID)
		if err != nil {
			return fmt.Errorf("CheckID failed: %v", err)
		}
		if !exists {
			return fmt.Errorf("threatModel.DataFlowDiagramID %v does not exist",
				params.DataFlowDiagramID)
		}
	}

	_, err = g.dao.Update(ctx, threatModelID, params)
	if err != nil {
		return fmt.Errorf("error updating threatModel: %v", err)
	}

	return nil
}

func (g *DefaultThreatModelService) GetAll(ctx context.Context) ([]*m.ThreatModel, error) {
	threatModels, err := g.dao.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("error retrieving threatModels: %v", err)
	}

	return threatModels, nil
}

func (g *DefaultThreatModelService) Delete(ctx context.Context, id m.ThreatModelID) error {
	err := g.dao.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting threatModel %s: %v", id, err)
	}

	return nil
}

func (g *DefaultThreatModelService) Query(ctx context.Context, q *m.ThreatModelQuery) ([]*m.ThreatModel, error) {
	result, err := g.dao.QueryExact(ctx, q)

	if err != nil {
		return nil, fmt.Errorf("error in QueryExact: %v", err)
	}

	return result, nil
}

func (g *DefaultThreatModelService) QuerySingle(ctx context.Context, q *m.ThreatModelQuery) (*m.ThreatModel, error) {
	result, err := g.dao.QueryExactSingle(ctx, q)

	if err != nil {
		return nil, fmt.Errorf("error in QueryExact: %v", err)
	}

	return result, nil
}
