// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package custom_experience

import (
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// MockUseCase is an autogenerated mock type for the UseCase type
type MockUseCase struct {
	mock.Mock
}

// Create provides a mock function with given fields: experience
func (_m *MockUseCase) Create(experience models.ExperienceCustomComp) (*models.ExperienceCustomComp, error) {
	ret := _m.Called(experience)

	var r0 *models.ExperienceCustomComp
	if rf, ok := ret.Get(0).(func(models.ExperienceCustomComp) *models.ExperienceCustomComp); ok {
		r0 = rf(experience)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ExperienceCustomComp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.ExperienceCustomComp) error); ok {
		r1 = rf(experience)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DropAllFromResume provides a mock function with given fields: resumeID
func (_m *MockUseCase) DropAllFromResume(resumeID uuid.UUID) error {
	ret := _m.Called(resumeID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(resumeID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllFromResume provides a mock function with given fields: ResumeID
func (_m *MockUseCase) GetAllFromResume(ResumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	ret := _m.Called(ResumeID)

	var r0 []models.ExperienceCustomComp
	if rf, ok := ret.Get(0).(func(uuid.UUID) []models.ExperienceCustomComp); ok {
		r0 = rf(ResumeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ExperienceCustomComp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(ResumeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *MockUseCase) GetById(id string) (*models.ExperienceCustomComp, error) {
	ret := _m.Called(id)

	var r0 *models.ExperienceCustomComp
	if rf, ok := ret.Get(0).(func(string) *models.ExperienceCustomComp); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ExperienceCustomComp)
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
