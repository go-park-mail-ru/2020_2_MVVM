package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Resume struct {
	tableName struct{}  `pg:"main.resume,discard_unknown_columns"`
	ResumeID  uuid.UUID `pg:"resume_id,pk,type:uuid" json:"id" form:"id" valid:"-"`

	CandID    uuid.UUID  `pg:"cand_id, fk, type:uuid" json:"cand_id" form:"cand_id" valid:"-"`
	Candidate *Candidate `pg:"rel:has-one" valid:"-"`

	Title                string                  `pg:"title, notnull" json:"title" form:"title" valid:"utfletternum~название резюме может содержать только буквы и цифры.,stringlength(4|128)~название резюме должно быть от 4 до 128 символов в длину." `
	SalaryMin            *int                    `pg:"salary_min" json:"salary_min" form:"salary_min" valid:"-"`
	SalaryMax            *int                    `pg:"salary_max" json:"salary_max" form:"salary_max" valid:"-"`
	Description          string                  `pg:"description, notnull" json:"description" form:"description" valid:"-"`
	Skills               string                  `pg:"skills, notnull" json:"skills" form:"skills" valid:"-"`
	Gender               string                  `pg:"gender, notnull" json:"gender" form:"gender" valid:"alpha,stringlength(4|6)"`
	EducationLevel       *string                 `pg:"education_level" json:"education_level" form:"education_level" valid:"-"`
	CareerLevel          *string                 `pg:"career_level" json:"career_level" form:"career_level" valid:"-"`
	Place                *string                 `pg:"place" json:"place" form:"place" valid:"-"`
	ExperienceMonth      *int                    `pg:"experience_month" json:"experience_month" form:"experience_month" valid:"-"`
	AreaSearch           *string                 `pg:"area_search" json:"area_search" form:"area_search" valid:"-"`
	DateCreate           time.Time               `pg:"date_create" json:"date_create" form:"date_create" valid:"-"`
	Education            []*Education            `pg:"rel:has-many" json:"education" valid:"-"`
	ExperienceCustomComp []*ExperienceCustomComp `pg:"rel:has-many" json:"custom_experience" valid:"-"`
	Avatar               string                  `pg:"-" json:"avatar" valid:"-"`
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
