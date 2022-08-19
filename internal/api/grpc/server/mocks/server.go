// Code generated by MockGen. DO NOT EDIT.
// Source: ./../../../../pkg/api/server/server_grpc.pb.go

// Package mock_server is a generated GoMock package.
package mock_server

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	server "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockCinemaBackendClient is a mock of CinemaBackendClient interface.
type MockCinemaBackendClient struct {
	ctrl     *gomock.Controller
	recorder *MockCinemaBackendClientMockRecorder
}

// MockCinemaBackendClientMockRecorder is the mock recorder for MockCinemaBackendClient.
type MockCinemaBackendClientMockRecorder struct {
	mock *MockCinemaBackendClient
}

// NewMockCinemaBackendClient creates a new mock instance.
func NewMockCinemaBackendClient(ctrl *gomock.Controller) *MockCinemaBackendClient {
	mock := &MockCinemaBackendClient{ctrl: ctrl}
	mock.recorder = &MockCinemaBackendClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCinemaBackendClient) EXPECT() *MockCinemaBackendClientMockRecorder {
	return m.recorder
}

// FilmRoom mocks base method.
func (m *MockCinemaBackendClient) FilmRoom(ctx context.Context, in *server.FilmRoomRequest, opts ...grpc.CallOption) (*server.FilmRoomResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FilmRoom", varargs...)
	ret0, _ := ret[0].(*server.FilmRoomResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilmRoom indicates an expected call of FilmRoom.
func (mr *MockCinemaBackendClientMockRecorder) FilmRoom(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilmRoom", reflect.TypeOf((*MockCinemaBackendClient)(nil).FilmRoom), varargs...)
}

// Films mocks base method.
func (m *MockCinemaBackendClient) Films(ctx context.Context, in *server.FilmsRequest, opts ...grpc.CallOption) (server.CinemaBackend_FilmsClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Films", varargs...)
	ret0, _ := ret[0].(server.CinemaBackend_FilmsClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Films indicates an expected call of Films.
func (mr *MockCinemaBackendClientMockRecorder) Films(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Films", reflect.TypeOf((*MockCinemaBackendClient)(nil).Films), varargs...)
}

// MyTickets mocks base method.
func (m *MockCinemaBackendClient) MyTickets(ctx context.Context, in *server.MyTicketsRequest, opts ...grpc.CallOption) (*server.MyTicketsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MyTickets", varargs...)
	ret0, _ := ret[0].(*server.MyTicketsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MyTickets indicates an expected call of MyTickets.
func (mr *MockCinemaBackendClientMockRecorder) MyTickets(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MyTickets", reflect.TypeOf((*MockCinemaBackendClient)(nil).MyTickets), varargs...)
}

// TicketCreate mocks base method.
func (m *MockCinemaBackendClient) TicketCreate(ctx context.Context, in *server.TicketCreateRequest, opts ...grpc.CallOption) (*server.TicketCreateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TicketCreate", varargs...)
	ret0, _ := ret[0].(*server.TicketCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TicketCreate indicates an expected call of TicketCreate.
func (mr *MockCinemaBackendClientMockRecorder) TicketCreate(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TicketCreate", reflect.TypeOf((*MockCinemaBackendClient)(nil).TicketCreate), varargs...)
}

// TicketDelete mocks base method.
func (m *MockCinemaBackendClient) TicketDelete(ctx context.Context, in *server.TicketDeleteRequest, opts ...grpc.CallOption) (*server.TicketDeleteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TicketDelete", varargs...)
	ret0, _ := ret[0].(*server.TicketDeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TicketDelete indicates an expected call of TicketDelete.
func (mr *MockCinemaBackendClientMockRecorder) TicketDelete(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TicketDelete", reflect.TypeOf((*MockCinemaBackendClient)(nil).TicketDelete), varargs...)
}

// UserAuth mocks base method.
func (m *MockCinemaBackendClient) UserAuth(ctx context.Context, in *server.UserAuthRequest, opts ...grpc.CallOption) (*server.UserAuthResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UserAuth", varargs...)
	ret0, _ := ret[0].(*server.UserAuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserAuth indicates an expected call of UserAuth.
func (mr *MockCinemaBackendClientMockRecorder) UserAuth(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAuth", reflect.TypeOf((*MockCinemaBackendClient)(nil).UserAuth), varargs...)
}

// MockCinemaBackend_FilmsClient is a mock of CinemaBackend_FilmsClient interface.
type MockCinemaBackend_FilmsClient struct {
	ctrl     *gomock.Controller
	recorder *MockCinemaBackend_FilmsClientMockRecorder
}

// MockCinemaBackend_FilmsClientMockRecorder is the mock recorder for MockCinemaBackend_FilmsClient.
type MockCinemaBackend_FilmsClientMockRecorder struct {
	mock *MockCinemaBackend_FilmsClient
}

// NewMockCinemaBackend_FilmsClient creates a new mock instance.
func NewMockCinemaBackend_FilmsClient(ctrl *gomock.Controller) *MockCinemaBackend_FilmsClient {
	mock := &MockCinemaBackend_FilmsClient{ctrl: ctrl}
	mock.recorder = &MockCinemaBackend_FilmsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCinemaBackend_FilmsClient) EXPECT() *MockCinemaBackend_FilmsClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockCinemaBackend_FilmsClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockCinemaBackend_FilmsClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockCinemaBackend_FilmsClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockCinemaBackend_FilmsClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockCinemaBackend_FilmsClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockCinemaBackend_FilmsClient)(nil).Context))
}

// Header mocks base method.
func (m *MockCinemaBackend_FilmsClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockCinemaBackend_FilmsClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockCinemaBackend_FilmsClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockCinemaBackend_FilmsClient) Recv() (*server.FilmsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*server.FilmsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockCinemaBackend_FilmsClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockCinemaBackend_FilmsClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockCinemaBackend_FilmsClient) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockCinemaBackend_FilmsClientMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockCinemaBackend_FilmsClient)(nil).RecvMsg), m)
}

// SendMsg mocks base method.
func (m_2 *MockCinemaBackend_FilmsClient) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockCinemaBackend_FilmsClientMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockCinemaBackend_FilmsClient)(nil).SendMsg), m)
}

// Trailer mocks base method.
func (m *MockCinemaBackend_FilmsClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockCinemaBackend_FilmsClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockCinemaBackend_FilmsClient)(nil).Trailer))
}

// MockCinemaBackendServer is a mock of CinemaBackendServer interface.
type MockCinemaBackendServer struct {
	ctrl     *gomock.Controller
	recorder *MockCinemaBackendServerMockRecorder
}

// MockCinemaBackendServerMockRecorder is the mock recorder for MockCinemaBackendServer.
type MockCinemaBackendServerMockRecorder struct {
	mock *MockCinemaBackendServer
}

// NewMockCinemaBackendServer creates a new mock instance.
func NewMockCinemaBackendServer(ctrl *gomock.Controller) *MockCinemaBackendServer {
	mock := &MockCinemaBackendServer{ctrl: ctrl}
	mock.recorder = &MockCinemaBackendServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCinemaBackendServer) EXPECT() *MockCinemaBackendServerMockRecorder {
	return m.recorder
}

// FilmRoom mocks base method.
func (m *MockCinemaBackendServer) FilmRoom(arg0 context.Context, arg1 *server.FilmRoomRequest) (*server.FilmRoomResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilmRoom", arg0, arg1)
	ret0, _ := ret[0].(*server.FilmRoomResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilmRoom indicates an expected call of FilmRoom.
func (mr *MockCinemaBackendServerMockRecorder) FilmRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilmRoom", reflect.TypeOf((*MockCinemaBackendServer)(nil).FilmRoom), arg0, arg1)
}

// Films mocks base method.
func (m *MockCinemaBackendServer) Films(arg0 *server.FilmsRequest, arg1 server.CinemaBackend_FilmsServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Films", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Films indicates an expected call of Films.
func (mr *MockCinemaBackendServerMockRecorder) Films(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Films", reflect.TypeOf((*MockCinemaBackendServer)(nil).Films), arg0, arg1)
}

// MyTickets mocks base method.
func (m *MockCinemaBackendServer) MyTickets(arg0 context.Context, arg1 *server.MyTicketsRequest) (*server.MyTicketsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MyTickets", arg0, arg1)
	ret0, _ := ret[0].(*server.MyTicketsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MyTickets indicates an expected call of MyTickets.
func (mr *MockCinemaBackendServerMockRecorder) MyTickets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MyTickets", reflect.TypeOf((*MockCinemaBackendServer)(nil).MyTickets), arg0, arg1)
}

// TicketCreate mocks base method.
func (m *MockCinemaBackendServer) TicketCreate(arg0 context.Context, arg1 *server.TicketCreateRequest) (*server.TicketCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TicketCreate", arg0, arg1)
	ret0, _ := ret[0].(*server.TicketCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TicketCreate indicates an expected call of TicketCreate.
func (mr *MockCinemaBackendServerMockRecorder) TicketCreate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TicketCreate", reflect.TypeOf((*MockCinemaBackendServer)(nil).TicketCreate), arg0, arg1)
}

// TicketDelete mocks base method.
func (m *MockCinemaBackendServer) TicketDelete(arg0 context.Context, arg1 *server.TicketDeleteRequest) (*server.TicketDeleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TicketDelete", arg0, arg1)
	ret0, _ := ret[0].(*server.TicketDeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TicketDelete indicates an expected call of TicketDelete.
func (mr *MockCinemaBackendServerMockRecorder) TicketDelete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TicketDelete", reflect.TypeOf((*MockCinemaBackendServer)(nil).TicketDelete), arg0, arg1)
}

// UserAuth mocks base method.
func (m *MockCinemaBackendServer) UserAuth(arg0 context.Context, arg1 *server.UserAuthRequest) (*server.UserAuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserAuth", arg0, arg1)
	ret0, _ := ret[0].(*server.UserAuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserAuth indicates an expected call of UserAuth.
func (mr *MockCinemaBackendServerMockRecorder) UserAuth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserAuth", reflect.TypeOf((*MockCinemaBackendServer)(nil).UserAuth), arg0, arg1)
}

// mustEmbedUnimplementedCinemaBackendServer mocks base method.
func (m *MockCinemaBackendServer) mustEmbedUnimplementedCinemaBackendServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCinemaBackendServer")
}

// mustEmbedUnimplementedCinemaBackendServer indicates an expected call of mustEmbedUnimplementedCinemaBackendServer.
func (mr *MockCinemaBackendServerMockRecorder) mustEmbedUnimplementedCinemaBackendServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCinemaBackendServer", reflect.TypeOf((*MockCinemaBackendServer)(nil).mustEmbedUnimplementedCinemaBackendServer))
}

// MockUnsafeCinemaBackendServer is a mock of UnsafeCinemaBackendServer interface.
type MockUnsafeCinemaBackendServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeCinemaBackendServerMockRecorder
}

// MockUnsafeCinemaBackendServerMockRecorder is the mock recorder for MockUnsafeCinemaBackendServer.
type MockUnsafeCinemaBackendServerMockRecorder struct {
	mock *MockUnsafeCinemaBackendServer
}

// NewMockUnsafeCinemaBackendServer creates a new mock instance.
func NewMockUnsafeCinemaBackendServer(ctrl *gomock.Controller) *MockUnsafeCinemaBackendServer {
	mock := &MockUnsafeCinemaBackendServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeCinemaBackendServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeCinemaBackendServer) EXPECT() *MockUnsafeCinemaBackendServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedCinemaBackendServer mocks base method.
func (m *MockUnsafeCinemaBackendServer) mustEmbedUnimplementedCinemaBackendServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCinemaBackendServer")
}

// mustEmbedUnimplementedCinemaBackendServer indicates an expected call of mustEmbedUnimplementedCinemaBackendServer.
func (mr *MockUnsafeCinemaBackendServerMockRecorder) mustEmbedUnimplementedCinemaBackendServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCinemaBackendServer", reflect.TypeOf((*MockUnsafeCinemaBackendServer)(nil).mustEmbedUnimplementedCinemaBackendServer))
}

// MockCinemaBackend_FilmsServer is a mock of CinemaBackend_FilmsServer interface.
type MockCinemaBackend_FilmsServer struct {
	ctrl     *gomock.Controller
	recorder *MockCinemaBackend_FilmsServerMockRecorder
}

// MockCinemaBackend_FilmsServerMockRecorder is the mock recorder for MockCinemaBackend_FilmsServer.
type MockCinemaBackend_FilmsServerMockRecorder struct {
	mock *MockCinemaBackend_FilmsServer
}

// NewMockCinemaBackend_FilmsServer creates a new mock instance.
func NewMockCinemaBackend_FilmsServer(ctrl *gomock.Controller) *MockCinemaBackend_FilmsServer {
	mock := &MockCinemaBackend_FilmsServer{ctrl: ctrl}
	mock.recorder = &MockCinemaBackend_FilmsServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCinemaBackend_FilmsServer) EXPECT() *MockCinemaBackend_FilmsServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockCinemaBackend_FilmsServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockCinemaBackend_FilmsServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockCinemaBackend_FilmsServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m_2 *MockCinemaBackend_FilmsServer) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockCinemaBackend_FilmsServerMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockCinemaBackend_FilmsServer)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockCinemaBackend_FilmsServer) Send(arg0 *server.FilmsResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockCinemaBackend_FilmsServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockCinemaBackend_FilmsServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockCinemaBackend_FilmsServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockCinemaBackend_FilmsServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockCinemaBackend_FilmsServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockCinemaBackend_FilmsServer) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockCinemaBackend_FilmsServerMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockCinemaBackend_FilmsServer)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockCinemaBackend_FilmsServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockCinemaBackend_FilmsServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockCinemaBackend_FilmsServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockCinemaBackend_FilmsServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockCinemaBackend_FilmsServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockCinemaBackend_FilmsServer)(nil).SetTrailer), arg0)
}
