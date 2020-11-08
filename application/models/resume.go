package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Resume struct {
	tableName struct{}  `pg:"main.resume,discard_unknown_columns"`
	ResumeID  uuid.UUID `pg:"resume_id,pk,type:uuid" json:"id" form:"id"`

	CandID    uuid.UUID  `pg:"cand_id, fk, type:uuid" json:"cand_id" form:"cand_id"`
	Candidate *Candidate `pg:"rel:has-one"`

	Title                string                  `pg:"title, notnull" json:"title" form:"title"`
	SalaryMin            *int                    `pg:"salary_min" json:"salary_min" form:"salary_min"`
	SalaryMax            *int                    `pg:"salary_max" json:"salary_max" form:"salary_max"`
	Description          string                  `pg:"description, notnull" json:"description" form:"description"`
	Skills               string                  `pg:"skills, notnull" json:"skills" form:"skills"`
	Gender               string                  `pg:"gender, notnull" json:"gender" form:"gender"`
	EducationLevel       *string                 `pg:"education_level" json:"education_level" form:"education_level"`
	CareerLevel          *string                 `pg:"career_level" json:"career_level" form:"career_level"`
	Place                *string                 `pg:"place" json:"place" form:"place"`
	ExperienceMonth      *int                    `pg:"experience_month" json:"experience_month" form:"experience_month"`
	AreaSearch           *string                 `pg:"area_search" json:"area_search" form:"area_search"`
	DateCreate           time.Time               `pg:"date_create" json:"date_create" form:"date_create"`
	Education            []*Education            `pg:"rel:has-many" json:"education"`
	ExperienceCustomComp []*ExperienceCustomComp `pg:"rel:has-many" json:"custom_experience"`
}

// TODO Фронт сам решит, что ему надо. Надо вернуть полное резюме
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
