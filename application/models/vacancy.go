package models

import (
	"github.com/google/uuid"
)

type Vacancy struct {
	tableName       struct{}  `pg:"main.vacancy,discard_unknown_columns"`
	ID              uuid.UUID `pg:"vac_id,pk,type:uuid"`
	FK1             uuid.UUID `pg:"empl_id, fk, type:uuid"`
	FK2             uuid.UUID `pg:"comp_id, fk, type:uuid"`
	Title           string    `pg:"title,notnull" json:"title"`
	SalaryMin       int       `pg:"salary_min" json:"salary_min"`
	SalaryMax       int       `pg:"salary_max" json:"salary_max"`
	Description     string    `pg:"description,notnull" json:"description"`
	Requirements    string    `pg:"requirements" json:"requirements"`
	Duties          string    `pg:"duties" json:"duties"`
	Skills          string    `pg:"skills" json:"skills"`
	Spheres         string    `pg:"spheres" json:"spheres"`
	Employment      string    `pg:"employment" json:"employment"`
	WeekWorkHours   int       `pg:"week_work_hours" json:"week_work_hours"`
	ExperienceMonth string    `pg:"experience_month" json:"experience_work"`
	Location        string    `pg:"location" json:"location"`
	CareerLevel     string    `pg:"career_level" json:"career_level"`
	EducationLevel  string    `pg:"education_level" json:"education_level"`
	//DateCreate      time.Time `pg:"date_create" json:"date_create"`
}
