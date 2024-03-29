// Code generated by MockGen. DO NOT EDIT.
// Source: filesystem.go

// Package mock_filesystem is a generated GoMock package.
package mock_filesystem

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileSystem is a mock of FileSystem interface.
type MockFileSystem struct {
	ctrl     *gomock.Controller
	recorder *MockFileSystemMockRecorder
}

// MockFileSystemMockRecorder is the mock recorder for MockFileSystem.
type MockFileSystemMockRecorder struct {
	mock *MockFileSystem
}

// NewMockFileSystem creates a new mock instance.
func NewMockFileSystem(ctrl *gomock.Controller) *MockFileSystem {
	mock := &MockFileSystem{ctrl: ctrl}
	mock.recorder = &MockFileSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileSystem) EXPECT() *MockFileSystemMockRecorder {
	return m.recorder
}

// GetAvailableDiskSpace mocks base method.
func (m *MockFileSystem) GetAvailableDiskSpace(path string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableDiskSpace", path)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableDiskSpace indicates an expected call of GetAvailableDiskSpace.
func (mr *MockFileSystemMockRecorder) GetAvailableDiskSpace(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableDiskSpace", reflect.TypeOf((*MockFileSystem)(nil).GetAvailableDiskSpace), path)
}
