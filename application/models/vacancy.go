package models

import "github.com/google/uuid"

type Vacancy struct {
	tableName          struct{}  `pg:"main.vacancy,discard_unknown_columns"`
	ID                 uuid.UUID `pg:"vacancy_id,pk,type:uuid"`
	//FK                 uuid.UUID `pg:"user_id, fk, type:uuid"`
	VacancyName        string    `pg:"vacancy_name,notnull" json:"vacancy_name"`
	CompanyName        string    `pg:"company_name,notnull" json:"company_name"`
	VacancyDescription string    `pg:"vacancy_description" json:"vacancy_description"`
	WorkExperience     string    `pg:"work_experience" json:"work_experience"`
	CompanyAddress     string    `pg:"company_address" json:"company_address"`
	Skills             string    `pg:"skills" json:"skills"`
	Salary             int       `pg:"salary" json:"salary"`
}
