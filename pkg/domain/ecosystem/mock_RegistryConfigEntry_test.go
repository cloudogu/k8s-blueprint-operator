// Code generated by mockery v2.42.1. DO NOT EDIT.

package ecosystem

import mock "github.com/stretchr/testify/mock"

// MockRegistryConfigEntry is an autogenerated mock type for the RegistryConfigEntry type
type MockRegistryConfigEntry struct {
	mock.Mock
}

type MockRegistryConfigEntry_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRegistryConfigEntry) EXPECT() *MockRegistryConfigEntry_Expecter {
	return &MockRegistryConfigEntry_Expecter{mock: &_m.Mock}
}

// NewMockRegistryConfigEntry creates a new instance of MockRegistryConfigEntry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRegistryConfigEntry(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRegistryConfigEntry {
	mock := &MockRegistryConfigEntry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
