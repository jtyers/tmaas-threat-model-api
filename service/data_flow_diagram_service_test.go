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

func TestGetDataFlowDiagram(t *testing.T) {
	dataFlowDiagram := m.DataFlowDiagram{
		DataFlowDiagramID: m.DataFlowDiagramID("1234-1234-1234-1234"),
	}

	var tests = []struct {
		name                   string
		inputDataFlowDiagramID m.DataFlowDiagramID
		daoReturnValue         m.DataFlowDiagram
		daoReturnError         error
		expectedResult         *m.DataFlowDiagram
		expectedError          error
	}{
		{
			"should get existing dataFlowDiagrams",
			dataFlowDiagram.DataFlowDiagramID,
			dataFlowDiagram,
			nil,
			&dataFlowDiagram,
			nil,
		},
		{
			"should return ErrNoSuchDataFlowDiagram for non-existent dataFlowDiagrams",
			dataFlowDiagram.DataFlowDiagramID,
			dataFlowDiagram,
			servicedao.ErrNoSuchDocument,
			nil,
			ErrNoSuchDataFlowDiagram,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockDataFlowDiagramDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().Get(ctx, test.inputDataFlowDiagramID.String()).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultDataFlowDiagramService(mockDao, nil, nil)
			g, err := service.GetDataFlowDiagram(ctx, test.inputDataFlowDiagramID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateDataFlowDiagram(t *testing.T) {
	dataFlowDiagram := m.DataFlowDiagram{
		RootComponentID: "cm_123456",
	}

	var tests = []struct {
		name                      string
		inputID                   m.DataFlowDiagramID
		input                     m.DataFlowDiagram
		daoReturnError            error
		validateUpdateReturnError error
		expectedResult            *m.DataFlowDiagram
		expectedError             error
	}{
		{
			"should update dataFlowDiagram",
			dataFlowDiagram.DataFlowDiagramID,
			dataFlowDiagram,
			nil,
			nil,
			&dataFlowDiagram,
			nil,
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			dataFlowDiagram.DataFlowDiagramID,
			dataFlowDiagram,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			dataFlowDiagram.DataFlowDiagramID,
			dataFlowDiagram,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			fmt.Errorf("error updating dataFlowDiagram: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockDataFlowDiagramDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				queryDataFlowDiagram := m.DataFlowDiagram{DataFlowDiagramID: test.inputID}

				mockDao.EXPECT().UpdateWhereExactSingle(ctx, &queryDataFlowDiagram, &test.input).Return(test.expectedResult, test.daoReturnError)

			}

			// when
			service := NewDefaultDataFlowDiagramService(mockDao, mockIDProvider, mockValidator)
			err := service.UpdateDataFlowDiagram(ctx, test.inputID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateDataFlowDiagram(t *testing.T) {
	dataFlowDiagram := m.DataFlowDiagram{
		RootComponentID: "cm_123456",
	}

	var tests = []struct {
		name                      string
		input                     m.DataFlowDiagram
		daoReturnError            error
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.DataFlowDiagram
		expectedError             error
	}{
		{
			"should create dataFlowDiagram",
			dataFlowDiagram,
			nil,
			nil,
			nil,
			&dataFlowDiagram,
			nil,
		},
		{
			"should pass through ValidateCreate errors with no wrapping or changes",
			dataFlowDiagram,
			nil,
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			dataFlowDiagram,
			nil,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO create fails",
			dataFlowDiagram,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			nil,
			fmt.Errorf("error creating dataFlowDiagram: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockDataFlowDiagramDao(ctrl)
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
					expectedInputForCreate.DataFlowDiagramID = m.DataFlowDiagramID(DataFlowDiagramIDPrefix + newID)

					mockDao.EXPECT().Create(ctx, &expectedInputForCreate).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultDataFlowDiagramService(mockDao, mockIDProvider, mockValidator)
			g, err := service.CreateDataFlowDiagram(ctx, test.input)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetDataFlowDiagrams(t *testing.T) {
	dataFlowDiagrams := []*m.DataFlowDiagram{
		{
			DataFlowDiagramID: m.DataFlowDiagramID("1234-1234-1234-1234"),
			RootComponentID:   "cm_123456",
		},
		{
			DataFlowDiagramID: m.DataFlowDiagramID("2345-2345-2345-3245"),
			RootComponentID:   "cm_234567",
		},
		{
			DataFlowDiagramID: m.DataFlowDiagramID("3456-3456-3456-3456"),
			RootComponentID:   "cm_345678",
		},
	}

	var tests = []struct {
		name           string
		daoReturnValue []*m.DataFlowDiagram
		daoReturnError error
		expectedResult []*m.DataFlowDiagram
		expectedError  error
	}{
		{
			"should get existing dataFlowDiagrams",
			dataFlowDiagrams,
			nil,
			dataFlowDiagrams,
			nil,
		},
		{
			"should pass through DAO errors",
			dataFlowDiagrams,
			fmt.Errorf("foo bar"),
			nil,
			fmt.Errorf("error retrieving dataFlowDiagrams: foo bar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockDataFlowDiagramDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().GetAll(ctx).Return(test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultDataFlowDiagramService(mockDao, nil, nil)
			g, err := service.GetDataFlowDiagrams(ctx)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
