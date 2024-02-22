// Code generated by mockery v2.20.0. DO NOT EDIT.

package etcd

import mock "github.com/stretchr/testify/mock"

// mockRegistryEntry is an autogenerated mock type for the registryEntry type
type mockRegistryEntry struct {
	mock.Mock
}

type mockRegistryEntry_Expecter struct {
	mock *mock.Mock
}

func (_m *mockRegistryEntry) EXPECT() *mockRegistryEntry_Expecter {
	return &mockRegistryEntry_Expecter{mock: &_m.Mock}
}

type mockConstructorTestingTnewMockRegistryEntry interface {
	mock.TestingT
	Cleanup(func())
}

// newMockRegistryEntry creates a new instance of mockRegistryEntry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockRegistryEntry(t mockConstructorTestingTnewMockRegistryEntry) *mockRegistryEntry {
	mock := &mockRegistryEntry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}