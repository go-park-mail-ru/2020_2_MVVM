// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// OfficialCompanyRepository is an autogenerated mock type for the OfficialCompanyRepository type
type OfficialCompanyRepository struct {
	mock.Mock
}

// CreateOfficialCompany provides a mock function with given fields: _a0, _a1
func (_m *OfficialCompanyRepository) CreateOfficialCompany(_a0 models.OfficialCompany, _a1 uuid.UUID) (*models.OfficialCompany, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *models.OfficialCompany
	if rf, ok := ret.Get(0).(func(models.OfficialCompany, uuid.UUID) *models.OfficialCompany); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OfficialCompany)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.OfficialCompany, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOfficialCompany provides a mock function with given fields: _a0, _a1
func (_m *OfficialCompanyRepository) DeleteOfficialCompany(_a0 uuid.UUID, _a1 uuid.UUID) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllCompaniesNames provides a mock function with given fields:
func (_m *OfficialCompanyRepository) GetAllCompaniesNames() ([]models.BriefCompany, error) {
	ret := _m.Called()

	var r0 []models.BriefCompany
	if rf, ok := ret.Get(0).(func() []models.BriefCompany); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.BriefCompany)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCompaniesList provides a mock function with given fields: _a0, _a1
func (_m *OfficialCompanyRepository) GetCompaniesList(_a0 uint, _a1 uint) ([]models.OfficialCompany, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []models.OfficialCompany
	if rf, ok := ret.Get(0).(func(uint, uint) []models.OfficialCompany); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OfficialCompany)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, uint) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMineCompany provides a mock function with given fields: _a0
func (_m *OfficialCompanyRepository) GetMineCompany(_a0 uuid.UUID) (*models.OfficialCompany, error) {
	ret := _m.Called(_a0)

	var r0 *models.OfficialCompany
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.OfficialCompany); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OfficialCompany)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOfficialCompany provides a mock function with given fields: _a0
func (_m *OfficialCompanyRepository) GetOfficialCompany(_a0 uuid.UUID) (*models.OfficialCompany, error) {
	ret := _m.Called(_a0)

	var r0 *models.OfficialCompany
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.OfficialCompany); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OfficialCompany)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchCompanies provides a mock function with given fields: _a0
func (_m *OfficialCompanyRepository) SearchCompanies(_a0 models.CompanySearchParams) ([]models.OfficialCompany, error) {
	ret := _m.Called(_a0)

	var r0 []models.OfficialCompany
	if rf, ok := ret.Get(0).(func(models.CompanySearchParams) []models.OfficialCompany); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.OfficialCompany)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.CompanySearchParams) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOfficialCompany provides a mock function with given fields: _a0, _a1
func (_m *OfficialCompanyRepository) UpdateOfficialCompany(_a0 models.OfficialCompany, _a1 uuid.UUID) (*models.OfficialCompany, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *models.OfficialCompany
	if rf, ok := ret.Get(0).(func(models.OfficialCompany, uuid.UUID) *models.OfficialCompany); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.OfficialCompany)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.OfficialCompany, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
