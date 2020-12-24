package models

import (
	"github.com/google/uuid"
	"time"
)

type ExperienceCustomComp struct {
	ID              uuid.UUID  `gorm:"column:exp_custom_id; primaryKey; default:uuid_generate_v4()" json:"id" form:"id"`
	CandID          uuid.UUID  `gorm:"column:cand_id" json:"cand_id" form:"cand_id"`
	ResumeID        uuid.UUID  `gorm:"column:resume_id" json:"resume_id" form:"resume_id"`
	NameJob         string     `gorm:"column:name_job; notnull" json:"name_job" form:"name_job"`
	Position        *string    `gorm:"column:position" json:"position" form:"position"`
	Begin           time.Time  `gorm:"column:begin;notnull" json:"begin" form:"begin"`
	Finish          *time.Time `gorm:"column:finish" json:"finish" form:"finish"`
	Duties          *string    `gorm:"column:duties" json:"duties" form:"duties"`
	ContinueToToday *bool      `gorm:"column:continue_to_today" json:"continue_to_today" form:"continue_to_today"`
}

func (r ExperienceCustomComp) TableName() string {
	return "main.experience_in_custom_company"
}

type ExperienceOfficialComp struct {
	ID              uuid.UUID  `gorm:"column:exp_official_id,pk,type:uuid" json:"id" form:"id"`
	CandID          uuid.UUID  `gorm:"column:cand_id,pk,type:uuid" json:"cand_id" form:"cand_id"`
	ResumeID        uuid.UUID  `gorm:"column:resume_id,pk,type:uuid" json:"resume_id" form:"resume_id"`
	CompanyID       uuid.UUID  `gorm:"column:company_id,pk,type:uuid" json:"company_id" form:"company_id"`
	Position        *string    `gorm:"column:position" json:"position" form:"position"`
	Begin           time.Time  `gorm:"column:begin,notnull" json:"begin" form:"begin"`
	Finish          *time.Time `gorm:"column:finish" json:"finish" form:"finish"`
	Duties          *string    `gorm:"column:duties" json:"duties" form:"duties"`
	ContinueToToday *string    `gorm:"column:continue_to_today" json:"continue_to_today" form:"continue_to_today"`
}

type ReqExperienceCustomComp struct {
	NameJob         string  `json:"name_job" form:"name_job"`
	Position        *string `json:"position" form:"position"`
	Begin           string  `json:"begin" form:"begin"`
	Finish          *string `json:"finish" form:"finish"`
	Duties          *string `json:"duties" form:"duties"`
	ContinueToToday bool   `json:"continue_to_today" form:"continue_to_today"`
}
