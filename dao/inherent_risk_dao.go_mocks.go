// Code generated by MockGen. DO NOT EDIT.
// Source: inherent_risk_dao.go

// Package dao is a generated GoMock package.
package dao

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/jtyers/tmaas-model"
)

// MockInherentRiskDao is a mock of InherentRiskDao interface.
type MockInherentRiskDao struct {
	ctrl     *gomock.Controller
	recorder *MockInherentRiskDaoMockRecorder
}

// MockInherentRiskDaoMockRecorder is the mock recorder for MockInherentRiskDao.
type MockInherentRiskDaoMockRecorder struct {
	mock *MockInherentRiskDao
}

// NewMockInherentRiskDao creates a new mock instance.
func NewMockInherentRiskDao(ctrl *gomock.Controller) *MockInherentRiskDao {
	mock := &MockInherentRiskDao{ctrl: ctrl}
	mock.recorder = &MockInherentRiskDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInherentRiskDao) EXPECT() *MockInherentRiskDaoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockInherentRiskDao) Create(ctx context.Context, data *entities.InherentRisk) (*entities.InherentRisk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].(*entities.InherentRisk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockInherentRiskDaoMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockInherentRiskDao)(nil).Create), ctx, data)
}

// Delete mocks base method.
func (m *MockInherentRiskDao) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockInherentRiskDaoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockInherentRiskDao)(nil).Delete), ctx, id)
}

// DeleteWhere mocks base method.
func (m *MockInherentRiskDao) DeleteWhere(ctx context.Context, query *entities.InherentRisk) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWhere", ctx, query)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWhere indicates an expected call of DeleteWhere.
func (mr *MockInherentRiskDaoMockRecorder) DeleteWhere(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWhere", reflect.TypeOf((*MockInherentRiskDao)(nil).DeleteWhere), ctx, query)
}

// Get mocks base method.
func (m *MockInherentRiskDao) Get(ctx context.Context, id string) (*entities.InherentRisk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entities.InherentRisk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockInherentRiskDaoMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInherentRiskDao)(nil).Get), ctx, id)
}

// GetAll mocks base method.
func (m *MockInherentRiskDao) GetAll(ctx context.Context) ([]*entities.InherentRisk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entities.InherentRisk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockInherentRiskDaoMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockInherentRiskDao)(nil).GetAll), ctx)
}

// QueryExact mocks base method.
func (m *MockInherentRiskDao) QueryExact(ctx context.Context, query *entities.InherentRisk) ([]*entities.InherentRisk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryExact", ctx, query)
	ret0, _ := ret[0].([]*entities.InherentRisk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryExact indicates an expected call of QueryExact.
func (mr *MockInherentRiskDaoMockRecorder) QueryExact(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryExact", reflect.TypeOf((*MockInherentRiskDao)(nil).QueryExact), ctx, query)
}

// QueryExactSingle mocks base method.
func (m *MockInherentRiskDao) QueryExactSingle(ctx context.Context, query *entities.InherentRisk) (*entities.InherentRisk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryExactSingle", ctx, query)
	ret0, _ := ret[0].(*entities.InherentRisk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryExactSingle indicates an expected call of QueryExactSingle.
func (mr *MockInherentRiskDaoMockRecorder) QueryExactSingle(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryExactSingle", reflect.TypeOf((*MockInherentRiskDao)(nil).QueryExactSingle), ctx, query)
}

// UpdateWhereExact mocks base method.
func (m *MockInherentRiskDao) UpdateWhereExact(ctx context.Context, queryExact, data *entities.InherentRisk) ([]*entities.InherentRisk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWhereExact", ctx, queryExact, data)
	ret0, _ := ret[0].([]*entities.InherentRisk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWhereExact indicates an expected call of UpdateWhereExact.
func (mr *MockInherentRiskDaoMockRecorder) UpdateWhereExact(ctx, queryExact, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWhereExact", reflect.TypeOf((*MockInherentRiskDao)(nil).UpdateWhereExact), ctx, queryExact, data)
}

// UpdateWhereExactSingle mocks base method.
func (m *MockInherentRiskDao) UpdateWhereExactSingle(ctx context.Context, queryExact, data *entities.InherentRisk) (*entities.InherentRisk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWhereExactSingle", ctx, queryExact, data)
	ret0, _ := ret[0].(*entities.InherentRisk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWhereExactSingle indicates an expected call of UpdateWhereExactSingle.
func (mr *MockInherentRiskDaoMockRecorder) UpdateWhereExactSingle(ctx, queryExact, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWhereExactSingle", reflect.TypeOf((*MockInherentRiskDao)(nil).UpdateWhereExactSingle), ctx, queryExact, data)
}