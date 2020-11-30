// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	common "github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	mock "github.com/stretchr/testify/mock"
)

// AuthClient is an autogenerated mock type for the AuthClient type
type AuthClient struct {
	mock.Mock
}

// Check provides a mock function with given fields: sessionID
func (_m *AuthClient) Check(sessionID string) (common.Session, error) {
	ret := _m.Called(sessionID)

	var r0 common.Session
	if rf, ok := ret.Get(0).(func(string) common.Session); ok {
		r0 = rf(sessionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sessionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: login, password
func (_m *AuthClient) Login(login string, password string) (common.Session, error) {
	ret := _m.Called(login, password)

	var r0 common.Session
	if rf, ok := ret.Get(0).(func(string, string) common.Session); ok {
		r0 = rf(login, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(login, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: sessionID
func (_m *AuthClient) Logout(sessionID string) error {
	ret := _m.Called(sessionID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(sessionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}