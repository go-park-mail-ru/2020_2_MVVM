package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

type SearchParams struct {
	KeyWords        *string  `json:"keywords"`
	SalaryMin       *int     `json:"salary_min"`
	SalaryMax       *int     `json:"salary_max"`
	Gender          []string `json:"gender"`
	EducationLevel  []string `json:"education_level"`
	CareerLevel     []string `json:"career_level"`
	ExperienceMonth []int    `json:"experience_month"`
	AreaSearch      []string `json:"area_search"`
}

type Response struct {
	Resume           models.Resume                 `json:"resume"`
	User             models.User                   `json:"user"`
	Educations       []models.Education            `json:"education"`
	CustomExperience []models.ExperienceCustomComp `json:"custom_experience"`
	IsFavorite       *uuid.UUID                    `json:"is_favorite"`
}
