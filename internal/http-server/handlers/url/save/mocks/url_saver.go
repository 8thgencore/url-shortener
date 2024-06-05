// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// URLSaver is an autogenerated mock type for the URLSaver type
type URLSaver struct {
	mock.Mock
}

// AliasExists provides a mock function with given fields: alias
func (_m *URLSaver) AliasExists(alias string) (bool, error) {
	ret := _m.Called(alias)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAliasByURL provides a mock function with given fields: urlToFind
func (_m *URLSaver) GetAliasByURL(urlToFind string) (string, error) {
	ret := _m.Called(urlToFind)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(urlToFind)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(urlToFind)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(urlToFind)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveURL provides a mock function with given fields: urlToSave, alias
func (_m *URLSaver) SaveURL(urlToSave string, alias string) (int64, error) {
	ret := _m.Called(urlToSave, alias)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (int64, error)); ok {
		return rf(urlToSave, alias)
	}
	if rf, ok := ret.Get(0).(func(string, string) int64); ok {
		r0 = rf(urlToSave, alias)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(urlToSave, alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// URLExists provides a mock function with given fields: urlToCheck
func (_m *URLSaver) URLExists(urlToCheck string) (bool, error) {
	ret := _m.Called(urlToCheck)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(urlToCheck)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(urlToCheck)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(urlToCheck)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewURLSaver creates a new instance of URLSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLSaver {
	mock := &URLSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}