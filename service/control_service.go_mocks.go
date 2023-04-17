// Code generated by MockGen. DO NOT EDIT.
// Source: control_service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/jtyers/tmaas-model"
)

// MockControlService is a mock of ControlService interface.
type MockControlService struct {
	ctrl     *gomock.Controller
	recorder *MockControlServiceMockRecorder
}

// MockControlServiceMockRecorder is the mock recorder for MockControlService.
type MockControlServiceMockRecorder struct {
	mock *MockControlService
}

// NewMockControlService creates a new mock instance.
func NewMockControlService(ctrl *gomock.Controller) *MockControlService {
	mock := &MockControlService{ctrl: ctrl}
	mock.recorder = &MockControlServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControlService) EXPECT() *MockControlServiceMockRecorder {
	return m.recorder
}

// CreateControl mocks base method.
func (m *MockControlService) CreateControl(ctx context.Context, control model.Control) (*model.Control, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateControl", ctx, control)
	ret0, _ := ret[0].(*model.Control)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateControl indicates an expected call of CreateControl.
func (mr *MockControlServiceMockRecorder) CreateControl(ctx, control interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateControl", reflect.TypeOf((*MockControlService)(nil).CreateControl), ctx, control)
}

// DeleteControl mocks base method.
func (m *MockControlService) DeleteControl(ctx context.Context, controlID model.ControlID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteControl", ctx, controlID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteControl indicates an expected call of DeleteControl.
func (mr *MockControlServiceMockRecorder) DeleteControl(ctx, controlID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteControl", reflect.TypeOf((*MockControlService)(nil).DeleteControl), ctx, controlID)
}

// GetControl mocks base method.
func (m *MockControlService) GetControl(ctx context.Context, id model.ControlID) (*model.Control, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetControl", ctx, id)
	ret0, _ := ret[0].(*model.Control)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetControl indicates an expected call of GetControl.
func (mr *MockControlServiceMockRecorder) GetControl(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetControl", reflect.TypeOf((*MockControlService)(nil).GetControl), ctx, id)
}

// GetControls mocks base method.
func (m *MockControlService) GetControls(ctx context.Context) ([]*model.Control, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetControls", ctx)
	ret0, _ := ret[0].([]*model.Control)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetControls indicates an expected call of GetControls.
func (mr *MockControlServiceMockRecorder) GetControls(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetControls", reflect.TypeOf((*MockControlService)(nil).GetControls), ctx)
}

// UpdateControl mocks base method.
func (m *MockControlService) UpdateControl(ctx context.Context, controlID model.ControlID, control model.Control) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateControl", ctx, controlID, control)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateControl indicates an expected call of UpdateControl.
func (mr *MockControlServiceMockRecorder) UpdateControl(ctx, controlID, control interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateControl", reflect.TypeOf((*MockControlService)(nil).UpdateControl), ctx, controlID, control)
}
