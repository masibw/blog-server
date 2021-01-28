// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository/post.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	gomock "github.com/golang/mock/gomock"
	entity "github.com/masibw/blog-server/domain/entity"
	reflect "reflect"
)

// MockPost is a mock of Post interface
type MockPost struct {
	ctrl     *gomock.Controller
	recorder *MockPostMockRecorder
}

// MockPostMockRecorder is the mock recorder for MockPost
type MockPostMockRecorder struct {
	mock *MockPost
}

// NewMockPost creates a new mock instance
func NewMockPost(ctrl *gomock.Controller) *MockPost {
	mock := &MockPost{ctrl: ctrl}
	mock.recorder = &MockPostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPost) EXPECT() *MockPostMockRecorder {
	return m.recorder
}

// FindByID mocks base method
func (m *MockPost) FindByID(id string) (*entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockPostMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockPost)(nil).FindByID), id)
}

// FindAll mocks base method
func (m *MockPost) FindAll(offset, pageSize int, conditions string, params []interface{}) ([]*entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", offset, pageSize, conditions, params)
	ret0, _ := ret[0].([]*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll
func (mr *MockPostMockRecorder) FindAll(offset, pageSize, conditions, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockPost)(nil).FindAll), offset, pageSize, conditions, params)
}

// FindByPermalink mocks base method
func (m *MockPost) FindByPermalink(permalink string) (*entity.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPermalink", permalink)
	ret0, _ := ret[0].(*entity.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPermalink indicates an expected call of FindByPermalink
func (mr *MockPostMockRecorder) FindByPermalink(permalink interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPermalink", reflect.TypeOf((*MockPost)(nil).FindByPermalink), permalink)
}

// Create mocks base method
func (m *MockPost) Create(post *entity.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockPostMockRecorder) Create(post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPost)(nil).Create), post)
}

// Update mocks base method
func (m *MockPost) Update(post *entity.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockPostMockRecorder) Update(post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPost)(nil).Update), post)
}

// Delete mocks base method
func (m *MockPost) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockPostMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPost)(nil).Delete), id)
}
