// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/jtyers/tmaas-model"
)

// MockThreatModelService is a mock of ThreatModelService interface.
type MockThreatModelService struct {
	ctrl     *gomock.Controller
	recorder *MockThreatModelServiceMockRecorder
}

// MockThreatModelServiceMockRecorder is the mock recorder for MockThreatModelService.
type MockThreatModelServiceMockRecorder struct {
	mock *MockThreatModelService
}

// NewMockThreatModelService creates a new mock instance.
func NewMockThreatModelService(ctrl *gomock.Controller) *MockThreatModelService {
	mock := &MockThreatModelService{ctrl: ctrl}
	mock.recorder = &MockThreatModelServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockThreatModelService) EXPECT() *MockThreatModelServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockThreatModelService) Create(ctx context.Context, params model.ThreatModelParams) (*model.ThreatModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(*model.ThreatModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockThreatModelServiceMockRecorder) Create(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockThreatModelService)(nil).Create), ctx, params)
}

// Delete mocks base method.
func (m *MockThreatModelService) Delete(ctx context.Context, id model.ThreatModelID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockThreatModelServiceMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockThreatModelService)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockThreatModelService) Get(ctx context.Context, id model.ThreatModelID) (*model.ThreatModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*model.ThreatModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockThreatModelServiceMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockThreatModelService)(nil).Get), ctx, id)
}

// GetAll mocks base method.
func (m *MockThreatModelService) GetAll(ctx context.Context) ([]*model.ThreatModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*model.ThreatModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockThreatModelServiceMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockThreatModelService)(nil).GetAll), ctx)
}

// Query mocks base method.
func (m *MockThreatModelService) Query(ctx context.Context, q *model.ThreatModelQuery) ([]*model.ThreatModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", ctx, q)
	ret0, _ := ret[0].([]*model.ThreatModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockThreatModelServiceMockRecorder) Query(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockThreatModelService)(nil).Query), ctx, q)
}

// QuerySingle mocks base method.
func (m *MockThreatModelService) QuerySingle(ctx context.Context, q *model.ThreatModelQuery) (*model.ThreatModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QuerySingle", ctx, q)
	ret0, _ := ret[0].(*model.ThreatModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QuerySingle indicates an expected call of QuerySingle.
func (mr *MockThreatModelServiceMockRecorder) QuerySingle(ctx, q interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QuerySingle", reflect.TypeOf((*MockThreatModelService)(nil).QuerySingle), ctx, q)
}

// Update mocks base method.
func (m *MockThreatModelService) Update(ctx context.Context, id model.ThreatModelID, params model.ThreatModelParams) (*model.ThreatModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, params)
	ret0, _ := ret[0].(*model.ThreatModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockThreatModelServiceMockRecorder) Update(ctx, id, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockThreatModelService)(nil).Update), ctx, id, params)
}
