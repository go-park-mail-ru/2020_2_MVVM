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
	Sender     string    `gorm:"column:sender; type:uuid" json:"sender"`
	Message    string    `gorm:"column:message; type:uuid" json:"message"`
	IsRead     bool      `gorm:"column:is_read; type:uuid" json:"is_read"`
	DateCreate time.Time `gorm:"column:date_create; type:uuid" json:"date_create"`
}

func (r Message) TableName() string {
	return "main.message"
}

//easyjson:json
type ListMessage []Message

type BriefChat struct {
	ChatID     uuid.UUID `gorm:"column:chat_id; type:uuid" json:"chat_id"`
	Sender     string    `gorm:"column:sender; type:uuid" json:"sender"`
	Message    string    `gorm:"column:message; type:uuid" json:"message"`
	IsRead     bool      `gorm:"column:is_read; type:uuid" json:"is_read"`
	DateCreate time.Time `gorm:"column:date_create; type:uuid" json:"date_create"`
	Name       string    `gorm:"column:name" json:"name"`
	Surname    string    `gorm:"column:surname" json:"surname"`
	//UserType string `gorm:"column:date_create; type:uuid" json:"date_create"`
	PathToAvatar string `gorm:"column:path_to_avatar" json:"avatar"`
}
//easyjson:json
type ListBriefChat []BriefChat

type TechMessage struct {
	MessageID   uuid.UUID `gorm:"column:message_id;primaryKey;default:uuid_generate_v4()" json:"message_id"`
	ChatID      uuid.UUID `gorm:"column:chat_id; type:uuid" json:"chat_id"`
	ResumeID    uuid.UUID `json:"resume_id"`
	ResumeName  string    `json:"resume_name"`
	VacancyID   uuid.UUID `json:"vacancy_id"`
	VacancyName string    `json:"vacancy_name"`
	CompanyID   uuid.UUID `json:"company_id"`
	CompanyName string    `json:"company_name"`
	Status      string    `json:"status"`
}
