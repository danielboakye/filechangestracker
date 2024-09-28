// Code generated by MockGen. DO NOT EDIT.
// Source: filechangestracker.go

// Package filechangestrackermock is a generated GoMock package.
package filechangestrackermock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileChangesTracker is a mock of FileChangesTracker interface.
type MockFileChangesTracker struct {
	ctrl     *gomock.Controller
	recorder *MockFileChangesTrackerMockRecorder
}

// MockFileChangesTrackerMockRecorder is the mock recorder for MockFileChangesTracker.
type MockFileChangesTrackerMockRecorder struct {
	mock *MockFileChangesTracker
}

// NewMockFileChangesTracker creates a new mock instance.
func NewMockFileChangesTracker(ctrl *gomock.Controller) *MockFileChangesTracker {
	mock := &MockFileChangesTracker{ctrl: ctrl}
	mock.recorder = &MockFileChangesTrackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileChangesTracker) EXPECT() *MockFileChangesTrackerMockRecorder {
	return m.recorder
}

// GetLogs mocks base method.
func (m *MockFileChangesTracker) GetLogs() ([]map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs")
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs.
func (mr *MockFileChangesTrackerMockRecorder) GetLogs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockFileChangesTracker)(nil).GetLogs))
}

// IsTimerThreadAlive mocks base method.
func (m *MockFileChangesTracker) IsTimerThreadAlive() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTimerThreadAlive")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTimerThreadAlive indicates an expected call of IsTimerThreadAlive.
func (mr *MockFileChangesTrackerMockRecorder) IsTimerThreadAlive() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTimerThreadAlive", reflect.TypeOf((*MockFileChangesTracker)(nil).IsTimerThreadAlive))
}

// Start mocks base method.
func (m *MockFileChangesTracker) Start(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockFileChangesTrackerMockRecorder) Start(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockFileChangesTracker)(nil).Start), ctx)
}

// Stop mocks base method.
func (m *MockFileChangesTracker) Stop(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockFileChangesTrackerMockRecorder) Stop(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockFileChangesTracker)(nil).Stop), ctx)
}
