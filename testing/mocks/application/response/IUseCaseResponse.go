// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
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

// GetAllCandidateResponses provides a mock function with given fields: candId, respIds
func (_m *IUseCaseResponse) GetAllCandidateResponses(candId uuid.UUID, respIds []uuid.UUID) ([]models.ResponseWithTitle, error) {
	ret := _m.Called(candId, respIds)

	var r0 []models.ResponseWithTitle
	if rf, ok := ret.Get(0).(func(uuid.UUID, []uuid.UUID) []models.ResponseWithTitle); ok {
		r0 = rf(candId, respIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResponseWithTitle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, []uuid.UUID) error); ok {
		r1 = rf(candId, respIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllEmployerResponses provides a mock function with given fields: emplId, respIds
func (_m *IUseCaseResponse) GetAllEmployerResponses(emplId uuid.UUID, respIds []uuid.UUID) ([]models.ResponseWithTitle, error) {
	ret := _m.Called(emplId, respIds)

	var r0 []models.ResponseWithTitle
	if rf, ok := ret.Get(0).(func(uuid.UUID, []uuid.UUID) []models.ResponseWithTitle); ok {
		r0 = rf(emplId, respIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ResponseWithTitle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, []uuid.UUID) error); ok {
		r1 = rf(emplId, respIds)
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

// GetRecommendedVacCnt provides a mock function with given fields: candId, daysFromNow
func (_m *IUseCaseResponse) GetRecommendedVacCnt(candId uuid.UUID, daysFromNow int) (uint, error) {
	ret := _m.Called(candId, daysFromNow)

	var r0 uint
	if rf, ok := ret.Get(0).(func(uuid.UUID, int) uint); ok {
		r0 = rf(candId, daysFromNow)
	} else {
		r0 = ret.Get(0).(uint)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, int) error); ok {
		r1 = rf(candId, daysFromNow)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecommendedVacancies provides a mock function with given fields: candId, start, limit, daysFromNow
func (_m *IUseCaseResponse) GetRecommendedVacancies(candId uuid.UUID, start uint, limit uint, daysFromNow int) ([]models.Vacancy, error) {
	ret := _m.Called(candId, start, limit, daysFromNow)

	var r0 []models.Vacancy
	if rf, ok := ret.Get(0).(func(uuid.UUID, uint, uint, int) []models.Vacancy); ok {
		r0 = rf(candId, start, limit, daysFromNow)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Vacancy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uint, uint, int) error); ok {
		r1 = rf(candId, start, limit, daysFromNow)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResponsesCnt provides a mock function with given fields: userId, userType
func (_m *IUseCaseResponse) GetResponsesCnt(userId uuid.UUID, userType string) (uint, error) {
	ret := _m.Called(userId, userType)

	var r0 uint
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) uint); ok {
		r0 = rf(userId, userType)
	} else {
		r0 = ret.Get(0).(uint)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(userId, userType)
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
