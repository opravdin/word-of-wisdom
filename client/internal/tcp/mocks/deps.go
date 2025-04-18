// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go
//
// Generated by this command:
//
//	mockgen -source=deps.go -package=mocks -destination=./mocks/deps.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/opravdin/word-of-wisdom/client/internal/domain"
	tcp "github.com/opravdin/word-of-wisdom/client/internal/tcp"
	gomock "go.uber.org/mock/gomock"
)

// MockTCPClient is a mock of TCPClient interface.
type MockTCPClient struct {
	ctrl     *gomock.Controller
	recorder *MockTCPClientMockRecorder
	isgomock struct{}
}

// MockTCPClientMockRecorder is the mock recorder for MockTCPClient.
type MockTCPClientMockRecorder struct {
	mock *MockTCPClient
}

// NewMockTCPClient creates a new mock instance.
func NewMockTCPClient(ctrl *gomock.Controller) *MockTCPClient {
	mock := &MockTCPClient{ctrl: ctrl}
	mock.recorder = &MockTCPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTCPClient) EXPECT() *MockTCPClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockTCPClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTCPClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTCPClient)(nil).Close))
}

// GetQuote mocks base method.
func (m *MockTCPClient) GetQuote(ctx context.Context) (*domain.Quote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuote", ctx)
	ret0, _ := ret[0].(*domain.Quote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuote indicates an expected call of GetQuote.
func (mr *MockTCPClientMockRecorder) GetQuote(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuote", reflect.TypeOf((*MockTCPClient)(nil).GetQuote), ctx)
}

// ProcessChallenge mocks base method.
func (m *MockTCPClient) ProcessChallenge(ctx context.Context, msg domain.Message, solveFn func(string, string, int, int, int, int, int) (string, error)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessChallenge", ctx, msg, solveFn)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessChallenge indicates an expected call of ProcessChallenge.
func (mr *MockTCPClientMockRecorder) ProcessChallenge(ctx, msg, solveFn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessChallenge", reflect.TypeOf((*MockTCPClient)(nil).ProcessChallenge), ctx, msg, solveFn)
}

// ProcessChallengeWithDefaults mocks base method.
func (m *MockTCPClient) ProcessChallengeWithDefaults(ctx context.Context, msg domain.Message, solveFn func(string, string) (string, error)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessChallengeWithDefaults", ctx, msg, solveFn)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessChallengeWithDefaults indicates an expected call of ProcessChallengeWithDefaults.
func (mr *MockTCPClientMockRecorder) ProcessChallengeWithDefaults(ctx, msg, solveFn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessChallengeWithDefaults", reflect.TypeOf((*MockTCPClient)(nil).ProcessChallengeWithDefaults), ctx, msg, solveFn)
}

// ReadMessage mocks base method.
func (m *MockTCPClient) ReadMessage(ctx context.Context) (domain.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage", ctx)
	ret0, _ := ret[0].(domain.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMessage indicates an expected call of ReadMessage.
func (mr *MockTCPClientMockRecorder) ReadMessage(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockTCPClient)(nil).ReadMessage), ctx)
}

// SendMessage mocks base method.
func (m *MockTCPClient) SendMessage(ctx context.Context, msg domain.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockTCPClientMockRecorder) SendMessage(ctx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockTCPClient)(nil).SendMessage), ctx, msg)
}

// MockTCPClientFactory is a mock of TCPClientFactory interface.
type MockTCPClientFactory struct {
	ctrl     *gomock.Controller
	recorder *MockTCPClientFactoryMockRecorder
	isgomock struct{}
}

// MockTCPClientFactoryMockRecorder is the mock recorder for MockTCPClientFactory.
type MockTCPClientFactoryMockRecorder struct {
	mock *MockTCPClientFactory
}

// NewMockTCPClientFactory creates a new mock instance.
func NewMockTCPClientFactory(ctrl *gomock.Controller) *MockTCPClientFactory {
	mock := &MockTCPClientFactory{ctrl: ctrl}
	mock.recorder = &MockTCPClientFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTCPClientFactory) EXPECT() *MockTCPClientFactoryMockRecorder {
	return m.recorder
}

// NewClient mocks base method.
func (m *MockTCPClientFactory) NewClient(ctx context.Context, serverAddr string) (tcp.TCPClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClient", ctx, serverAddr)
	ret0, _ := ret[0].(tcp.TCPClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewClient indicates an expected call of NewClient.
func (mr *MockTCPClientFactoryMockRecorder) NewClient(ctx, serverAddr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClient", reflect.TypeOf((*MockTCPClientFactory)(nil).NewClient), ctx, serverAddr)
}
