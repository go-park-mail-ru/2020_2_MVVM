package models

import (
	"github.com/google/uuid"
	"time"
)

type Response struct {
	ID         uuid.UUID `gorm:"column:response_id;primaryKey;default:uuid_generate_v4()" json:"response_id"`
	ResumeID   uuid.UUID `gorm:"column:resume_id; type:uuid" json:"resume_id"`
	VacancyID  uuid.UUID `gorm:"column:vacancy_id; type:uuid" json:"vacancy_id"`
	Initial    string    `gorm:"column:initial" json:"initial"`
	Status     string    `gorm:"column:status" json:"status"`
	DateCreate time.Time `gorm:"column:date_create" json:"date_create"`
}

func (r Response) TableName() string {
	return "main.response"
}

type ResponseWithTitle struct {
	ResponseID  uuid.UUID `json:"response_id"`
	ResumeID    uuid.UUID `json:"resume_id"`
	ResumeName  string    `json:"resume_name"`
	CandName    string    `json:"cand_name"`
	CandSurname string    `json:"cand_surname"`
	VacancyID   uuid.UUID `json:"vacancy_id"`
	VacancyName string    `json:"vacancy_name"`
	CompanyID   uuid.UUID `json:"company_id"`
	CompanyName string    `json:"company_name"`
	Initial     string    `json:"initial"`
	Status      string    `json:"status"`
	DateCreate  time.Time `json:"date_create"`
}

//easyjson:json
type ListResponseWithTitle []ResponseWithTitle
