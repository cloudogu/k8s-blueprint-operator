// Code generated by mockery v2.42.1. DO NOT EDIT.

package ecosystem

import mock "github.com/stretchr/testify/mock"

// MockEcosystemConfigEntry is an autogenerated mock type for the EcosystemConfigEntry type
type MockEcosystemConfigEntry struct {
	mock.Mock
}

type MockEcosystemConfigEntry_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEcosystemConfigEntry) EXPECT() *MockEcosystemConfigEntry_Expecter {
	return &MockEcosystemConfigEntry_Expecter{mock: &_m.Mock}
}

// NewMockEcosystemConfigEntry creates a new instance of MockEcosystemConfigEntry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEcosystemConfigEntry(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEcosystemConfigEntry {
	mock := &MockEcosystemConfigEntry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}