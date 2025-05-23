// Code generated by MockGen. DO NOT EDIT.
// Source: internal/requests/graphql.go

// Package mock_requests is a generated GoMock package.
package mock_requests

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGraphQLClient is a mock of GraphQLClient interface.
type MockGraphQLClient struct {
	ctrl     *gomock.Controller
	recorder *MockGraphQLClientMockRecorder
}

// MockGraphQLClientMockRecorder is the mock recorder for MockGraphQLClient.
type MockGraphQLClientMockRecorder struct {
	mock *MockGraphQLClient
}

// NewMockGraphQLClient creates a new mock instance.
func NewMockGraphQLClient(ctrl *gomock.Controller) *MockGraphQLClient {
	mock := &MockGraphQLClient{ctrl: ctrl}
	mock.recorder = &MockGraphQLClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGraphQLClient) EXPECT() *MockGraphQLClientMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockGraphQLClient) Do(query string, variables map[string]interface{}, response interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", query, variables, response)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do.
func (mr *MockGraphQLClientMockRecorder) Do(query, variables, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockGraphQLClient)(nil).Do), query, variables, response)
}

// Query mocks base method.
func (m *MockGraphQLClient) Query(name string, q interface{}, variables map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", name, q, variables)
	ret0, _ := ret[0].(error)
	return ret0
}

// Query indicates an expected call of Query.
func (mr *MockGraphQLClientMockRecorder) Query(name, q, variables interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockGraphQLClient)(nil).Query), name, q, variables)
}
