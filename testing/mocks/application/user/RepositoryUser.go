// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	mock "github.com/stretchr/testify/mock"
)

// RepositoryUser is an autogenerated mock type for the RepositoryUser type
type RepositoryUser struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: _a0
func (_m *RepositoryUser) CreateUser(_a0 models.User) (*models.User, error) {
	ret := _m.Called(_a0)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(models.User) *models.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCandByID provides a mock function with given fields: id
func (_m *RepositoryUser) GetCandByID(id string) (*models.User, error) {
	ret := _m.Called(id)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCandidateByID provides a mock function with given fields: id
func (_m *RepositoryUser) GetCandidateByID(id string) (*models.Candidate, error) {
	ret := _m.Called(id)

	var r0 *models.Candidate
	if rf, ok := ret.Get(0).(func(string) *models.Candidate); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Candidate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEmplByID provides a mock function with given fields: id
func (_m *RepositoryUser) GetEmplByID(id string) (*models.User, error) {
	ret := _m.Called(id)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEmployerByID provides a mock function with given fields: id
func (_m *RepositoryUser) GetEmployerByID(id string) (*models.Employer, error) {
	ret := _m.Called(id)

	var r0 *models.Employer
	if rf, ok := ret.Get(0).(func(string) *models.Employer); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Employer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: id
func (_m *RepositoryUser) GetUserByID(id string) (*models.User, error) {
	ret := _m.Called(id)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: _a0
func (_m *RepositoryUser) Login(_a0 models.UserLogin) (*models.User, error) {
	ret := _m.Called(_a0)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(models.UserLogin) *models.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.UserLogin) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: _a0
func (_m *RepositoryUser) UpdateUser(_a0 models.User) (*models.User, error) {
	ret := _m.Called(_a0)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(models.User) *models.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}