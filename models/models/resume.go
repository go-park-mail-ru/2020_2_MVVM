package models

import (
	"github.com/google/uuid"
)

type Resume struct {
	ResumeID uuid.UUID `gorm:"column:resume_id; primaryKey; default:uuid_generate_v4()" json:"id" form:"id" valid:"-"`

	CandID               uuid.UUID              `gorm:"column:cand_id" json:"cand_id" form:"cand_id" valid:"-"`
	Candidate            Candidate              `gorm:"foreignKey:CandID;save_associations:false" valid:"-" json:"-"`
	Title                string                 `gorm:"column:title; notnull" json:"title" binding:"required" valid:"required, stringlength(4|128)~название резюме должно быть от 4 до 128 символов в длину." `
	SalaryMin            *int                   `gorm:"column:salary_min" json:"salary_min" valid:"-"`
	SalaryMax            *int                   `gorm:"column:salary_max" json:"salary_max" valid:"-"`
	Sphere               *int                   `gorm:"column:sphere" json:"sphere" valid:"numeric~сфера деятельности должна содержать только код"`
	Description          string                 `gorm:"column:description; notnull" json:"description" binding:"required" valid:"required"`
	Skills               string                 `gorm:"column:skills; notnull" json:"skills" binding:"required" valid:"required"`
	Gender               string                 `gorm:"column:gender;" json:"gender" binding:"required" valid:"alpha;stringlength(4|6)"`
	EducationLevel       *string                `gorm:"column:education_level" json:"education_level" valid:"-"`
	CareerLevel          *string                `gorm:"column:career_level" json:"career_level" valid:"-"`
	Place                *string                `gorm:"column:place" json:"place" valid:"-"`
	ExperienceMonth      *int                   `gorm:"column:experience_month" json:"experience_month" valid:"-"`
	AreaSearch           *string                `gorm:"column:area_search" json:"area_search" valid:"-"`
	DateCreate           string                 `gorm:"column:date_create" json:"date_create" valid:"-"`
	Education            []Education            `gorm:"-;" json:"education" valid:"-"`
	ExperienceCustomComp []ExperienceCustomComp `gorm:"foreignKey:ResumeID" json:"custom_experience" valid:"-"`
	Avatar               string                 `gorm:"column:path_to_avatar" json:"avatar" valid:"-"`
	CandName             string                 `json:"cand_name" valid:"utfletter~имя должно содержать только буквы,stringlength(3|25)~длина имени должна быть от 3 до 25 символов., required"`
	CandSurname          string                 `json:"cand_surname" valid:"utfletter~фамилия должна содержать только буквы,stringlength(3|25)~длина фамилии должна быть от 3 до 25 символов., required"`
	CandEmail            string                 `json:"cand_email" valid:"email, required"`
}

func (r Resume) TableName() string {
	return "main.resume"
}

func (r *Resume) Brief() (*BriefResumeInfo, error) {
	//if r.Candidate == nil || r.Candidate.User == nil {
	//	return nil, fmt.Errorf("failed to create brief resume description")
	//}

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
		Avatar:      r.Avatar,
	}, nil
}

// TODO ВСЕГДА ИСПОЛЬЗОВАТЬ ТОЛЬКО ОДИН ID
type BriefResumeInfo struct {
	Avatar      string    `json:"avatar"`
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

//easyjson:json
type ListBriefResumeInfo []BriefResumeInfo

type LinkToPdf struct {
	Link string `json:"link_to_pdf"`
}
