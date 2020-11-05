package models

import (
	"github.com/google/uuid"
	"time"
)

type Response struct {
	tableName struct{} `pg:"main.response,discard_unknown_columns"`

	ID         uuid.UUID `pg:"response_id,pk,type:uuid" json:"response_id"`
	ResumeID   uuid.UUID `pg:"resume_id, fk, type:uuid" json:"resume_id"`
	VacancyID  uuid.UUID `pg:"vacancy_id, fk, type:uuid" json:"vacancy_id"`
	Initial    string    `pg:"initial" json:"initial"`
	Status     string    `pg:"isapply" json:"status"`
	DateCreate time.Time `pg:"date_create" json:"date_create"`
}

type ResponseWithTitle struct {
	ID          uuid.UUID `json:"response_id"`
	ResumeID    uuid.UUID `json:"resume_id"`
	ResumeName  string    `json:"resume_name"`
	VacancyID   uuid.UUID `json:"vacancy_id"`
	VacancyName string    `json:"vacancy_name"`
	CompanyName string    `json:"company_name"`
	Initial     string    `json:"initial"`
	Status      string    `json:"status"`
	DateCreate  time.Time `json:"date_create"`
}
