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

func TestGetThreat(t *testing.T) {
	threat := m.Threat{
		ThreatID:    m.ThreatID("1234-1234-1234-1234"),
		Description: "my-first-threat",
	}

	var tests = []struct {
		name           string
		inputThreatID  m.ThreatID
		daoReturnValue m.Threat
		daoReturnError error
		expectedResult *m.Threat
		expectedError  error
	}{
		{
			"should get existing threats",
			threat.ThreatID,
			threat,
			nil,
			&threat,
			nil,
		},
		{
			"should return ErrNoSuchThreat for non-existent threats",
			threat.ThreatID,
			threat,
			servicedao.ErrNoSuchDocument,
			nil,
			ErrNoSuchThreat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().Get(ctx, test.inputThreatID.String()).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultThreatService(mockDao, nil, nil)
			g, err := service.GetThreat(ctx, test.inputThreatID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateThreat(t *testing.T) {
	threat := m.Threat{
		Description: "my-first-threat",
		InID:        "cm_123456",
		Stride:      m.StrideSpoofing,
	}

	var tests = []struct {
		name                      string
		inputID                   m.ThreatID
		input                     m.Threat
		daoReturnError            error
		validateUpdateReturnError error
		expectedResult            *m.Threat
		expectedError             error
	}{
		{
			"should update threat",
			threat.ThreatID,
			threat,
			nil,
			nil,
			&threat,
			nil,
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			threat.ThreatID,
			threat,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			threat.ThreatID,
			threat,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			fmt.Errorf("error updating threat: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				queryThreat := m.Threat{ThreatID: test.inputID}

				mockDao.EXPECT().UpdateWhereExactSingle(ctx, &queryThreat, &test.input).Return(test.expectedResult, test.daoReturnError)

			}

			// when
			service := NewDefaultThreatService(mockDao, mockIDProvider, mockValidator)
			err := service.UpdateThreat(ctx, test.inputID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestDeleteThreat(t *testing.T) {
	threat := m.Threat{
		Description: "my-first-threat",
		InID:        "cm_123456",
		Stride:      m.StrideSpoofing,
	}

	var tests = []struct {
		name           string
		inputID        m.ThreatID
		daoReturnError error
		expectedError  error
	}{
		{
			"should delete threat",
			threat.ThreatID,
			nil,
			nil,
		},
		{
			"should fail if DAO delete fails",
			threat.ThreatID,
			fmt.Errorf("dao failure"),
			fmt.Errorf("error deleting threat: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatDao(ctrl)
			ctx := context.Background()

			queryThreat := m.Threat{ThreatID: test.inputID}

			mockDao.EXPECT().DeleteWhere(ctx, &queryThreat).Return(test.daoReturnError)

			// when
			service := NewDefaultThreatService(mockDao, nil, nil)
			err := service.DeleteThreat(ctx, test.inputID)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateThreat(t *testing.T) {
	threat := m.Threat{
		Description: "my-first-threat",
		InID:        "cm_123456",
		Stride:      m.StrideSpoofing,
	}

	var tests = []struct {
		name                      string
		input                     m.Threat
		daoReturnError            error
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.Threat
		expectedError             error
	}{
		{
			"should create threat",
			threat,
			nil,
			nil,
			nil,
			&threat,
			nil,
		},
		{
			"should pass through ValidateCreate errors with no wrapping or changes",
			threat,
			nil,
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			threat,
			nil,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO create fails",
			threat,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			nil,
			fmt.Errorf("error creating threat: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatDao(ctrl)
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
					expectedInputForCreate.ThreatID = m.ThreatID(ThreatIDPrefix + newID)

					mockDao.EXPECT().Create(ctx, &expectedInputForCreate).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultThreatService(mockDao, mockIDProvider, mockValidator)
			g, err := service.CreateThreat(ctx, test.input)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetThreats(t *testing.T) {
	threats := []*m.Threat{
		{
			ThreatID:    m.ThreatID("1234-1234-1234-1234"),
			Description: "my-first-threat",
		},
		{
			ThreatID:    m.ThreatID("2345-2345-2345-3245"),
			Description: "my-second-threat",
		},
		{
			ThreatID:    m.ThreatID("3456-3456-3456-3456"),
			Description: "my-third-threat",
		},
	}

	var tests = []struct {
		name           string
		daoReturnValue []*m.Threat
		daoReturnError error
		expectedResult []*m.Threat
		expectedError  error
	}{
		{
			"should get existing threats",
			threats,
			nil,
			threats,
			nil,
		},
		{
			"should pass through DAO errors",
			threats,
			fmt.Errorf("foo bar"),
			nil,
			fmt.Errorf("error retrieving threats: foo bar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().GetAll(ctx).Return(test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultThreatService(mockDao, nil, nil)
			g, err := service.GetThreats(ctx)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
