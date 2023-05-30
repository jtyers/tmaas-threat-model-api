package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	m "github.com/jtyers/tmaas-model"
	"github.com/jtyers/tmaas-model/validator"
	servicedao "github.com/jtyers/tmaas-service-dao"
	"github.com/jtyers/tmaas-service-util/idchecker"
	"github.com/jtyers/tmaas-threat-model-api/dao"
	"github.com/stretchr/testify/require"
)

func TestGetThreatModel(t *testing.T) {
	threatModel := m.ThreatModel{
		ThreatModelID: m.ThreatModelID("1234-1234-1234-1234"),
	}

	var tests = []struct {
		name               string
		inputThreatModelID m.ThreatModelID
		daoReturnValue     m.ThreatModel
		daoReturnError     error
		expectedResult     *m.ThreatModel
		expectedError      error
	}{
		{
			"should get existing threatModels",
			threatModel.ThreatModelID,
			threatModel,
			nil,
			&threatModel,
			nil,
		},
		{
			"should return ErrNoSuchThreatModel for non-existent threatModels",
			threatModel.ThreatModelID,
			threatModel,
			servicedao.ErrNoSuchDocument,
			nil,
			ErrNoSuchThreatModel,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatModelDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().Get(ctx, test.inputThreatModelID).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultThreatModelService(mockDao, nil, nil, nil)
			g, err := service.GetThreatModel(ctx, test.inputThreatModelID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateThreatModel(t *testing.T) {
	threatModel := m.ThreatModel{
		ThreatModelID:     m.ThreatModelID("1234-1234-1234-1234"),
		DataFlowDiagramID: m.DataFlowDiagramID("1234-1234-1234-1234"),
	}

	var tests = []struct {
		name                      string
		inputID                   m.ThreatModelID
		input                     m.ThreatModel
		daoReturnError            error
		validateUpdateReturnError error
		checkIDResult             bool
		checkIDError              error
		expectedResult            *m.ThreatModel
		expectedError             error
	}{
		{
			"should update threatModel",
			threatModel.ThreatModelID,
			threatModel,
			nil,
			nil,
			true,
			nil,
			&threatModel,
			nil,
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			threatModel.ThreatModelID,
			threatModel,
			nil,
			fmt.Errorf("invalid"),
			true,
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			threatModel.ThreatModelID,
			threatModel,
			fmt.Errorf("dao failure"),
			nil,
			true,
			nil,
			nil,
			fmt.Errorf("error updating threatModel: dao failure"),
		},
		{
			"should fail if IDChecker fails",
			threatModel.ThreatModelID,
			threatModel,
			nil,
			nil,
			false,
			fmt.Errorf("failure"),
			nil,
			fmt.Errorf("CheckID failed: failure"),
		},
		{
			"should fail if CheckID() returns false",
			threatModel.ThreatModelID,
			threatModel,
			nil,
			nil,
			false,
			nil,
			nil,
			fmt.Errorf("threatModel.DataFlowDiagramID %v does not exist", threatModel.DataFlowDiagramID),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatModelDao(ctrl)
			ctx := context.Background()
			mockIDChecker := idchecker.NewMockIDChecker(ctrl)
			mockValidator := validator.NewMockStructValidator(ctrl)

			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				mockIDChecker.EXPECT().CheckID(ctx, test.input.DataFlowDiagramID).Return(test.checkIDResult, test.checkIDError)

				if test.checkIDResult && test.checkIDError == nil {
					mockDao.EXPECT().Update(ctx, test.inputID, &test.input).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultThreatModelService(mockDao, mockValidator, mockIDChecker, nil)
			err := service.UpdateThreatModel(ctx, test.inputID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateThreatModel(t *testing.T) {
	threatModel := m.ThreatModel{}

	var tests = []struct {
		name                      string
		input                     m.ThreatModel
		daoReturnError            error
		checkIDResult             bool
		checkIDError              error
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.ThreatModel
		expectedError             error
	}{
		{
			"should create threatModel",
			threatModel,
			nil,
			true,
			nil,
			nil,
			nil,
			&threatModel,
			nil,
		},
		{
			"should pass through ValidateCreate errors with no wrapping or changes",
			threatModel,
			nil,
			true,
			nil,
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			threatModel,
			nil,
			true,
			nil,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO create fails",
			threatModel,
			fmt.Errorf("dao failure"),
			true,
			nil,
			nil,
			nil,
			nil,
			fmt.Errorf("error creating threatModel: dao failure"),
		},
		{
			"should fail if IDChecker fails",
			threatModel,
			nil,
			false,
			fmt.Errorf("failure"),
			nil,
			nil,
			nil,
			fmt.Errorf("CheckID failed: failure"),
		},
		{
			"should fail if CheckID() returns false",
			threatModel,
			nil,
			false,
			nil,
			nil,
			nil,
			nil,
			fmt.Errorf("threatModel.DataFlowDiagramID %v does not exist", threatModel.DataFlowDiagramID),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatModelDao(ctrl)
			ctx := context.Background()
			mockValidator := validator.NewMockStructValidator(ctrl)
			mockIDChecker := idchecker.NewMockIDChecker(ctrl)

			mockValidator.EXPECT().ValidateForCreate(test.input).Return(test.validateCreateReturnError)

			if test.validateCreateReturnError == nil {
				mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

				if test.validateUpdateReturnError == nil {
					mockIDChecker.EXPECT().CheckID(ctx, test.input.DataFlowDiagramID).Return(test.checkIDResult, test.checkIDError)

					if test.checkIDResult && test.checkIDError == nil {

						mockDao.EXPECT().Create(ctx, &test.input).Return(test.expectedResult, test.daoReturnError)

					}
				}
			}

			// when
			service := NewDefaultThreatModelService(mockDao, mockValidator, mockIDChecker, nil)
			g, err := service.CreateThreatModel(ctx, test.input)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetThreatModels(t *testing.T) {
	threatModels := []*m.ThreatModel{
		{
			ThreatModelID: m.ThreatModelID("1234-1234-1234-1234"),
		},
		{
			ThreatModelID: m.ThreatModelID("2345-2345-2345-3245"),
		},
		{
			ThreatModelID: m.ThreatModelID("3456-3456-3456-3456"),
		},
	}

	var tests = []struct {
		name           string
		daoReturnValue []*m.ThreatModel
		daoReturnError error
		expectedResult []*m.ThreatModel
		expectedError  error
	}{
		{
			"should get existing threatModels",
			threatModels,
			nil,
			threatModels,
			nil,
		},
		{
			"should pass through DAO errors",
			threatModels,
			fmt.Errorf("foo bar"),
			nil,
			fmt.Errorf("error retrieving threatModels: foo bar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatModelDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().GetAll(ctx).Return(test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultThreatModelService(mockDao, nil, nil, nil)
			g, err := service.GetThreatModels(ctx)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
