// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// IUseCaseResponse is an autogenerated mock type for the IUseCaseResponse type
type IUseCaseResponse struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *IUseCaseResponse) Create(_a0 models.Response) (*models.Response, error) {
	ret := _m.Called(_a0)

	var r0 *models.Response
	if rf, ok := ret.Get(0).(func(models.Response) *models.Response); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Response) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllCandidateResponses provides a mock function with given fields: _a0
func (_m *IUseCaseResponse) GetAllCandidateResponses(_a0 uuid.UUID) ([]models.ResponseWithTitle, error) {
	ret := _m.Called(_a0)

	var r0 []models.ResponseWithTitle
	if rf, ok := ret.Get(0).(func(uuid.UUID) []models.ResponseWithTitle); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResponseWithTitle)
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

// GetAllEmployerResponses provides a mock function with given fields: _a0
func (_m *IUseCaseResponse) GetAllEmployerResponses(_a0 uuid.UUID) ([]models.ResponseWithTitle, error) {
	ret := _m.Called(_a0)

	var r0 []models.ResponseWithTitle
	if rf, ok := ret.Get(0).(func(uuid.UUID) []models.ResponseWithTitle); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResponseWithTitle)
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

// GetAllResumeWithoutResponse provides a mock function with given fields: candID, vacancyID
func (_m *IUseCaseResponse) GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.BriefResumeInfo, error) {
	ret := _m.Called(candID, vacancyID)

	var r0 []models.BriefResumeInfo
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) []models.BriefResumeInfo); ok {
		r0 = rf(candID, vacancyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.BriefResumeInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(candID, vacancyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllVacancyWithoutResponse provides a mock function with given fields: emplID, resumeID
func (_m *IUseCaseResponse) GetAllVacancyWithoutResponse(emplID uuid.UUID, resumeID uuid.UUID) ([]models.Vacancy, error) {
	ret := _m.Called(emplID, resumeID)

	var r0 []models.Vacancy
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) []models.Vacancy); ok {
		r0 = rf(emplID, resumeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Vacancy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(emplID, resumeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: _a0, userType
func (_m *IUseCaseResponse) UpdateStatus(_a0 models.Response, userType string) (*models.Response, error) {
	ret := _m.Called(_a0, userType)

	var r0 *models.Response
	if rf, ok := ret.Get(0).(func(models.Response, string) *models.Response); ok {
		r0 = rf(_a0, userType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Response, string) error); ok {
		r1 = rf(_a0, userType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}