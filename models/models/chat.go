package models

import (
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ChatID     uuid.UUID `gorm:"column:chat_id;primaryKey;default:uuid_generate_v4()" json:"chat_id"`
	ResponseID uuid.UUID `gorm:"column:response_id; type:uuid" json:"response_id"`
	CandID     uuid.UUID `gorm:"column:user_id_cand; type:uuid" json:"cand_id"`
	EmplID     uuid.UUID `gorm:"column:user_id_empl; type:uuid" json:"empl_id"`
}

func (r Chat) TableName() string {
	return "main.chat"
}

type Message struct {
	MessageID  uuid.UUID `gorm:"column:message_id;primaryKey;default:uuid_generate_v4()" json:"message_id"`
	ChatID     uuid.UUID `gorm:"column:chat_id; type:uuid" json:"chat_id"`
	Sender     string    `gorm:"column:sender;" json:"sender"`
	Message    string    `gorm:"column:message;" json:"message"`
	IsRead     bool      `gorm:"column:is_read;" json:"is_read"`
	DateCreate time.Time `gorm:"column:date_create" json:"date_create"`
}

func (r Message) TableName() string {
	return "main.message"
}

type MessageBrief struct {
	Sender     string    `json:"sender"`
	Message    string    `json:"message"`
	IsRead     bool      `gorm:"column:is_read;" json:"is_read"`
	DateCreate time.Time `json:"date_create"`
}

type TechMessage struct {
	MessageID  uuid.UUID `gorm:"column:message_id;primaryKey;default:uuid_generate_v4()" json:"message_id"`
	ChatID     uuid.UUID `gorm:"column:chat_id; type:uuid" json:"chat_id"`
	ResponseID uuid.UUID `gorm:"column:response_id; type:uuid" json:"response_id"`
	DateCreate time.Time `gorm:"column:date_create" json:"date_create"`
}

func (r TechMessage) TableName() string {
	return "main.tech_message"
}

type TechMessageBrief struct {
	DateCreate time.Time `gorm:"column:date_create; type:uuid" json:"date_create"`
	//IsRead          bool      `gorm:"column:is_read;" json:"is_read"`
	ResumeID        uuid.UUID `gorm:"column:resume_id; type:uuid" json:"resume_id"`
	ResumeTitle     string    `gorm:"column:resume_title" json:"resume_title"`
	CompanyID       uuid.UUID `gorm:"column:company_id; type:uuid" json:"company_id"`
	CompanyName     string    `gorm:"column:company_name" json:"company_name"`
	VacancyID       uuid.UUID `gorm:"column:vacancy_id; type:uuid" json:"vacancy_id"`
	VacancyTitle    string    `gorm:"column:vacancy_title" json:"vacancy_title"`
	ResponseID      uuid.UUID `gorm:"column:response_id; type:uuid" json:"response_id"`
	ResponseInitial string    `gorm:"column:response_initial" json:"response_initial"`
	ResponseStatus  string    `gorm:"column:response_status" json:"response_status"`
}

type ChatHistory struct {
	ChatID            uuid.UUID          `json:"chat_id"`
	TechnicalMessages []TechMessageBrief `json:"technical_messages"`
	Dialog            []MessageBrief     `json:"dialog"`
}

type ChatSummary struct {
	ChatID      uuid.UUID   `json:"chat_id"`
	TotalUnread uint        `json:"total_unread"`
	Name        string      `json:"name"`
	Surname     string      `json:"surname"`
	Avatar      string      `json:"avatar"`
	Type        string      `json:"type"`
	Message     interface{} `json:"message"`
}

//easyjson:json
type ListChatSummary []ChatSummary
