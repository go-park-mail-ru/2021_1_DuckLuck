// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category (interfaces: UseCase)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// GetCatalogCategories mocks base method.
func (m *MockUseCase) GetCatalogCategories() ([]*models.CategoriesCatalog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCatalogCategories")
	ret0, _ := ret[0].([]*models.CategoriesCatalog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCatalogCategories indicates an expected call of GetCatalogCategories.
func (mr *MockUseCaseMockRecorder) GetCatalogCategories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCatalogCategories", reflect.TypeOf((*MockUseCase)(nil).GetCatalogCategories))
}

// GetSubCategoriesById mocks base method.
func (m *MockUseCase) GetSubCategoriesById(arg0 uint64) ([]*models.CategoriesCatalog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubCategoriesById", arg0)
	ret0, _ := ret[0].([]*models.CategoriesCatalog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubCategoriesById indicates an expected call of GetSubCategoriesById.
func (mr *MockUseCaseMockRecorder) GetSubCategoriesById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubCategoriesById", reflect.TypeOf((*MockUseCase)(nil).GetSubCategoriesById), arg0)
}
