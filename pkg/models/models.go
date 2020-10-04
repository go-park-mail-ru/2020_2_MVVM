package models

import (
	"github.com/google/uuid"
)

type User struct {
	tableName struct{} `pg:"main.users,discard_unknown_columns"`

	ID      uuid.UUID `pg:"user_id,pk,type:uuid"`
	Name    string    `pg:"name,notnull" json:"name"`
	Surname string    `pg:"surname,notnull" json:"surname"`
}
