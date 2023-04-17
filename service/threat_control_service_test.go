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

func TestGetThreatControl(t *testing.T) {
	threatControl := m.ThreatControl{
		ThreatControlID: m.ThreatControlID("1234-1234-1234-1234"),
		ControlID:       m.ControlID("xxxxyyyyzzzz"),
	}

	var tests = []struct {
		name                 string
		inputThreatControlID m.ThreatControlID
		daoReturnValue       m.ThreatControl
		daoReturnError       error
		expectedResult       *m.ThreatControl
		expectedError        error
	}{
		{
			"should get existing threatControls",
			threatControl.ThreatControlID,
			threatControl,
			nil,
			&threatControl,
			nil,
		},
		{
			"should return ErrNoSuchThreatControl for non-existent threatControls",
			threatControl.ThreatControlID,
			threatControl,
			servicedao.ErrNoSuchDocument,
			nil,
			ErrNoSuchThreatControl,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatControlDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().Get(ctx, test.inputThreatControlID.String()).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultThreatControlService(mockDao, nil, nil)
			g, err := service.GetThreatControl(ctx, test.inputThreatControlID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateThreatControl(t *testing.T) {
	threatControl := m.ThreatControl{
		ControlID: m.ControlID("xxxxyyyyzzzz"),
	}

	var tests = []struct {
		name                      string
		inputID                   m.ThreatControlID
		input                     m.ThreatControl
		daoReturnError            error
		validateUpdateReturnError error
		expectedResult            *m.ThreatControl
		expectedError             error
	}{
		{
			"should update threatControl",
			threatControl.ThreatControlID,
			threatControl,
			nil,
			nil,
			&threatControl,
			nil,
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			threatControl.ThreatControlID,
			threatControl,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			threatControl.ThreatControlID,
			threatControl,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			fmt.Errorf("error updating threatControl: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatControlDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				queryThreatControl := m.ThreatControl{ThreatControlID: test.inputID}

				mockDao.EXPECT().UpdateWhereExactSingle(ctx, &queryThreatControl, &test.input).Return(test.expectedResult, test.daoReturnError)

			}

			// when
			service := NewDefaultThreatControlService(mockDao, mockIDProvider, mockValidator)
			err := service.UpdateThreatControl(ctx, test.inputID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateThreatControl(t *testing.T) {
	threatControl := m.ThreatControl{
		ControlID: m.ControlID("xxxxyyyyzzzz"),
	}

	var tests = []struct {
		name                      string
		input                     m.ThreatControl
		daoReturnError            error
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.ThreatControl
		expectedError             error
	}{
		{
			"should create threatControl",
			threatControl,
			nil,
			nil,
			nil,
			&threatControl,
			nil,
		},
		{
			"should pass through ValidateCreate errors with no wrapping or changes",
			threatControl,
			nil,
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			threatControl,
			nil,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO create fails",
			threatControl,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			nil,
			fmt.Errorf("error creating threatControl: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatControlDao(ctrl)
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
					expectedInputForCreate.ThreatControlID = m.ThreatControlID(ThreatControlIDPrefix + newID)

					mockDao.EXPECT().Create(ctx, &expectedInputForCreate).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultThreatControlService(mockDao, mockIDProvider, mockValidator)
			g, err := service.CreateThreatControl(ctx, test.input)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetThreatControls(t *testing.T) {
	threatControls := []*m.ThreatControl{
		{
			ThreatControlID: m.ThreatControlID("1234-1234-1234-1234"),
			ControlID:       m.ControlID("zzzzyyyyxxxx"),
		},
		{
			ThreatControlID: m.ThreatControlID("2345-2345-2345-3245"),
			ControlID:       m.ControlID("xxxxyyyyzzzz"),
		},
		{
			ThreatControlID: m.ThreatControlID("3456-3456-3456-3456"),
			ControlID:       m.ControlID("aaaabbbbcccc"),
		},
	}

	var tests = []struct {
		name           string
		daoReturnValue []*m.ThreatControl
		daoReturnError error
		expectedResult []*m.ThreatControl
		expectedError  error
	}{
		{
			"should get existing threatControls",
			threatControls,
			nil,
			threatControls,
			nil,
		},
		{
			"should pass through DAO errors",
			threatControls,
			fmt.Errorf("foo bar"),
			nil,
			fmt.Errorf("error retrieving threatControls: foo bar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockThreatControlDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().GetAll(ctx).Return(test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultThreatControlService(mockDao, nil, nil)
			g, err := service.GetThreatControls(ctx)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
