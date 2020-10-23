package models

import (
	"github.com/google/uuid"
)

type User struct {
	tableName struct{} `pg:"main.candidates,discard_unknown_columns"`

	ID            uuid.UUID `pg:"cand_id,pk,type:uuid" json:"id"`
	Nickname      string    `pg:"nickname,notnull" json:"nickname"`
	Name          string    `pg:"name,notnull" json:"name"`
	Surname       string    `pg:"surname" json:"surname"`
	Email         string    `pg:"email,notnull" json:"email"`
	PasswordHash  []byte    `pg:"password_hash,notnull" json:"-"`
	Phone         string    `pg:"phone" json:"phone"`
	AreaSearch    string    `pg:"area_search" json:"area_search"`
	SocialNetwork []string  `pg:"social_network" json:"social_network"`
	Avatar        string    `pg:"avatar" json:"avatar"`
}

type UserLogin struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// THIS DATA IS PUBLIC IN JWT
type JWTUserData struct {
	ID       uuid.UUID
	Nickname string
	Email    string
}
