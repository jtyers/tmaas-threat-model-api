package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/validator"
	servicedao "github.com/jtyers/tmaas-service-dao"
	"github.com/jtyers/tmaas-service-util/id"
	"github.com/jtyers/tmaas-api/dao"
	"github.com/stretchr/testify/require"
)

func TestGetComponent(t *testing.T) {
	component := m.Component{
		ComponentID: m.ComponentID("1234-1234-1234-1234"),
		Name:        "my-first-component",
	}

	var tests = []struct {
		name             string
		inputComponentID m.ComponentID
		daoReturnValue   m.Component
		daoReturnError   error
		expectedResult   *m.Component
		expectedError    error
	}{
		{
			"should get existing components",
			component.ComponentID,
			component,
			nil,
			&component,
			nil,
		},
		{
			"should return ErrNoSuchComponent for non-existent components",
			component.ComponentID,
			component,
			servicedao.ErrNoSuchDocument,
			nil,
			ErrNoSuchComponent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockComponentDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().Get(ctx, test.inputComponentID.String()).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultComponentService(mockDao, nil, nil)
			g, err := service.GetComponent(ctx, test.inputComponentID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateComponent(t *testing.T) {
	component := m.Component{
		Name:              "my-first-component",
		ParentComponentID: "cm_123456",
	}

	var tests = []struct {
		name                      string
		inputID                   m.ComponentID
		input                     m.Component
		daoReturnError            error
		validateUpdateReturnError error
		expectedResult            *m.Component
		expectedError             error
	}{
		{
			"should update component",
			component.ComponentID,
			component,
			nil,
			nil,
			&component,
			nil,
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			component.ComponentID,
			component,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			component.ComponentID,
			component,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			fmt.Errorf("error updating component: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockComponentDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				queryComponent := m.Component{ComponentID: test.inputID}

				mockDao.EXPECT().UpdateWhereExactSingle(ctx, &queryComponent, &test.input).Return(test.expectedResult, test.daoReturnError)

			}

			// when
			service := NewDefaultComponentService(mockDao, mockIDProvider, mockValidator)
			err := service.UpdateComponent(ctx, test.inputID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateComponent(t *testing.T) {
	component := m.Component{
		Name:              "my-first-component",
		ParentComponentID: "cm_123456",
	}

	var tests = []struct {
		name                      string
		input                     m.Component
		daoReturnError            error
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.Component
		expectedError             error
	}{
		{
			"should create component",
			component,
			nil,
			nil,
			nil,
			&component,
			nil,
		},
		{
			"should pass through ValidateCreate errors with no wrapping or changes",
			component,
			nil,
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			component,
			nil,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO create fails",
			component,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			nil,
			fmt.Errorf("error creating component: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockComponentDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForCreate(test.input).Return(test.validateCreateReturnError)

			if test.validateCreateReturnError == nil {
				mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

				if test.validateUpdateReturnError == nil {
					newID := "1234-1234"
					mockIDProvider.EXPECT().GenerateID().Return(newID)

					// NOTE: this is not a pointer, so a copy of the original struct
					expectedInputForCreate := test.input
					expectedInputForCreate.ComponentID = m.ComponentID(ComponentIDPrefix + newID)

					mockDao.EXPECT().Create(ctx, &expectedInputForCreate).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultComponentService(mockDao, mockIDProvider, mockValidator)
			g, err := service.CreateComponent(ctx, test.input)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetComponents(t *testing.T) {
	components := []*m.Component{
		{
			ComponentID: m.ComponentID("1234-1234-1234-1234"),
			Name:        "my-first-component",
		},
		{
			ComponentID: m.ComponentID("2345-2345-2345-3245"),
			Name:        "my-second-component",
		},
		{
			ComponentID: m.ComponentID("3456-3456-3456-3456"),
			Name:        "my-third-component",
		},
	}

	var tests = []struct {
		name           string
		dfdID          m.DataFlowDiagramID
		daoReturnValue []*m.Component
		daoReturnError error
		expectedResult []*m.Component
		expectedError  error
	}{
		{
			"should get existing components",
			m.DataFlowDiagramID("dfd_1"),
			components,
			nil,
			components,
			nil,
		},
		{
			"should pass through DAO errors",
			m.DataFlowDiagramID("dfd_1"),
			components,
			fmt.Errorf("foo bar"),
			nil,
			fmt.Errorf("error retrieving components: foo bar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockComponentDao(ctrl)
			ctx := context.Background()

			queryComponent := &m.Component{DataFlowDiagramID: test.dfdID}
			mockDao.EXPECT().QueryExact(ctx, queryComponent).Return(test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultComponentService(mockDao, nil, nil)
			g, err := service.GetComponents(ctx, test.dfdID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
