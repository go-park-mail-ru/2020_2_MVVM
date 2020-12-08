package models

import (
	"github.com/google/uuid"
)

type Vacancy struct {
	ID              uuid.UUID `gorm:"column:vac_id; primaryKey;default:uuid_generate_v4()" json:"vac_id"`
	EmpID           uuid.UUID `gorm:"column:empl_id;type:uuid" json:"empl_id"`
	CompID          uuid.UUID `gorm:"column:comp_id;type:uuid" json:"comp_id"`
	Title           string    `gorm:"column:title;notnull" json:"title"`
	Gender          string    `gorm:"column:gender;default:\"male\"" json:"gender"`
	SalaryMin       int       `gorm:"column:salary_min" json:"salary_min"`
	SalaryMax       int       `gorm:"column:salary_max" json:"salary_max"`
	Description     string    `gorm:"column:description;notnull" json:"description"`
	Requirements    string    `gorm:"column:requirements" json:"requirements"`
	Duties          string    `gorm:"column:duties" json:"duties"`
	Skills          string    `gorm:"column:skills" json:"skills"`
	Sphere          int       `gorm:"column:sphere" json:"sphere"`
	Employment      string    `gorm:"column:employment" json:"employment"`
	ExperienceMonth int       `gorm:"column:experience_month" json:"experience_month"`
	AreaSearch      string    `gorm:"column:area_search" json:"area_search"`
	Location        string    `gorm:"column:location" json:"location"`
	CareerLevel     string    `gorm:"column:career_level" json:"career_level"`
	EducationLevel  string    `gorm:"column:education_level" json:"education_level"`
	DateCreate      string    `gorm:"column:date_create" json:"date_create"`
	EmpEmail        string    `gorm:"column:empl_email" json:"email"`
	EmpPhone        string    `gorm:"column:empl_phone" json:"phone"`
	Avatar          string    `gorm:"column:path_to_avatar" json:"avatar"`
}

type VacancySearchParams struct {
	KeyWords        string   `json:"keywords"`
	SalaryMin       int      `json:"salary_min"`
	SalaryMax       int      `json:"salary_max"`
	Gender          string   `json:"gender"`
	ExperienceMonth []int    `json:"experience_month"`
	Employment      []string `json:"employment"`
	EducationLevel  []string `json:"education_level"`
	CareerLevel     []string `json:"career_level"`
	Sphere          []int    `json:"sphere"`
	AreaSearch      []string `json:"area_search"`
	OrderBy         string   `json:"order_by"`
	ByAsc           bool     `json:"byAsc"`
	DaysFromNow     int      `json:"days_from_now"`
	StartDate       string   `json:"start_date"`
	KeywordsGeo     string   `json:"keywordsGeo"` //for main
}

func (v Vacancy) TableName() string {
	return "main.vacancy"
}

//easyjson:json
type ListVacancy []Vacancy
