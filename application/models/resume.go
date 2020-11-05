package models

import (
	"github.com/google/uuid"
	"time"
)

type Resume struct {
	tableName struct{} `pg:"main.resume,discard_unknown_columns"`

	ID              uuid.UUID `pg:"resume_id,pk,type:uuid" json:"id" form:"id"`
	UserID          uuid.UUID `pg:"cand_id, fk, type:uuid" json:"cand_id" form:"user_id"`
	Title           string    `pg:"title, notnull" json:"title" form:"title"`
	SalaryMin       *int      `pg:"salary_min" json:"salary_min" form:"salary_min"`
	SalaryMax       *int      `pg:"salary_max" json:"salary_max" form:"salary_max"`
	Description     string    `pg:"description, notnull" json:"description" form:"description"`
	Skills          string    `pg:"skills, notnull" json:"skills" form:"skills"`
	Gender          string    `pg:"gender, notnull" json:"gender" form:"gender"`
	EducationLevel  *string   `pg:"education_level" json:"education_level" form:"education_level"`
	CareerLevel     *string   `pg:"career_level" json:"career_level" form:"career_level"`
	Place           *string   `pg:"place" json:"place" form:"place"`
	ExperienceMonth *int      `pg:"experience_month" json:"experience_month" form:"experience_month"`
	AreaSearch      *string   `pg:"area_search" json:"area_search" form:"area_search"`
	DateCreate      time.Time `pg:"date_create" json:"date_create" form:"date_create"`
}

type AdditionInResume struct {
	Education        []Education               `json:"education" form:"education"`
	CustomExperience []ReqExperienceCustomComp `json:"custom_experience" form:"custom_experience"`
}

type RespResume struct {
	Resume           Resume                 `json:"resume"`
	Educations       []Education            `json:"education"`
	CustomExperience []ExperienceCustomComp `json:"custom_experience"`
	IsFavorite       *uuid.UUID             `json:"is_favorite"`
}
type Resp struct {
	Resume []RespResume `json:"resume"`
}

type ResumeWithCandidate struct {
	tableName struct{} `pg:"main.resume,discard_unknown_columns"`

	ResumeID          uuid.UUID          `pg:"resume_id,pk,type:uuid" json:"resume_id"`
	CandID            uuid.UUID          `pg:"cand_id, fk, type:uuid" json:"cand_id"`
	Title             string             `pg:"title" json:"title"`
	Description       string             `pg:"description" json:"description"`
	Place             string             `pg:"place" json:"place"`
	AreaSearch        string             `pg:"area_search" json:"location"`
	CandidateWithUser *CandidateWithUser `pg:"rel:has-one"`
}

type BriefRespResume struct {
	ResumeID    uuid.UUID `json:"resume_id"`
	CandID      uuid.UUID `json:"cand_id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Place       string    `json:"place"`
	AreaSearch  string    `json:"location"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
}

type SearchResume struct {
	KeyWords        string   `json:"keywords"`
	SalaryMin       int      `json:"salary_min"`
	SalaryMax       int      `json:"salary_max"`
	Gender          []string `json:"gender"`
	EducationLevel  []string `json:"education_level"`
	CareerLevel     []string `json:"career_level"`
	ExperienceMonth []int    `json:"experience_month"`
	AreaSearch      []string `json:"area_search"`
}
