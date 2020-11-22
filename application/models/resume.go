package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Resume struct {
	ResumeID  uuid.UUID `gorm:"column:resume_id" json:"id" form:"id" valid:"-"`

	CandID    uuid.UUID  `gorm:"foreignKey:cand_id" json:"cand_id" form:"cand_id" valid:"-"`
	Candidate *Candidate `gorm:"foreignKey:cand_id" valid:"-"`

	Title                string                  `gorm:"column:title; notnull" json:"title" form:"title" valid:"stringlength(4|128)~название резюме должно быть от 4 до 128 символов в длину." `
	SalaryMin            *int                    `gorm:"column:salary_min" json:"salary_min" form:"salary_min" valid:"-"`
	SalaryMax            *int                    `gorm:"column:salary_max" json:"salary_max" form:"salary_max" valid:"-"`
	Description          string                  `gorm:"column:description; notnull" json:"description" form:"description" valid:"-"`
	Skills               string                  `gorm:"column:skills; notnull" json:"skills" form:"skills" valid:"-"`
	Gender               string                  `gorm:"column:gender; notnull" json:"gender" form:"gender" valid:"alpha;stringlength(4|6)"`
	EducationLevel       *string                 `gorm:"column:education_level" json:"education_level" form:"education_level" valid:"-"`
	CareerLevel          *string                 `gorm:"column:career_level" json:"career_level" form:"career_level" valid:"-"`
	Place                *string                 `gorm:"column:place" json:"place" form:"place" valid:"-"`
	ExperienceMonth      *int                    `gorm:"column:experience_month" json:"experience_month" form:"experience_month" valid:"-"`
	AreaSearch           *string                 `gorm:"column:area_search" json:"area_search" form:"area_search" valid:"-"`
	DateCreate           time.Time               `gorm:"column:date_create" json:"date_create" form:"date_create" valid:"-"`
	Education            []*Education            `gorm:"foreignKey:ed_id;" json:"education" valid:"-"`
	ExperienceCustomComp []*ExperienceCustomComp `gorm:"foreignKey:exp_custom_id;" json:"custom_experience" valid:"-"`
	Avatar               string                  `gorm:"column:-" json:"avatar" valid:"-"`
}

func (r Resume) TableName() string {
	return "main.resume"
}


func (r *Resume) Brief() (*BriefResumeInfo, error) {
	if r.Candidate == nil || r.Candidate.User == nil {
		return nil, fmt.Errorf("failed to create brief resume description")
	}

	return &BriefResumeInfo{
		ResumeID:    r.ResumeID,
		CandID:      r.CandID,
		UserID:      r.Candidate.UserID,
		Title:       r.Title,
		Description: r.Description,
		Place:       r.Place,
		AreaSearch:  r.AreaSearch,
		Name:        r.Candidate.User.Name,
		Surname:     r.Candidate.User.Surname,
		Email:       r.Candidate.User.Email,
	}, nil
}

// TODO ВСЕГДА ИСПОЛЬЗОВАТЬ ТОЛЬКО ОДИН ID
type BriefResumeInfo struct {
	ResumeID    uuid.UUID `json:"resume_id"`
	CandID      uuid.UUID `json:"cand_id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Place       *string   `json:"place"`
	AreaSearch  *string   `json:"location"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
}
