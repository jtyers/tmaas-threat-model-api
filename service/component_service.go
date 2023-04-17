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
	ErrTooManyComponents = errors.New("too many components")
	ErrNoSuchComponent   = errors.New("no such component")
)

type ComponentService interface {
	// Retrieve a component by componentID.
	GetComponent(ctx context.Context, id m.ComponentID) (*m.Component, error)

	// Retrieve all component for a given DFD.
	GetComponents(ctx context.Context, id m.DataFlowDiagramID) ([]*m.Component, error)

	// Creates a component. `component` should not have ID or ComponentID set.
	CreateComponent(ctx context.Context, component m.Component) (*m.Component, error)

	// Updates a component
	UpdateComponent(ctx context.Context, componentID m.ComponentID, component m.Component) error
}

type DefaultComponentService struct {
	dao dao.ComponentDao

	randomIDProvider id.RandomIDProvider
	validator        validator.StructValidator
}

func NewDefaultComponentService(dao dao.ComponentDao, randomIDProvider id.RandomIDProvider, validator validator.StructValidator) *DefaultComponentService {
	return &DefaultComponentService{dao, randomIDProvider, validator}
}

func (g *DefaultComponentService) GetComponent(ctx context.Context, componentID m.ComponentID) (*m.Component, error) {
	component, err := g.dao.Get(ctx, componentID.String())

	if err != nil {
		if err == servicedao.ErrNoSuchDocument {
			return nil, ErrNoSuchComponent
		} else {
			return nil, fmt.Errorf("error retrieving component: %v", err)
		}
	}

	return component, nil
}

// CreateComponent Creates a new Component in Firestore.
//
// The component supplied should not have its ID or ComponentID fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created component is returned to the caller, with ID and ComponentID set.
func (g *DefaultComponentService) CreateComponent(ctx context.Context, component m.Component) (*m.Component, error) {
	if component.ComponentID != "" {
		return nil, fmt.Errorf("cannot create a component that already has ComponentID set")
	}

	err := g.validator.ValidateForCreate(component)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(component)
	if err != nil {
		return nil, err
	}

	component.ComponentID = m.ComponentID(ComponentIDPrefix + g.randomIDProvider.GenerateID())

	result, err := g.dao.Create(ctx, &component)
	if err != nil {
		return nil, fmt.Errorf("error creating component: %v", err)
	}

	return result, nil
}

func (g *DefaultComponentService) UpdateComponent(ctx context.Context, componentID m.ComponentID, component m.Component) error {
	err := g.validator.ValidateForUpdate(component)
	if err != nil {
		return err
	}

	if component.ComponentID != componentID {
		return fmt.Errorf("given component IDs do not match")
	}

	queryComponent := m.Component{ComponentID: componentID}

	_, err = g.dao.UpdateWhereExactSingle(ctx, &queryComponent, &component)
	if err != nil {
		return fmt.Errorf("error updating component: %v", err)
	}

	return nil
}

func (g *DefaultComponentService) GetComponents(ctx context.Context, id m.DataFlowDiagramID) ([]*m.Component, error) {
	queryComponent := &m.Component{DataFlowDiagramID: id}
	components, err := g.dao.QueryExact(ctx, queryComponent)

	if err != nil {
		return nil, fmt.Errorf("error retrieving components: %v", err)
	}

	return components, nil
}
