// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository/tag.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/masibw/blog-server/domain/entity"
)

// MockTag is a mock of Tag interface
type MockTag struct {
	ctrl     *gomock.Controller
	recorder *MockTagMockRecorder
}

// MockTagMockRecorder is the mock recorder for MockTag
type MockTagMockRecorder struct {
	mock *MockTag
}

// NewMockTag creates a new mock instance
func NewMockTag(ctrl *gomock.Controller) *MockTag {
	mock := &MockTag{ctrl: ctrl}
	mock.recorder = &MockTagMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTag) EXPECT() *MockTagMockRecorder {
	return m.recorder
}

// FindByID mocks base method
func (m *MockTag) FindByID(id string) (*entity.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*entity.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockTagMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockTag)(nil).FindByID), id)
}

// FindAll mocks base method
func (m *MockTag) FindAll(offset, pageSize int) ([]*entity.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", offset, pageSize)
	ret0, _ := ret[0].([]*entity.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll
func (mr *MockTagMockRecorder) FindAll(offset, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockTag)(nil).FindAll), offset, pageSize)
}

// FindByName mocks base method
func (m *MockTag) FindByName(name string) (*entity.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByName", name)
	ret0, _ := ret[0].(*entity.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByName indicates an expected call of FindByName
func (mr *MockTagMockRecorder) FindByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByName", reflect.TypeOf((*MockTag)(nil).FindByName), name)
}

// Store mocks base method
func (m *MockTag) Store(post *entity.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockTagMockRecorder) Store(post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTag)(nil).Store), post)
}

// Delete mocks base method
func (m *MockTag) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTagMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTag)(nil).Delete), id)
}
