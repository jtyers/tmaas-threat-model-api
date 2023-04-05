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
	ErrTooManyControls = errors.New("too many controls")
	ErrNoSuchControl   = errors.New("no such control")
)

type ControlService interface {
	// Retrieve a control by controlID.
	GetControl(ctx context.Context, id m.ControlID) (*m.Control, error)

	// Retrieve all control.
	GetControls(ctx context.Context) ([]*m.Control, error)

	// Creates a control. `control` should not have ID or ControlID set.
	CreateControl(ctx context.Context, control m.Control) (*m.Control, error)

	// Updates a control
	UpdateControl(ctx context.Context, controlID m.ControlID, control m.Control) error

	// Deletes a control.
	DeleteControl(ctx context.Context, controlID m.ControlID) error
}

type DefaultControlService struct {
	dao dao.ControlDao

	randomIDProvider id.RandomIDProvider
	validator        validator.StructValidator
}

var _ ControlService = (*DefaultControlService)(nil)

func NewDefaultControlService(dao dao.ControlDao, randomIDProvider id.RandomIDProvider, validator validator.StructValidator) *DefaultControlService {
	return &DefaultControlService{dao, randomIDProvider, validator}
}

func (g *DefaultControlService) GetControl(ctx context.Context, controlID m.ControlID) (*m.Control, error) {
	control, err := g.dao.Get(ctx, controlID.String())

	if err != nil {
		if err == servicedao.ErrNoSuchDocument {
			return nil, ErrNoSuchControl
		} else {
			return nil, fmt.Errorf("error retrieving control: %v", err)
		}
	}

	return control, nil
}

// CreateControl Creates a new Control in Firestore.
//
// The control supplied should not have its ID or ControlID fields set to anything
// other than "". An error is emitted if this is not the case.
//
// The created control is returned to the caller, with ID and ControlID set.
func (g *DefaultControlService) CreateControl(ctx context.Context, control m.Control) (*m.Control, error) {
	if control.ControlID != "" {
		return nil, fmt.Errorf("cannot create a control that already has ControlID set")
	}

	err := g.validator.ValidateForCreate(control)
	if err != nil {
		return nil, err
	}
	err = g.validator.ValidateForUpdate(control)
	if err != nil {
		return nil, err
	}

	control.ControlID = m.ControlID(ControlIDPrefix + g.randomIDProvider.GenerateID())

	result, err := g.dao.Create(ctx, &control)
	if err != nil {
		return nil, fmt.Errorf("error creating control: %v", err)
	}

	return result, nil
}

func (g *DefaultControlService) UpdateControl(ctx context.Context, controlID m.ControlID, control m.Control) error {
	err := g.validator.ValidateForUpdate(control)
	if err != nil {
		return err
	}

	if control.ControlID != controlID {
		return fmt.Errorf("given control IDs do not match")
	}

	queryControl := m.Control{ControlID: controlID}

	_, err = g.dao.UpdateWhereExactSingle(ctx, &queryControl, &control)
	if err != nil {
		return fmt.Errorf("error updating control: %v", err)
	}

	return nil
}

func (g *DefaultControlService) GetControls(ctx context.Context) ([]*m.Control, error) {
	controls, err := g.dao.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("error retrieving controls: %v", err)
	}

	return controls, nil
}

func (g *DefaultControlService) DeleteControl(ctx context.Context, controlID m.ControlID) error {
	queryControl := m.Control{ControlID: controlID}
	err := g.dao.DeleteWhere(ctx, &queryControl)
	if err != nil {
		return fmt.Errorf("error deleting control: %v", err)
	}

	return nil
}
