// Code generated by MockGen. DO NOT EDIT.
// Source: component_service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/jtyers/tmaas-model"
)

// MockComponentService is a mock of ComponentService interface.
type MockComponentService struct {
	ctrl     *gomock.Controller
	recorder *MockComponentServiceMockRecorder
}

// MockComponentServiceMockRecorder is the mock recorder for MockComponentService.
type MockComponentServiceMockRecorder struct {
	mock *MockComponentService
}

// NewMockComponentService creates a new mock instance.
func NewMockComponentService(ctrl *gomock.Controller) *MockComponentService {
	mock := &MockComponentService{ctrl: ctrl}
	mock.recorder = &MockComponentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComponentService) EXPECT() *MockComponentServiceMockRecorder {
	return m.recorder
}

// CreateComponent mocks base method.
func (m *MockComponentService) CreateComponent(ctx context.Context, component model.Component) (*model.Component, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComponent", ctx, component)
	ret0, _ := ret[0].(*model.Component)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComponent indicates an expected call of CreateComponent.
func (mr *MockComponentServiceMockRecorder) CreateComponent(ctx, component interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComponent", reflect.TypeOf((*MockComponentService)(nil).CreateComponent), ctx, component)
}

// GetComponent mocks base method.
func (m *MockComponentService) GetComponent(ctx context.Context, id model.ComponentID) (*model.Component, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComponent", ctx, id)
	ret0, _ := ret[0].(*model.Component)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComponent indicates an expected call of GetComponent.
func (mr *MockComponentServiceMockRecorder) GetComponent(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComponent", reflect.TypeOf((*MockComponentService)(nil).GetComponent), ctx, id)
}

// GetComponents mocks base method.
func (m *MockComponentService) GetComponents(ctx context.Context, id model.DataFlowDiagramID) ([]*model.Component, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComponents", ctx, id)
	ret0, _ := ret[0].([]*model.Component)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComponents indicates an expected call of GetComponents.
func (mr *MockComponentServiceMockRecorder) GetComponents(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComponents", reflect.TypeOf((*MockComponentService)(nil).GetComponents), ctx, id)
}

// UpdateComponent mocks base method.
func (m *MockComponentService) UpdateComponent(ctx context.Context, componentID model.ComponentID, component model.Component) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComponent", ctx, componentID, component)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateComponent indicates an expected call of UpdateComponent.
func (mr *MockComponentServiceMockRecorder) UpdateComponent(ctx, componentID, component interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComponent", reflect.TypeOf((*MockComponentService)(nil).UpdateComponent), ctx, componentID, component)
}
