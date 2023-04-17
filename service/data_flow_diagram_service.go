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
	dao "github.com/jtyers/tmaas-api/dao"
)

var (
	ErrTooManyDataFlowDiagrams = errors.New("too many dataFlowDiagrams")
	ErrNoSuchDataFlowDiagram   = errors.New("no such dataFlowDiagram")
)

type DataFlowDiagramService interface {
	// Retrieve a dataFlowDiagram by dataFlowDiagramID.
	GetDataFlowDiagram(ctx context.Context, id m.DataFlowDiagramID) (*m.DataFlowDiagram, error)

	// Retrieve all dataFlowDiagram.
	GetDataFlowDiagrams(ctx context.Context) ([]*m.DataFlowDiagram, error)

	// Creates a dataFlowDiagram. `dataFlowDiagram` should not have ID or DataFlowDiagramID set.
	CreateDataFlowDiagram(ctx context.Context, dataFlowDiagram m.DataFlowDiagram) (*m.DataFlowDiagram, error)

	// Updates a dataFlowDiagram
	UpdateDataFlowDiagram(ctx context.Context, dataFlowDiagramID m.DataFlowDiagramID, dataFlowDiagram m.DataFlowDiagram) error
}

type DefaultDataFlowDiagramService struct {
	dao dao.DataFlowDiagramDao

	randomIDProvider id.RandomIDProvider
	validator        validator.StructValidator
}

func NewDefaultDataFlowDiagramService(dao dao.DataFlowDiagramDao, randomIDProvider id.RandomIDProvider, validator validator.StructValidator) *DefaultDataFlowDiagramService {
	return &DefaultDataFlowDiagramService{dao, randomIDProvider, validator}
}

func (g *DefaultDataFlowDiagramService) GetDataFlowDiagram(ctx context.Context, dataFlowDiagramID m.DataFlowDiagramID) (*m.DataFlowDiagram, error) {
	dataFlowDiagram, err := g.dao.Get(ctx, dataFlowDiagramID.String())

	if err != nil {
		if err == servicedao.ErrNoSuchDocument {
			return nil, ErrNoSuchDataFlowDiagram
		} else {
			return nil, fmt.Errorf("error retrieving dataFlowDiagram: %v", err)
		}
	}

	return dataFlowDiagram, nil
}

// CreateDataFlowDiagram Creates a new DataFlowDiagram in Firestore.
//
// The dataFlowDiagram supplied should not have its ID or DataFlowDiagramID fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created dataFlowDiagram is returned to the caller, with ID and DataFlowDiagramID set.
func (g *DefaultDataFlowDiagramService) CreateDataFlowDiagram(ctx context.Context, dataFlowDiagram m.DataFlowDiagram) (*m.DataFlowDiagram, error) {
	if dataFlowDiagram.DataFlowDiagramID != "" {
		return nil, fmt.Errorf("cannot create a dataFlowDiagram that already has DataFlowDiagramID set")
	}

	err := g.validator.ValidateForCreate(dataFlowDiagram)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(dataFlowDiagram)
	if err != nil {
		return nil, err
	}

	dataFlowDiagram.DataFlowDiagramID = m.DataFlowDiagramID(DataFlowDiagramIDPrefix + g.randomIDProvider.GenerateID())

	result, err := g.dao.Create(ctx, &dataFlowDiagram)
	if err != nil {
		return nil, fmt.Errorf("error creating dataFlowDiagram: %v", err)
	}

	return result, nil
}

func (g *DefaultDataFlowDiagramService) UpdateDataFlowDiagram(ctx context.Context, dataFlowDiagramID m.DataFlowDiagramID, dataFlowDiagram m.DataFlowDiagram) error {
	err := g.validator.ValidateForUpdate(dataFlowDiagram)
	if err != nil {
		return err
	}

	if dataFlowDiagram.DataFlowDiagramID != dataFlowDiagramID {
		return fmt.Errorf("given dataFlowDiagram IDs do not match")
	}

	queryDataFlowDiagram := m.DataFlowDiagram{DataFlowDiagramID: dataFlowDiagramID}

	_, err = g.dao.UpdateWhereExactSingle(ctx, &queryDataFlowDiagram, &dataFlowDiagram)
	if err != nil {
		return fmt.Errorf("error updating dataFlowDiagram: %v", err)
	}

	return nil
}

func (g *DefaultDataFlowDiagramService) GetDataFlowDiagrams(ctx context.Context) ([]*m.DataFlowDiagram, error) {
	dataFlowDiagrams, err := g.dao.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("error retrieving dataFlowDiagrams: %v", err)
	}

	return dataFlowDiagrams, nil
}
