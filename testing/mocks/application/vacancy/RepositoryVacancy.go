// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	dtovacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/models/vacancy"
	mock "github.com/stretchr/testify/mock"

	models "github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"

	uuid "github.com/google/uuid"
)

// RepositoryVacancy is an autogenerated mock type for the RepositoryVacancy type
type RepositoryVacancy struct {
	mock.Mock
}

// AddRecommendation provides a mock function with given fields: _a0, _a1
func (_m *RepositoryVacancy) AddRecommendation(_a0 uuid.UUID, _a1 int) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, int) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateVacancy provides a mock function with given fields: _a0
func (_m *RepositoryVacancy) CreateVacancy(_a0 models.Vacancy) (*models.Vacancy, error) {
	ret := _m.Called(_a0)

	var r0 *models.Vacancy
	if rf, ok := ret.Get(0).(func(models.Vacancy) *models.Vacancy); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Vacancy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Vacancy) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPreferredSalary provides a mock function with given fields: _a0
func (_m *RepositoryVacancy) GetPreferredSalary(_a0 uuid.UUID) (*float64, error) {
	ret := _m.Called(_a0)

	var r0 *float64
	if rf, ok := ret.Get(0).(func(uuid.UUID) *float64); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*float64)
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

// GetPreferredSpheres provides a mock function with given fields: userID
func (_m *RepositoryVacancy) GetPreferredSpheres(userID uuid.UUID) ([]dtovacancy.Pair, error) {
	ret := _m.Called(userID)

	var r0 []dtovacancy.Pair
	if rf, ok := ret.Get(0).(func(uuid.UUID) []dtovacancy.Pair); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dtovacancy.Pair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecommendation provides a mock function with given fields: start, limit, salary, spheres
func (_m *RepositoryVacancy) GetRecommendation(start int, limit int, salary float64, spheres []int) ([]models.Vacancy, error) {
	ret := _m.Called(start, limit, salary, spheres)

	var r0 []models.Vacancy
	if rf, ok := ret.Get(0).(func(int, int, float64, []int) []models.Vacancy); ok {
		r0 = rf(start, limit, salary, spheres)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Vacancy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, float64, []int) error); ok {
		r1 = rf(start, limit, salary, spheres)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVacancyById provides a mock function with given fields: _a0
func (_m *RepositoryVacancy) GetVacancyById(_a0 uuid.UUID) (*models.Vacancy, error) {
	ret := _m.Called(_a0)

	var r0 *models.Vacancy
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Vacancy); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Vacancy)
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

// GetVacancyList provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *RepositoryVacancy) GetVacancyList(_a0 uint, _a1 uint, _a2 uuid.UUID, _a3 int) ([]models.Vacancy, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 []models.Vacancy
	if rf, ok := ret.Get(0).(func(uint, uint, uuid.UUID, int) []models.Vacancy); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Vacancy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, uint, uuid.UUID, int) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchVacancies provides a mock function with given fields: _a0
func (_m *RepositoryVacancy) SearchVacancies(_a0 models.VacancySearchParams) ([]models.Vacancy, error) {
	ret := _m.Called(_a0)

	var r0 []models.Vacancy
	if rf, ok := ret.Get(0).(func(models.VacancySearchParams) []models.Vacancy); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Vacancy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.VacancySearchParams) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateVacancy provides a mock function with given fields: _a0
func (_m *RepositoryVacancy) UpdateVacancy(_a0 models.Vacancy) (*models.Vacancy, error) {
	ret := _m.Called(_a0)

	var r0 *models.Vacancy
	if rf, ok := ret.Get(0).(func(models.Vacancy) *models.Vacancy); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Vacancy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Vacancy) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
