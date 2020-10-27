package models

import (
	"github.com/google/uuid"
	"time"
)

type Resume struct {
	tableName struct{} `pg:"main.resume,discard_unknown_columns"`

	ID              uuid.UUID `pg:"resume_id,pk,type:uuid" json:"id"`
	UserID          uuid.UUID `pg:"cand_id, fk, type:uuid" json:"user_id"`
	Title           string    `pg:"title, notnull" json:"title"`
	SalaryMin       *int      `pg:"salary_min" json:"salary_min"`
	SalaryMax       *int      `pg:"salary_max" json:"salary_max"`
	Description     string    `pg:"description, notnull" json:"description"`
	Skills          string    `pg:"skills, notnull" json:"skills"`
	Gender          string    `pg:"gender, notnull" json:"gender"`
	EducationLevel  *string   `pg:"education_level" json:"education_level"`
	CareerLevel     *string   `pg:"career_level" json:"career_level"`
	Place           *string   `pg:"place" json:"place"`
	ExperienceMonth *int      `pg:"experience_month" json:"experience_month"`
	AreaSearch      *string   `pg:"area_search" json:"area_search"`
	DateCreate      time.Time `pg:"date_create" json:"date_create"`
}

type AdditionInResume struct {
	Education        []Education               `json:"education"`
	CustomExperience []ReqExperienceCustomComp `json:"custom_experience"`
}

type SearchResume struct {
	Title           string    `json:"title"`
	SalaryMin       int      `json:"salary_min"`
	SalaryMax       int      `json:"salary_max"`
	Gender          string    `json:"gender"`
	EducationLevel  string   `json:"education_level"`
	CareerLevel     string   `json:"career_level"`
	Place           string   `json:"place"`
	ExperienceMonth int      `json:"experience_month"`
	AreaSearch      string   `json:"area_search"`
}
