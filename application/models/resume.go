package models

import (
	"github.com/google/uuid"
	"time"
)

type Resume struct {
	tableName struct{} `pg:"main.resume,discard_unknown_columns"`

	ID              uuid.UUID `pg:"resume_id,pk,type:uuid" json:"id"`
	UserID          uuid.UUID `pg:"cand_id, fk, type:uuid" json:"user_id"`
	Title           *string   `pg:"title, notnull" json:"title"`
	SalaryMin       *int      `pg:"salary_min" json:"salary_min"`
	SalaryMax       *int      `pg:"salary_max" json:"salary_max"`
	AboutMe         *string   `pg:"about_me, notnull" json:"about_me"`
	Gender          *string   `pg:"gender, notnull" json:"gender"`
	EducationLevel  *string   `pg:"education_level" json:"education_level"`
	CareerLevel     *string   `pg:"career_level" json:"career_level"`
	ExperienceMonth *int      `pg:"experience_month" json:"experience_month"`
	DateCreate      time.Time `pg:"date_create" json:"date_create"`
}

type ReqResume struct {
	Educations       ListEducations          `json:"educations"`
	CustomExperience ListReqCustomExperience `json:"custom_experience"`
}
