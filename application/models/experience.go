package models

import (
	"github.com/google/uuid"
	"time"
)

type ExperienceCustomComp struct {
	tableName struct{} `pg:"main.experience_in_custom_company,discard_unknown_columns"`

	ID              uuid.UUID  `pg:"exp_custom_id,pk,type:uuid" json:"id"`
	CandID          uuid.UUID  `pg:"cand_id,pk,type:uuid" json:"cand_id"`
	ResumeID        uuid.UUID  `pg:"resume_id,pk,type:uuid" json:"resume_id"`
	NameJob         string     `pg:"name_job, notnull" json:"name_job"`
	Position        *string    `pg:"position" json:"position"`
	Begin           time.Time  `pg:"begin,notnull" json:"begin"`
	Finish          *time.Time `pg:"finish" json:"finish"`
	Duties          *string    `pg:"duties" json:"duties"`
	ContinueToToday *bool      `pg:"continue_to_today" json:"continue_to_today"`
}

type ExperienceOfficialComp struct {
	tableName struct{} `pg:"main.experience_in_official_company,discard_unknown_columns"`

	ID              uuid.UUID  `pg:"exp_official_id,pk,type:uuid" json:"id"`
	CandID          uuid.UUID  `pg:"cand_id,pk,type:uuid" json:"cand_id"`
	ResumeID        uuid.UUID  `pg:"resume_id,pk,type:uuid" json:"resume_id"`
	CompanyID       uuid.UUID  `pg:"company_id,pk,type:uuid" json:"company_id"`
	Position        *string    `pg:"position" json:"position"`
	Begin           time.Time  `pg:"begin,notnull" json:"begin"`
	Finish          *time.Time `pg:"finish" json:"finish"`
	Duties          *string    `pg:"duties" json:"duties"`
	ContinueToToday *string    `pg:"continue_to_today" json:"continue_to_today"`
}

type ReqExperienceCustomComp struct {
	NameJob         string  `json:"name_job"`
	Position        *string `json:"position"`
	Begin           string  `json:"begin"`
	Finish          *string `json:"finish"`
	Duties          *string `json:"duties"`
	ContinueToToday bool   `json:"continue_to_today"`
}
