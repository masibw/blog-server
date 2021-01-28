// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository/posts_tags.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	gomock "github.com/golang/mock/gomock"
	entity "github.com/masibw/blog-server/domain/entity"
	reflect "reflect"
)

// MockPostsTags is a mock of PostsTags interface
type MockPostsTags struct {
	ctrl     *gomock.Controller
	recorder *MockPostsTagsMockRecorder
}

// MockPostsTagsMockRecorder is the mock recorder for MockPostsTags
type MockPostsTagsMockRecorder struct {
	mock *MockPostsTags
}

// NewMockPostsTags creates a new mock instance
func NewMockPostsTags(ctrl *gomock.Controller) *MockPostsTags {
	mock := &MockPostsTags{ctrl: ctrl}
	mock.recorder = &MockPostsTagsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPostsTags) EXPECT() *MockPostsTagsMockRecorder {
	return m.recorder
}

// FindByPostIDAndTagName mocks base method
func (m *MockPostsTags) FindByPostIDAndTagName(postID, tagName string) (*entity.PostsTags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPostIDAndTagName", postID, tagName)
	ret0, _ := ret[0].(*entity.PostsTags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPostIDAndTagName indicates an expected call of FindByPostIDAndTagName
func (mr *MockPostsTagsMockRecorder) FindByPostIDAndTagName(postID, tagName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPostIDAndTagName", reflect.TypeOf((*MockPostsTags)(nil).FindByPostIDAndTagName), postID, tagName)
}

// Store mocks base method
func (m *MockPostsTags) Store(postsTags []*entity.PostsTags) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", postsTags)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockPostsTagsMockRecorder) Store(postsTags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockPostsTags)(nil).Store), postsTags)
}

// Delete mocks base method
func (m *MockPostsTags) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockPostsTagsMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPostsTags)(nil).Delete), id)
}

// DeleteByPostID mocks base method
func (m *MockPostsTags) DeleteByPostID(postID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByPostID", postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByPostID indicates an expected call of DeleteByPostID
func (mr *MockPostsTagsMockRecorder) DeleteByPostID(postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByPostID", reflect.TypeOf((*MockPostsTags)(nil).DeleteByPostID), postID)
}
