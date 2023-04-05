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
	"github.com/jtyers/tmaas-threat-service/dao"
	"github.com/stretchr/testify/require"
)

func TestGetInherentRisk(t *testing.T) {
	inherentRisk := m.InherentRisk{
		InherentRiskID: m.InherentRiskID("1234-1234-1234-1234"),
		Risk:           m.RiskLevelHigh,
	}

	var tests = []struct {
		name                string
		inputInherentRiskID m.InherentRiskID
		daoReturnValue      m.InherentRisk
		daoReturnError      error
		expectedResult      *m.InherentRisk
		expectedError       error
	}{
		{
			"should get existing inherentRisks",
			inherentRisk.InherentRiskID,
			inherentRisk,
			nil,
			&inherentRisk,
			nil,
		},
		{
			"should return ErrNoSuchInherentRisk for non-existent inherentRisks",
			inherentRisk.InherentRiskID,
			inherentRisk,
			servicedao.ErrNoSuchDocument,
			nil,
			ErrNoSuchInherentRisk,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockInherentRiskDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().Get(ctx, test.inputInherentRiskID.String()).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultInherentRiskService(mockDao, nil, nil)
			g, err := service.GetInherentRisk(ctx, test.inputInherentRiskID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateInherentRisk(t *testing.T) {
	inherentRisk := m.InherentRisk{
		Risk: m.RiskLevelHigh,
	}

	var tests = []struct {
		name                      string
		inputID                   m.InherentRiskID
		input                     m.InherentRisk
		daoReturnError            error
		validateUpdateReturnError error
		expectedResult            *m.InherentRisk
		expectedError             error
	}{
		{
			"should update inherentRisk",
			inherentRisk.InherentRiskID,
			inherentRisk,
			nil,
			nil,
			&inherentRisk,
			nil,
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			inherentRisk.InherentRiskID,
			inherentRisk,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			inherentRisk.InherentRiskID,
			inherentRisk,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			fmt.Errorf("error updating inherentRisk: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockInherentRiskDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				queryInherentRisk := m.InherentRisk{InherentRiskID: test.inputID}

				mockDao.EXPECT().UpdateWhereExactSingle(ctx, &queryInherentRisk, &test.input).Return(test.expectedResult, test.daoReturnError)

			}

			// when
			service := NewDefaultInherentRiskService(mockDao, mockIDProvider, mockValidator)
			err := service.UpdateInherentRisk(ctx, test.inputID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateInherentRisk(t *testing.T) {
	inherentRisk := m.InherentRisk{
		Risk: m.RiskLevelHigh,
	}

	var tests = []struct {
		name                      string
		input                     m.InherentRisk
		daoReturnError            error
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.InherentRisk
		expectedError             error
	}{
		{
			"should create inherentRisk",
			inherentRisk,
			nil,
			nil,
			nil,
			&inherentRisk,
			nil,
		},
		{
			"should pass through ValidateCreate errors with no wrapping or changes",
			inherentRisk,
			nil,
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			inherentRisk,
			nil,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO create fails",
			inherentRisk,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			nil,
			fmt.Errorf("error creating inherentRisk: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockInherentRiskDao(ctrl)
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
					expectedInputForCreate.InherentRiskID = m.InherentRiskID(InherentRiskIDPrefix + newID)

					mockDao.EXPECT().Create(ctx, &expectedInputForCreate).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultInherentRiskService(mockDao, mockIDProvider, mockValidator)
			g, err := service.CreateInherentRisk(ctx, test.input)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetInherentRisks(t *testing.T) {
	inherentRisks := []*m.InherentRisk{
		{
			InherentRiskID: m.InherentRiskID("1234-1234-1234-1234"),
			Risk:           m.RiskLevelHigh,
		},
		{
			InherentRiskID: m.InherentRiskID("2345-2345-2345-3245"),
			Risk:           m.RiskLevelMedium,
		},
		{
			InherentRiskID: m.InherentRiskID("3456-3456-3456-3456"),
			Risk:           m.RiskLevelLow,
		},
	}

	var tests = []struct {
		name           string
		daoReturnValue []*m.InherentRisk
		daoReturnError error
		expectedResult []*m.InherentRisk
		expectedError  error
	}{
		{
			"should get existing inherentRisks",
			inherentRisks,
			nil,
			inherentRisks,
			nil,
		},
		{
			"should pass through DAO errors",
			inherentRisks,
			fmt.Errorf("foo bar"),
			nil,
			fmt.Errorf("error retrieving inherentRisks: foo bar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockInherentRiskDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().GetAll(ctx).Return(test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultInherentRiskService(mockDao, nil, nil)
			g, err := service.GetInherentRisks(ctx)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
