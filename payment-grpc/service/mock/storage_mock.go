// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	paymentpb "github.com/Edbeer/payment-grpc/proto"
	types "github.com/Edbeer/payment-grpc/types"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetPaymentByID mocks base method.
func (m *MockStorage) GetPaymentByID(ctx context.Context, req *paymentpb.PaidRequest) (*types.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentByID", ctx, req)
	ret0, _ := ret[0].(*types.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaymentByID indicates an expected call of GetPaymentByID.
func (mr *MockStorageMockRecorder) GetPaymentByID(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentByID", reflect.TypeOf((*MockStorage)(nil).GetPaymentByID), ctx, req)
}

// SavePayment mocks base method.
func (m *MockStorage) SavePayment(ctx context.Context, tx *sql.Tx, payment *types.Payment) (*types.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePayment", ctx, tx, payment)
	ret0, _ := ret[0].(*types.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SavePayment indicates an expected call of SavePayment.
func (mr *MockStorageMockRecorder) SavePayment(ctx, tx, payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePayment", reflect.TypeOf((*MockStorage)(nil).SavePayment), ctx, tx, payment)
}
