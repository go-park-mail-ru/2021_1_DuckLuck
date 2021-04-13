// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order (interfaces: Handler)

// Package mock is a generated GoMock package.
package mock

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// AddCompletedOrder mocks base method.
func (m *MockHandler) AddCompletedOrder(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddCompletedOrder", arg0, arg1)
}

// AddCompletedOrder indicates an expected call of AddCompletedOrder.
func (mr *MockHandlerMockRecorder) AddCompletedOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCompletedOrder", reflect.TypeOf((*MockHandler)(nil).AddCompletedOrder), arg0, arg1)
}

// GetOrderFromCart mocks base method.
func (m *MockHandler) GetOrderFromCart(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetOrderFromCart", arg0, arg1)
}

// GetOrderFromCart indicates an expected call of GetOrderFromCart.
func (mr *MockHandlerMockRecorder) GetOrderFromCart(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderFromCart", reflect.TypeOf((*MockHandler)(nil).GetOrderFromCart), arg0, arg1)
}