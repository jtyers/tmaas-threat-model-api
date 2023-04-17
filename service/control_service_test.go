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

func TestGetControl(t *testing.T) {
	control := m.Control{
		ControlID: m.ControlID("1234-1234-1234-1234"),
		Name:      "my-first-control",
	}

	var tests = []struct {
		name           string
		inputControlID m.ControlID
		daoReturnValue m.Control
		daoReturnError error
		expectedResult *m.Control
		expectedError  error
	}{
		{
			"should get existing controls",
			control.ControlID,
			control,
			nil,
			&control,
			nil,
		},
		{
			"should return ErrNoSuchControl for non-existent controls",
			control.ControlID,
			control,
			servicedao.ErrNoSuchDocument,
			nil,
			ErrNoSuchControl,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockControlDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().Get(ctx, test.inputControlID.String()).Return(&test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultControlService(mockDao, nil, nil)
			g, err := service.GetControl(ctx, test.inputControlID)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpdateControl(t *testing.T) {
	control := m.Control{
		Name: "my-first-control",
		Type: m.ControlTypePreventative,
	}

	var tests = []struct {
		name                      string
		inputID                   m.ControlID
		input                     m.Control
		daoReturnError            error
		validateUpdateReturnError error
		expectedResult            *m.Control
		expectedError             error
	}{
		{
			"should update control",
			control.ControlID,
			control,
			nil,
			nil,
			&control,
			nil,
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			control.ControlID,
			control,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO update fails",
			control.ControlID,
			control,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			fmt.Errorf("error updating control: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockControlDao(ctrl)
			ctx := context.Background()
			mockIDProvider := id.NewMockRandomIDProvider(ctrl)

			mockValidator := validator.NewMockStructValidator(ctrl)
			mockValidator.EXPECT().ValidateForUpdate(test.input).Return(test.validateUpdateReturnError)

			if test.validateUpdateReturnError == nil {
				queryControl := m.Control{ControlID: test.inputID}

				mockDao.EXPECT().UpdateWhereExactSingle(ctx, &queryControl, &test.input).Return(test.expectedResult, test.daoReturnError)

			}

			// when
			service := NewDefaultControlService(mockDao, mockIDProvider, mockValidator)
			err := service.UpdateControl(ctx, test.inputID, test.input)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestDeleteControl(t *testing.T) {
	control := m.Control{
		Name: "my-first-control",
		Type: m.ControlTypePreventative,
	}

	var tests = []struct {
		name           string
		inputID        m.ControlID
		daoReturnError error
		expectedError  error
	}{
		{
			"should delete control",
			control.ControlID,
			nil,
			nil,
		},
		{
			"should fail if DAO delete fails",
			control.ControlID,
			fmt.Errorf("dao failure"),
			fmt.Errorf("error deleting control: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockControlDao(ctrl)
			ctx := context.Background()

			queryControl := m.Control{ControlID: test.inputID}

			mockDao.EXPECT().DeleteWhere(ctx, &queryControl).Return(test.daoReturnError)

			// when
			service := NewDefaultControlService(mockDao, nil, nil)
			err := service.DeleteControl(ctx, test.inputID)

			// then
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateControl(t *testing.T) {
	control := m.Control{
		Name: "my-first-control",
		Type: m.ControlTypePreventative,
	}

	var tests = []struct {
		name                      string
		input                     m.Control
		daoReturnError            error
		validateCreateReturnError error
		validateUpdateReturnError error
		expectedResult            *m.Control
		expectedError             error
	}{
		{
			"should create control",
			control,
			nil,
			nil,
			nil,
			&control,
			nil,
		},
		{
			"should pass through ValidateCreate errors with no wrapping or changes",
			control,
			nil,
			fmt.Errorf("invalid"),
			nil,
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should pass through ValidateForUpdate errors with no wrapping or changes",
			control,
			nil,
			nil,
			fmt.Errorf("invalid"),
			nil,
			fmt.Errorf("invalid"),
		},
		{
			"should fail if DAO create fails",
			control,
			fmt.Errorf("dao failure"),
			nil,
			nil,
			nil,
			fmt.Errorf("error creating control: dao failure"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockControlDao(ctrl)
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
					expectedInputForCreate.ControlID = m.ControlID(ControlIDPrefix + newID)

					mockDao.EXPECT().Create(ctx, &expectedInputForCreate).Return(test.expectedResult, test.daoReturnError)

				}
			}

			// when
			service := NewDefaultControlService(mockDao, mockIDProvider, mockValidator)
			g, err := service.CreateControl(ctx, test.input)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetControls(t *testing.T) {
	controls := []*m.Control{
		{
			ControlID: m.ControlID("1234-1234-1234-1234"),
			Name:      "my-first-control",
			Type:      m.ControlTypePreventative,
		},
		{
			ControlID: m.ControlID("2345-2345-2345-3245"),
			Name:      "my-second-control",
			Type:      m.ControlTypePreventative,
		},
		{
			ControlID: m.ControlID("3456-3456-3456-3456"),
			Name:      "my-third-control",
			Type:      m.ControlTypePreventative,
		},
	}

	var tests = []struct {
		name           string
		daoReturnValue []*m.Control
		daoReturnError error
		expectedResult []*m.Control
		expectedError  error
	}{
		{
			"should get existing controls",
			controls,
			nil,
			controls,
			nil,
		},
		{
			"should pass through DAO errors",
			controls,
			fmt.Errorf("foo bar"),
			nil,
			fmt.Errorf("error retrieving controls: foo bar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDao := dao.NewMockControlDao(ctrl)
			ctx := context.Background()

			mockDao.EXPECT().GetAll(ctx).Return(test.daoReturnValue, test.daoReturnError)

			// when
			service := NewDefaultControlService(mockDao, nil, nil)
			g, err := service.GetControls(ctx)

			// then
			require.Equal(t, test.expectedResult, g)
			require.Equal(t, test.expectedError, err)
		})
	}
}
