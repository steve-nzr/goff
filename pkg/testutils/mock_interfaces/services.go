// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/steve-nzr/goff-server/internal/domain/interfaces (interfaces: IdentifierGenerator)

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	gomock "github.com/golang/mock/gomock"
	customtypes "github.com/steve-nzr/goff-server/internal/domain/customtypes"
	reflect "reflect"
)

// MockIdentifierGenerator is a mock of IdentifierGenerator interface
type MockIdentifierGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockIdentifierGeneratorMockRecorder
}

// MockIdentifierGeneratorMockRecorder is the mock recorder for MockIdentifierGenerator
type MockIdentifierGeneratorMockRecorder struct {
	mock *MockIdentifierGenerator
}

// NewMockIdentifierGenerator creates a new mock instance
func NewMockIdentifierGenerator(ctrl *gomock.Controller) *MockIdentifierGenerator {
	mock := &MockIdentifierGenerator{ctrl: ctrl}
	mock.recorder = &MockIdentifierGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIdentifierGenerator) EXPECT() *MockIdentifierGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockIdentifierGenerator) Generate() customtypes.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(customtypes.ID)
	return ret0
}

// Generate indicates an expected call of Generate
func (mr *MockIdentifierGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockIdentifierGenerator)(nil).Generate))
}