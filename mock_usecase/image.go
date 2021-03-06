// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/image.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockImage is a mock of Image interface
type MockImage struct {
	ctrl     *gomock.Controller
	recorder *MockImageMockRecorder
}

// MockImageMockRecorder is the mock recorder for MockImage
type MockImageMockRecorder struct {
	mock *MockImage
}

// NewMockImage creates a new mock instance
func NewMockImage(ctrl *gomock.Controller) *MockImage {
	mock := &MockImage{ctrl: ctrl}
	mock.recorder = &MockImageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImage) EXPECT() *MockImageMockRecorder {
	return m.recorder
}

// CreatePresignedURL mocks base method
func (m *MockImage) CreatePresignedURL(fileName, contentType *string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePresignedURL", fileName, contentType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePresignedURL indicates an expected call of CreatePresignedURL
func (mr *MockImageMockRecorder) CreatePresignedURL(fileName, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePresignedURL", reflect.TypeOf((*MockImage)(nil).CreatePresignedURL), fileName, contentType)
}
