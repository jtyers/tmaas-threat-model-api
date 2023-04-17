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

			mockDao.EXPECT().Get(ctx, test.inputThreatModelID.String()).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultThreatModelService(mockDao, nil, nil)
			g, err := service.GetThreatModel(ctx, test.inputThreatModelID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateThreatModel(t *testing.T) {
	threatModel := m.ThreatModel{
		ThreatModelID: m.ThreatModelID("1234-1234-1234-1234"),
	}

	var tests = []struct {
		name                      string
		inputID                   m.ThreatModelID
		input                     m.ThreatModel
		daoReturnError            error
		validateUpdateReturnError error
		expectedResult            *m.ThreatModel
		expectedError             error
	}{
		{
			"should update threatModel",
			threatModel.ThreatModelID,
			threatModel,
			nil,
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
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			threatModel.ThreatModelID,
			threatModel,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			fmt.Errorf("error updating threatModel: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatModelDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				queryThreatModel := m.ThreatModel{ThreatModelID: test.inputID}

				mockDao.EXPECT().UpdateWhereExactSingle(ctx, &queryThreatModel, &test.input).Return(test.expectedResult, test.daoReturnError)

			}

			// when
			service := NewDefaultThreatModelService(mockDao, mockIDProvider, mockValidator)
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
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.ThreatModel
		expectedError             error
	}{
		{
			"should create threatModel",
			threatModel,
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
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			threatModel,
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
			nil,
			nil,
			nil,
			fmt.Errorf("error creating threatModel: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatModelDao(ctrl)
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
					expectedInputForCreate.ThreatModelID = m.ThreatModelID(ThreatModelIDPrefix + newID)

					mockDao.EXPECT().Create(ctx, &expectedInputForCreate).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultThreatModelService(mockDao, mockIDProvider, mockValidator)
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
			service := NewDefaultThreatModelService(mockDao, nil, nil)
			g, err := service.GetThreatModels(ctx)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
