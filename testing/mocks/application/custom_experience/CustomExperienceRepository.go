// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// CustomExperienceRepository is an autogenerated mock type for the CustomExperienceRepository type
type CustomExperienceRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: experience
func (_m *CustomExperienceRepository) Create(experience models.ExperienceCustomComp) (*models.ExperienceCustomComp, error) {
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
func (_m *CustomExperienceRepository) DropAllFromResume(resumeID uuid.UUID) error {
	ret := _m.Called(resumeID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(resumeID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
