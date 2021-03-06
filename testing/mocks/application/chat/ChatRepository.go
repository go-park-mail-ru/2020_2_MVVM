// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	time "time"

	models "github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ChatRepository is an autogenerated mock type for the ChatRepository type
type ChatRepository struct {
	mock.Mock
}

// CreateChatAndTechMes provides a mock function with given fields: response
func (_m *ChatRepository) CreateChatAndTechMes(response models.Response) (*models.Chat, error) {
	ret := _m.Called(response)

	var r0 *models.Chat
	if rf, ok := ret.Get(0).(func(models.Response) *models.Chat); ok {
		r0 = rf(response)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Chat)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Response) error); ok {
		r1 = rf(response)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateMessage provides a mock function with given fields: mes, sender
func (_m *ChatRepository) CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error) {
	ret := _m.Called(mes, sender)

	var r0 *models.Message
	if rf, ok := ret.Get(0).(func(models.Message, uuid.UUID) *models.Message); ok {
		r0 = rf(mes, sender)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Message, uuid.UUID) error); ok {
		r1 = rf(mes, sender)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTechMesToUpdate provides a mock function with given fields: response
func (_m *ChatRepository) CreateTechMesToUpdate(response models.Response) (*models.Chat, error) {
	ret := _m.Called(response)

	var r0 *models.Chat
	if rf, ok := ret.Get(0).(func(models.Response) *models.Chat); ok {
		r0 = rf(response)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Chat)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Response) error); ok {
		r1 = rf(response)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTotalUnreadMes provides a mock function with given fields: userID, userType
func (_m *ChatRepository) GetTotalUnreadMes(userID uuid.UUID, userType string) (*uint, error) {
	ret := _m.Called(userID, userType)

	var r0 *uint
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) *uint); ok {
		r0 = rf(userID, userType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(userID, userType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListChats provides a mock function with given fields: userID, userType
func (_m *ChatRepository) ListChats(userID uuid.UUID, userType string) ([]models.ChatSummary, error) {
	ret := _m.Called(userID, userType)

	var r0 []models.ChatSummary
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) []models.ChatSummary); ok {
		r0 = rf(userID, userType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ChatSummary)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(userID, userType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkMessagesAsRead provides a mock function with given fields: chatID, utype, from, to, offset, limit
func (_m *ChatRepository) MarkMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error {
	ret := _m.Called(chatID, utype, from, to, offset, limit)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, *time.Time, *time.Time, *uint, *uint) error); ok {
		r0 = rf(chatID, utype, from, to, offset, limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MarkTechnicalMessagesAsRead provides a mock function with given fields: chatID, utype, from, to, offset, limit
func (_m *ChatRepository) MarkTechnicalMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error {
	ret := _m.Called(chatID, utype, from, to, offset, limit)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, *time.Time, *time.Time, *uint, *uint) error); ok {
		r0 = rf(chatID, utype, from, to, offset, limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MessagesForChat provides a mock function with given fields: chatID, from, to, offset, limit
func (_m *ChatRepository) MessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.MessageBrief, error) {
	ret := _m.Called(chatID, from, to, offset, limit)

	var r0 *[]models.MessageBrief
	if rf, ok := ret.Get(0).(func(uuid.UUID, *time.Time, *time.Time, *uint, *uint) *[]models.MessageBrief); ok {
		r0 = rf(chatID, from, to, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.MessageBrief)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, *time.Time, *time.Time, *uint, *uint) error); ok {
		r1 = rf(chatID, from, to, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OnlyUnreadMessagesForChat provides a mock function with given fields: chatID, userType
func (_m *ChatRepository) OnlyUnreadMessagesForChat(chatID uuid.UUID, userType string) (*[]models.MessageBrief, error) {
	ret := _m.Called(chatID, userType)

	var r0 *[]models.MessageBrief
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) *[]models.MessageBrief); ok {
		r0 = rf(chatID, userType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.MessageBrief)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(chatID, userType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OnlyUnreadTechnicalMessagesForChat provides a mock function with given fields: chatID, userType
func (_m *ChatRepository) OnlyUnreadTechnicalMessagesForChat(chatID uuid.UUID, userType string) (*[]models.TechMessageBrief, error) {
	ret := _m.Called(chatID, userType)

	var r0 *[]models.TechMessageBrief
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) *[]models.TechMessageBrief); ok {
		r0 = rf(chatID, userType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.TechMessageBrief)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(chatID, userType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TechnicalMessagesForChat provides a mock function with given fields: chatID, from, to, offset, limit
func (_m *ChatRepository) TechnicalMessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.TechMessageBrief, error) {
	ret := _m.Called(chatID, from, to, offset, limit)

	var r0 *[]models.TechMessageBrief
	if rf, ok := ret.Get(0).(func(uuid.UUID, *time.Time, *time.Time, *uint, *uint) *[]models.TechMessageBrief); ok {
		r0 = rf(chatID, from, to, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.TechMessageBrief)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, *time.Time, *time.Time, *uint, *uint) error); ok {
		r1 = rf(chatID, from, to, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
