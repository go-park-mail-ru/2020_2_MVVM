package models

import (
	"github.com/google/uuid"
	"time"
)

type ExperienceCustomComp struct {
	tableName struct{} `pg:"main.experience_in_custom_company,discard_unknown_columns"`

	ID              uuid.UUID  `pg:"exp_custom_id,pk,type:uuid" json:"id" form:"id"`
	CandID          uuid.UUID  `pg:"cand_id,pk,type:uuid" json:"cand_id" form:"cand_id"`
	ResumeID        uuid.UUID  `pg:"resume_id,pk,type:uuid" json:"resume_id" form:"resume_id"`
	NameJob         string     `pg:"name_job, notnull" json:"name_job" form:"name_job"`
	Position        *string    `pg:"position" json:"position" form:"position"`
	Begin           time.Time  `pg:"begin,notnull" json:"begin" form:"begin"`
	Finish          *time.Time `pg:"finish" json:"finish" form:"finish"`
	Duties          *string    `pg:"duties" json:"duties" form:"duties"`
	ContinueToToday *bool      `pg:"continue_to_today" json:"continue_to_today" form:"continue_to_today"`
}

type ExperienceOfficialComp struct {
	tableName struct{} `pg:"main.experience_in_official_company,discard_unknown_columns"`

	ID              uuid.UUID  `pg:"exp_official_id,pk,type:uuid" json:"id" form:"id"`
	CandID          uuid.UUID  `pg:"cand_id,pk,type:uuid" json:"cand_id" form:"cand_id"`
	ResumeID        uuid.UUID  `pg:"resume_id,pk,type:uuid" json:"resume_id" form:"resume_id"`
	CompanyID       uuid.UUID  `pg:"company_id,pk,type:uuid" json:"company_id" form:"company_id"`
	Position        *string    `pg:"position" json:"position" form:"position"`
	Begin           time.Time  `pg:"begin,notnull" json:"begin" form:"begin"`
	Finish          *time.Time `pg:"finish" json:"finish" form:"finish"`
	Duties          *string    `pg:"duties" json:"duties" form:"duties"`
	ContinueToToday *string    `pg:"continue_to_today" json:"continue_to_today" form:"continue_to_today"`
}

type ReqExperienceCustomComp struct {
	NameJob         string  `json:"name_job" form:"name_job"`
	Position        *string `json:"position" form:"position"`
	Begin           string  `json:"begin" form:"begin"`
	Finish          *string `json:"finish" form:"finish"`
	Duties          *string `json:"duties" form:"duties"`
	ContinueToToday bool   `json:"continue_to_today" form:"continue_to_today"`
}
