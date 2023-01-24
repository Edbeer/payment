// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Edbeer/auth-grpc/proto (interfaces: AuthServiceClient)

// Package mock_proto is a generated GoMock package.
package mock_proto

import (
	context "context"
	reflect "reflect"

	authpb "github.com/Edbeer/auth-grpc/proto"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAuthServiceClient is a mock of AuthServiceClient interface.
type MockAuthServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceClientMockRecorder
}

// MockAuthServiceClientMockRecorder is the mock recorder for MockAuthServiceClient.
type MockAuthServiceClientMockRecorder struct {
	mock *MockAuthServiceClient
}

// NewMockAuthServiceClient creates a new mock instance.
func NewMockAuthServiceClient(ctrl *gomock.Controller) *MockAuthServiceClient {
	mock := &MockAuthServiceClient{ctrl: ctrl}
	mock.recorder = &MockAuthServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceClient) EXPECT() *MockAuthServiceClientMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockAuthServiceClient) CreateAccount(arg0 context.Context, arg1 *authpb.CreateRequest, arg2 ...grpc.CallOption) (*authpb.AccountWithToken, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAccount", varargs...)
	ret0, _ := ret[0].(*authpb.AccountWithToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockAuthServiceClientMockRecorder) CreateAccount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAuthServiceClient)(nil).CreateAccount), varargs...)
}

// CreateStatement mocks base method.
func (m *MockAuthServiceClient) CreateStatement(arg0 context.Context, arg1 ...grpc.CallOption) (authpb.AuthService_CreateStatementClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateStatement", varargs...)
	ret0, _ := ret[0].(authpb.AuthService_CreateStatementClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStatement indicates an expected call of CreateStatement.
func (mr *MockAuthServiceClientMockRecorder) CreateStatement(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStatement", reflect.TypeOf((*MockAuthServiceClient)(nil).CreateStatement), varargs...)
}

// DeleteAccount mocks base method.
func (m *MockAuthServiceClient) DeleteAccount(arg0 context.Context, arg1 *authpb.DeleteRequest, arg2 ...grpc.CallOption) (*authpb.DeleteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAccount", varargs...)
	ret0, _ := ret[0].(*authpb.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockAuthServiceClientMockRecorder) DeleteAccount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockAuthServiceClient)(nil).DeleteAccount), varargs...)
}

// DepositAccount mocks base method.
func (m *MockAuthServiceClient) DepositAccount(arg0 context.Context, arg1 *authpb.DepositRequest, arg2 ...grpc.CallOption) (*authpb.DepositResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DepositAccount", varargs...)
	ret0, _ := ret[0].(*authpb.DepositResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DepositAccount indicates an expected call of DepositAccount.
func (mr *MockAuthServiceClientMockRecorder) DepositAccount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepositAccount", reflect.TypeOf((*MockAuthServiceClient)(nil).DepositAccount), varargs...)
}

// GetAccount mocks base method.
func (m *MockAuthServiceClient) GetAccount(arg0 context.Context, arg1 *authpb.GetRequest, arg2 ...grpc.CallOption) (authpb.AuthService_GetAccountClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccount", varargs...)
	ret0, _ := ret[0].(authpb.AuthService_GetAccountClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockAuthServiceClientMockRecorder) GetAccount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAuthServiceClient)(nil).GetAccount), varargs...)
}

// GetAccountByID mocks base method.
func (m *MockAuthServiceClient) GetAccountByID(arg0 context.Context, arg1 *authpb.GetIDRequest, arg2 ...grpc.CallOption) (*authpb.Account, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccountByID", varargs...)
	ret0, _ := ret[0].(*authpb.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountByID indicates an expected call of GetAccountByID.
func (mr *MockAuthServiceClientMockRecorder) GetAccountByID(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountByID", reflect.TypeOf((*MockAuthServiceClient)(nil).GetAccountByID), varargs...)
}

// GetStatement mocks base method.
func (m *MockAuthServiceClient) GetStatement(arg0 context.Context, arg1 *authpb.StatementGet, arg2 ...grpc.CallOption) (authpb.AuthService_GetStatementClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetStatement", varargs...)
	ret0, _ := ret[0].(authpb.AuthService_GetStatementClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatement indicates an expected call of GetStatement.
func (mr *MockAuthServiceClientMockRecorder) GetStatement(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatement", reflect.TypeOf((*MockAuthServiceClient)(nil).GetStatement), varargs...)
}

// UpdateAccount mocks base method.
func (m *MockAuthServiceClient) UpdateAccount(arg0 context.Context, arg1 *authpb.UpdateRequest, arg2 ...grpc.CallOption) (*authpb.Account, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateAccount", varargs...)
	ret0, _ := ret[0].(*authpb.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockAuthServiceClientMockRecorder) UpdateAccount(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockAuthServiceClient)(nil).UpdateAccount), varargs...)
}

// UpdateBalance mocks base method.
func (m *MockAuthServiceClient) UpdateBalance(arg0 context.Context, arg1 *authpb.UpdateBalanceRequest, arg2 ...grpc.CallOption) (*authpb.Account, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateBalance", varargs...)
	ret0, _ := ret[0].(*authpb.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockAuthServiceClientMockRecorder) UpdateBalance(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockAuthServiceClient)(nil).UpdateBalance), varargs...)
}
