// Code generated by MockGen. DO NOT EDIT.
// Source: osquerymanager.go

// Package osquerymanagermock is a generated GoMock package.
package osquerymanagermock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOSQueryManager is a mock of OSQueryManager interface.
type MockOSQueryManager struct {
	ctrl     *gomock.Controller
	recorder *MockOSQueryManagerMockRecorder
}

// MockOSQueryManagerMockRecorder is the mock recorder for MockOSQueryManager.
type MockOSQueryManagerMockRecorder struct {
	mock *MockOSQueryManager
}

// NewMockOSQueryManager creates a new mock instance.
func NewMockOSQueryManager(ctrl *gomock.Controller) *MockOSQueryManager {
	mock := &MockOSQueryManager{ctrl: ctrl}
	mock.recorder = &MockOSQueryManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOSQueryManager) EXPECT() *MockOSQueryManagerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockOSQueryManager) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockOSQueryManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockOSQueryManager)(nil).Close))
}

// Query mocks base method.
func (m *MockOSQueryManager) Query(sql string) ([]map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", sql)
	ret0, _ := ret[0].([]map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockOSQueryManagerMockRecorder) Query(sql interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockOSQueryManager)(nil).Query), sql)
}
