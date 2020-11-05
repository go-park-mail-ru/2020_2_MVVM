package models

import (
	"github.com/google/uuid"
)

type Vacancy struct {
	tableName       struct{}  `pg:"main.vacancy,discard_unknown_columns"`
	ID              uuid.UUID `pg:"vac_id,pk,type:uuid" json:"vac_id"`
	EmpID           uuid.UUID `pg:"empl_id, fk, type:uuid" json:"empl_id"`
	CompID          uuid.UUID `pg:"comp_id, fk, type:uuid" json:"comp_id"`
	Title           string    `pg:"title,notnull" json:"title"`
	SalaryMin       int       `pg:"salary_min" json:"salary_min"`
	SalaryMax       int       `pg:"salary_max" json:"salary_max"`
	Description     string    `pg:"description,notnull" json:"description"`
	Requirements    string    `pg:"requirements" json:"requirements"`
	Duties          string    `pg:"duties" json:"duties"`
	Skills          string    `pg:"skills" json:"skills"`
	Sphere          int       `pg:"sphere" json:"sphere"`
	Employment      string    `pg:"employment" json:"employment"`
	ExperienceMonth string    `pg:"experience_month" json:"experience_month"`
	Location        string    `pg:"location" json:"location"`
	CareerLevel     string    `pg:"career_level" json:"career_level"`
	EducationLevel  string    `pg:"education_level" json:"education_level"`
	DateCreate      string    `pg:"date_create" json:"date_create"`
	EmpEmail        string    `pg:"empl_email" json:"email"`
	EmpPhone        string    `pg:"empl_phone" json:"phone"`
}

type VacancySearchParams struct {
	KeyWords        string   `json:"keywords"`
	SalaryMin       int      `json:"salary_min"`
	SalaryMax       int      `json:"salary_max"`
	ExperienceMonth []int    `json:"experience_month"`
	Employment      []string `json:"employment"`
	EducationLevel  []string `json:"education_level"`
	CareerLevel     []string `json:"career_level"`
	Spheres         []int    `json:"spheres"`
	Location        []string `json:"location"`
	OrderBy         string   `json:"order_by"`
	ByAsc           bool     `json:"byAsc"`
	DaysFromNow     int      `json:"days_from_now"`
	StartDate       string   `json:"start_date"`
}
