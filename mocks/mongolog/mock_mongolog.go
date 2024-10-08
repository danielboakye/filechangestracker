// Code generated by MockGen. DO NOT EDIT.
// Source: mongolog.go

// Package mongologmock is a generated GoMock package.
package mongologmock

import (
	context "context"
	reflect "reflect"

	mongolog "github.com/danielboakye/filechangestracker/internal/mongolog"
	gomock "github.com/golang/mock/gomock"
)

// MockLogStore is a mock of LogStore interface.
type MockLogStore struct {
	ctrl     *gomock.Controller
	recorder *MockLogStoreMockRecorder
}

// MockLogStoreMockRecorder is the mock recorder for MockLogStore.
type MockLogStoreMockRecorder struct {
	mock *MockLogStore
}

// NewMockLogStore creates a new mock instance.
func NewMockLogStore(ctrl *gomock.Controller) *MockLogStore {
	mock := &MockLogStore{ctrl: ctrl}
	mock.recorder = &MockLogStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogStore) EXPECT() *MockLogStoreMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockLogStore) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockLogStoreMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLogStore)(nil).Close), ctx)
}

// ReadLogsPaginated mocks base method.
func (m *MockLogStore) ReadLogsPaginated(ctx context.Context, page, pageSize int64) ([]mongolog.LogEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadLogsPaginated", ctx, page, pageSize)
	ret0, _ := ret[0].([]mongolog.LogEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadLogsPaginated indicates an expected call of ReadLogsPaginated.
func (mr *MockLogStoreMockRecorder) ReadLogsPaginated(ctx, page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadLogsPaginated", reflect.TypeOf((*MockLogStore)(nil).ReadLogsPaginated), ctx, page, pageSize)
}

// Write mocks base method.
func (m *MockLogStore) Write(ctx context.Context, logDetail map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", ctx, logDetail)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockLogStoreMockRecorder) Write(ctx, logDetail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockLogStore)(nil).Write), ctx, logDetail)
}
