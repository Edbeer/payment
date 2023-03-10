// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "github.com/Edbeer/auth-grpc/types"
	proto "github.com/Edbeer/payment-proto/auth-grpc/proto"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
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

// CreateAccount mocks base method.
func (m *MockStorage) CreateAccount(ctx context.Context, account *proto.CreateRequest) (*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, account)
	ret0, _ := ret[0].(*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockStorageMockRecorder) CreateAccount(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockStorage)(nil).CreateAccount), ctx, account)
}

// DeleteAccount mocks base method.
func (m *MockStorage) DeleteAccount(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", ctx, req)
	ret0, _ := ret[0].(*proto.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockStorageMockRecorder) DeleteAccount(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockStorage)(nil).DeleteAccount), ctx, req)
}

// DepositAccount mocks base method.
func (m *MockStorage) DepositAccount(ctx context.Context, req *proto.DepositRequest) (*proto.DepositResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepositAccount", ctx, req)
	ret0, _ := ret[0].(*proto.DepositResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DepositAccount indicates an expected call of DepositAccount.
func (mr *MockStorageMockRecorder) DepositAccount(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepositAccount", reflect.TypeOf((*MockStorage)(nil).DepositAccount), ctx, req)
}

// GetAccount mocks base method.
func (m *MockStorage) GetAccount(ctx context.Context) ([]*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx)
	ret0, _ := ret[0].([]*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockStorageMockRecorder) GetAccount(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockStorage)(nil).GetAccount), ctx)
}

// GetAccountByID mocks base method.
func (m *MockStorage) GetAccountByID(ctx context.Context, req *proto.GetIDRequest) (*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountByID", ctx, req)
	ret0, _ := ret[0].(*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountByID indicates an expected call of GetAccountByID.
func (mr *MockStorageMockRecorder) GetAccountByID(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountByID", reflect.TypeOf((*MockStorage)(nil).GetAccountByID), ctx, req)
}

// SaveBalance mocks base method.
func (m *MockStorage) SaveBalance(ctx context.Context, req *proto.UpdateBalanceRequest) (*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBalance", ctx, req)
	ret0, _ := ret[0].(*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveBalance indicates an expected call of SaveBalance.
func (mr *MockStorageMockRecorder) SaveBalance(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBalance", reflect.TypeOf((*MockStorage)(nil).SaveBalance), ctx, req)
}

// UpdateAccount mocks base method.
func (m *MockStorage) UpdateAccount(ctx context.Context, account *proto.UpdateRequest) (*types.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", ctx, account)
	ret0, _ := ret[0].(*types.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockStorageMockRecorder) UpdateAccount(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockStorage)(nil).UpdateAccount), ctx, account)
}

// UpdateStatement mocks base method.
func (m *MockStorage) UpdateStatement(ctx context.Context, req *proto.StatementRequest) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatement", ctx, req)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatement indicates an expected call of UpdateStatement.
func (mr *MockStorageMockRecorder) UpdateStatement(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatement", reflect.TypeOf((*MockStorage)(nil).UpdateStatement), ctx, req)
}

// MockRedisStorage is a mock of RedisStorage interface.
type MockRedisStorage struct {
	ctrl     *gomock.Controller
	recorder *MockRedisStorageMockRecorder
}

// MockRedisStorageMockRecorder is the mock recorder for MockRedisStorage.
type MockRedisStorageMockRecorder struct {
	mock *MockRedisStorage
}

// NewMockRedisStorage creates a new mock instance.
func NewMockRedisStorage(ctrl *gomock.Controller) *MockRedisStorage {
	mock := &MockRedisStorage{ctrl: ctrl}
	mock.recorder = &MockRedisStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisStorage) EXPECT() *MockRedisStorageMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockRedisStorage) CreateSession(ctx context.Context, session *types.Session, expire int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, session, expire)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockRedisStorageMockRecorder) CreateSession(ctx, session, expire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockRedisStorage)(nil).CreateSession), ctx, session, expire)
}

// DeleteSession mocks base method.
func (m *MockRedisStorage) DeleteSession(ctx context.Context, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", ctx, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockRedisStorageMockRecorder) DeleteSession(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockRedisStorage)(nil).DeleteSession), ctx, refreshToken)
}

// GetUserID mocks base method.
func (m *MockRedisStorage) GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserID", ctx, refreshToken)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserID indicates an expected call of GetUserID.
func (mr *MockRedisStorageMockRecorder) GetUserID(ctx, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserID", reflect.TypeOf((*MockRedisStorage)(nil).GetUserID), ctx, refreshToken)
}
