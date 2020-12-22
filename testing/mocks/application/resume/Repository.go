// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package resume

import (
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	modelsresume "github.com/go-park-mail-ru/2020_2_MVVM.git/models/resume"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// AddFavorite provides a mock function with given fields: favoriteForEmpl
func (_m *Repository) AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoriteID, error) {
	ret := _m.Called(favoriteForEmpl)

	var r0 *models.FavoriteID
	if rf, ok := ret.Get(0).(func(models.FavoritesForEmpl) *models.FavoriteID); ok {
		r0 = rf(favoriteForEmpl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FavoriteID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.FavoritesForEmpl) error); ok {
		r1 = rf(favoriteForEmpl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: _a0
func (_m *Repository) Create(_a0 models.Resume) (*models.Resume, error) {
	ret := _m.Called(_a0)

	var r0 *models.Resume
	if rf, ok := ret.Get(0).(func(models.Resume) *models.Resume); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Resume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Resume) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: resId, candId
func (_m *Repository) Delete(resId uuid.UUID, candId uuid.UUID) error {
	ret := _m.Called(resId, candId)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(resId, candId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Drop provides a mock function with given fields: _a0
func (_m *Repository) Drop(_a0 models.Resume) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Resume) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllEmplFavoriteResume provides a mock function with given fields: emplID
func (_m *Repository) GetAllEmplFavoriteResume(emplID uuid.UUID) ([]models.Resume, error) {
	ret := _m.Called(emplID)

	var r0 []models.Resume
	if rf, ok := ret.Get(0).(func(uuid.UUID) []models.Resume); ok {
		r0 = rf(emplID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Resume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(emplID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllUserResume provides a mock function with given fields: userID
func (_m *Repository) GetAllUserResume(userID uuid.UUID) ([]models.Resume, error) {
	ret := _m.Called(userID)

	var r0 []models.Resume
	if rf, ok := ret.Get(0).(func(uuid.UUID) []models.Resume); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Resume)
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

// GetById provides a mock function with given fields: id
func (_m *Repository) GetById(id uuid.UUID) (*models.Resume, error) {
	ret := _m.Called(id)

	var r0 *models.Resume
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Resume); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Resume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByIdWithCand provides a mock function with given fields: id
func (_m *Repository) GetByIdWithCand(id uuid.UUID) (*models.Resume, error) {
	ret := _m.Called(id)

	var r0 *models.Resume
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Resume); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Resume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFavoriteByID provides a mock function with given fields: favoriteID
func (_m *Repository) GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error) {
	ret := _m.Called(favoriteID)

	var r0 *models.FavoritesForEmpl
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.FavoritesForEmpl); ok {
		r0 = rf(favoriteID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FavoritesForEmpl)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(favoriteID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFavoriteForResume provides a mock function with given fields: userID, resumeID
func (_m *Repository) GetFavoriteForResume(userID uuid.UUID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error) {
	ret := _m.Called(userID, resumeID)

	var r0 *models.FavoritesForEmpl
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) *models.FavoritesForEmpl); ok {
		r0 = rf(userID, resumeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FavoritesForEmpl)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(userID, resumeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: start, limit
func (_m *Repository) List(start uint, limit uint) ([]models.Resume, error) {
	ret := _m.Called(start, limit)

	var r0 []models.Resume
	if rf, ok := ret.Get(0).(func(uint, uint) []models.Resume); ok {
		r0 = rf(start, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Resume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, uint) error); ok {
		r1 = rf(start, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFavorite provides a mock function with given fields: favoriteForEmpl
func (_m *Repository) RemoveFavorite(favoriteForEmpl uuid.UUID) error {
	ret := _m.Called(favoriteForEmpl)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(favoriteForEmpl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: searchParams
func (_m *Repository) Search(searchParams *modelsresume.SearchParams) ([]models.Resume, error) {
	ret := _m.Called(searchParams)

	var r0 []models.Resume
	if rf, ok := ret.Get(0).(func(*modelsresume.SearchParams) []models.Resume); ok {
		r0 = rf(searchParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Resume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*modelsresume.SearchParams) error); ok {
		r1 = rf(searchParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0
func (_m *Repository) Update(_a0 models.Resume) (*models.Resume, error) {
	ret := _m.Called(_a0)

	var r0 *models.Resume
	if rf, ok := ret.Get(0).(func(models.Resume) *models.Resume); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Resume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Resume) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
