package models

import "github.com/google/uuid"

type Resume struct {
	tableName struct{} `pg:"main.resume,discard_unknown_columns"`

	ID              uuid.UUID `pg:"resume_id,pk,type:uuid" json:"id"`
	UserID          uuid.UUID `pg:"user_id, fk, type:uuid" json:"user_id"`
	SalaryMin       *int       `pg:"salary_min" json:"salary_min"`
	SalaryMax       *int       `pg:"salary_max" json:"salary_max"`
	Description     *string    `pg:"description" json:"description"`
	Gender          *string    `pg:"gender" json:"gender"`
	Level           *string    `pg:"level" json:"level"`
	ExperienceMonth *int       `pg:"experience_month" json:"experience_month"`
	Education       *string    `pg:"education" json:"education"`
}
