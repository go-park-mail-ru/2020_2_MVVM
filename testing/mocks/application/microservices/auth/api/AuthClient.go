// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/microservises/auth"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// AuthClient is an autogenerated mock type for the AuthClient type
type AuthClient struct {
	mock.Mock
}

// Check provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Check(ctx context.Context, in *auth.SessionID, opts ...grpc.CallOption) (*auth.SessionInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.SessionInfo
	if rf, ok := ret.Get(0).(func(context.Context, *auth.SessionID, ...grpc.CallOption) *auth.SessionInfo); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.SessionInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.SessionID, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Login(ctx context.Context, in *auth.Credentials, opts ...grpc.CallOption) (*auth.SessionInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.SessionInfo
	if rf, ok := ret.Get(0).(func(context.Context, *auth.Credentials, ...grpc.CallOption) *auth.SessionInfo); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.SessionInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.Credentials, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: ctx, in, opts
func (_m *AuthClient) Logout(ctx context.Context, in *auth.SessionID, opts ...grpc.CallOption) (*auth.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *auth.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *auth.SessionID, ...grpc.CallOption) *auth.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *auth.SessionID, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
