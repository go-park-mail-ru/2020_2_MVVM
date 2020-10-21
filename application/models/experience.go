package models

import (
	"github.com/google/uuid"
	"time"
)

type ExperienceCustomComp struct {
	tableName struct{} `pg:"main.experience_in_custom_company,discard_unknown_columns"`

	ID          uuid.UUID `pg:"exp_custom_id,pk,type:uuid" json:"id"`
	CandID      uuid.UUID `pg:"cand_id,pk,type:uuid" json:"cand_id"`
	ResumeID    uuid.UUID `pg:"resume_id,pk,type:uuid" json:"resume_id"`
	CompanyID   uuid.UUID `pg:"company_id,pk,type:uuid" json:"company_id"`
	Position    *string    `pg:"position,notnull" json:"position"`
	Begin       time.Time `pg:"begin,notnull" json:"begin"`
	Finish      time.Time `pg:"finish,notnull" json:"finish"`
	Description *string    `pg:"description" json:"description"`
}


type ExperienceOfficialComp struct {
	tableName struct{} `pg:"main.experience_in_official_company,discard_unknown_columns"`

	ID          uuid.UUID `pg:"exp_official_id,pk,type:uuid" json:"id"`
	CandID      uuid.UUID `pg:"cand_id,pk,type:uuid" json:"cand_id"`
	ResumeID    uuid.UUID `pg:"resume_id,pk,type:uuid" json:"resume_id"`
	CompanyID   uuid.UUID `pg:"company_id,pk,type:uuid" json:"company_id"`
	Position    string    `pg:"position,notnull" json:"position"`
	Begin       time.Time `pg:"begin,notnull" json:"begin"`
	Finish      time.Time `pg:"finish,notnull" json:"finish"`
	Description string    `pg:"description" json:"description"`
}

type CustomExperienceWithCompanies struct {
	CompanyName string   `json:"name" binding:"required"`
	Location    *string   `json:"location"`
	Sphere      []string  `json:"sphere"`
	Position    *string   `json:"position"`
	Begin       time.Time `json:"begin"`
	Finish      time.Time `json:"finish"`
	Description *string   `json:"description"`
}

type ListReqCustomExperience struct {
	ListReqCustomExperience []CustomExperienceWithCompanies `json:"custom_experience"`
}
