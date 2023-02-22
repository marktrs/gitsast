// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/repository.go

// Package testutil is a generated GoMock package.
package testutil

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/marktrs/gitsast/internal/model"
)

// MockIRepositoryRepo is a mock of IRepositoryRepo interface.
type MockIRepositoryRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryRepoMockRecorder
}

// MockIRepositoryRepoMockRecorder is the mock recorder for MockIRepositoryRepo.
type MockIRepositoryRepoMockRecorder struct {
	mock *MockIRepositoryRepo
}

// NewMockIRepositoryRepo creates a new mock instance.
func NewMockIRepositoryRepo(ctrl *gomock.Controller) *MockIRepositoryRepo {
	mock := &MockIRepositoryRepo{ctrl: ctrl}
	mock.recorder = &MockIRepositoryRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRepositoryRepo) EXPECT() *MockIRepositoryRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockIRepositoryRepo) Add(ctx context.Context, repo *model.Repository) (*model.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, repo)
	ret0, _ := ret[0].(*model.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockIRepositoryRepoMockRecorder) Add(ctx, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockIRepositoryRepo)(nil).Add), ctx, repo)
}

// GetById mocks base method.
func (m *MockIRepositoryRepo) GetById(ctx context.Context, id string) (*model.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, id)
	ret0, _ := ret[0].(*model.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockIRepositoryRepoMockRecorder) GetById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockIRepositoryRepo)(nil).GetById), ctx, id)
}

// List mocks base method.
func (m *MockIRepositoryRepo) List(ctx context.Context, f *model.RepositoryFilter) ([]*model.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, f)
	ret0, _ := ret[0].([]*model.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockIRepositoryRepoMockRecorder) List(ctx, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIRepositoryRepo)(nil).List), ctx, f)
}

// Remove mocks base method.
func (m *MockIRepositoryRepo) Remove(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockIRepositoryRepoMockRecorder) Remove(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockIRepositoryRepo)(nil).Remove), ctx, id)
}

// Update mocks base method.
func (m *MockIRepositoryRepo) Update(ctx context.Context, id string, repo map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, repo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIRepositoryRepoMockRecorder) Update(ctx, id, repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIRepositoryRepo)(nil).Update), ctx, id, repo)
}