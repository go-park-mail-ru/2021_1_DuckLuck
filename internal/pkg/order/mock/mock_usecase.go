// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order (interfaces: UseCase)

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

// AddCompletedOrder mocks base method.
func (m *MockUseCase) AddCompletedOrder(arg0 *models.Order, arg1 uint64, arg2 *models.PreviewCart) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCompletedOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCompletedOrder indicates an expected call of AddCompletedOrder.
func (mr *MockUseCaseMockRecorder) AddCompletedOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCompletedOrder", reflect.TypeOf((*MockUseCase)(nil).AddCompletedOrder), arg0, arg1, arg2)
}

// GetPreviewOrder mocks base method.
func (m *MockUseCase) GetPreviewOrder(arg0 uint64, arg1 *models.PreviewCart) (*models.PreviewOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreviewOrder", arg0, arg1)
	ret0, _ := ret[0].(*models.PreviewOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPreviewOrder indicates an expected call of GetPreviewOrder.
func (mr *MockUseCaseMockRecorder) GetPreviewOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreviewOrder", reflect.TypeOf((*MockUseCase)(nil).GetPreviewOrder), arg0, arg1)
}
