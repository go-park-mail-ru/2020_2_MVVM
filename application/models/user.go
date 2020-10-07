package models

import (
	"github.com/google/uuid"
)

type User struct {
	tableName struct{} `pg:"main.users,discard_unknown_columns"`

	ID           uuid.UUID `pg:"user_id,pk,type:uuid" json:"id"`
	Nickname     string    `pg:"nickname,notnull" json:"nickname"`
	Name         string    `pg:"name,notnull" json:"name"`
	Surname      string    `pg:"surname" json:"surname"`
	Email        string    `pg:"email,notnull" json:"email"`
	PasswordHash []byte    `pg:"password_hash,notnull" json:"-"`
}

// THIS DATA IS PUBLIC IN JWT
type JWTUserData struct {
	ID       uuid.UUID
	Nickname string
	Email    string
}
