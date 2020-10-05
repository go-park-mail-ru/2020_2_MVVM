package models

import "github.com/google/uuid"

type Resume struct {
	tableName struct{} `pg:"main.resume,discard_unknown_columns"`

	ID          uuid.UUID `pg:"resume_id,pk,type:uuid"`
	FK          uuid.UUID `pg:"user_id, fk, type:uuid"`
	Title       string    `pg:"title,notnull" json:"title"`
	Salary      int       `pg:"salary" json:"salary"`
	Description string    `pg:"description" json:"description"`
	Skills      string    `pg:"skills" json:"skills"`
	Views       int       `pg:"views" json:"views"`
}
