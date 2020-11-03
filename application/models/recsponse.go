package models

import (
	"github.com/google/uuid"
	"time"
)

type Response struct {
	tableName struct{} `pg:"main.respomse,discard_unknown_columns"`

	ID         uuid.UUID `pg:"response_id,pk,type:uuid" json:"response_id"`
	ResumeID   uuid.UUID `pg:"resume_id, fk, type:uuid" json:"resume_id"`
	VacancyID  uuid.UUID `pg:"vacancy_id, fk, type:uuid" json:"vacancy_id"`
	Initial    string    `pg:"initial" json:"initial"`
	Status     string    `pg:"isApply" json:"status"`
	DateCreate time.Time `pg:"date_create" json:"date_create"`
}
