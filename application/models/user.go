package models

import (
	"github.com/google/uuid"
)

type User struct {
	tableName struct{} `pg:"main.users,discard_unknown_columns"`

	ID            uuid.UUID `pg:"user_id,pk,type:uuid" json:"id"`
	UserType      string    `pg:"user_type,notnull" json:"user_type"`
	Name          string    `pg:"name,notnull" json:"name"`
	Surname       string    `pg:"surname" json:"surname"`
	Email         string    `pg:"email,notnull" json:"email"`
	PasswordHash  []byte    `pg:"password_hash,notnull" json:"-"`
	Phone         *string   `pg:"phone" json:"phone"`
	SocialNetwork *string   `pg:"social_network" json:"social_network"`
}

type Employer struct {
	tableName struct{} `pg:"main.employers,discard_unknown_columns"`

	ID        uuid.UUID           `pg:"empl_id,pk,type:uuid" json:"empl_id"`
	UserID    uuid.UUID           `pg:"user_id,type:uuid" json:"user_id"`
	CompanyID uuid.UUID           `pg:"comp_id,type:uuid" json:"comp_id"`
	Favorites []*FavoritesForEmpl `pg:"has-many"`
}

type Candidate struct {
	tableName struct{} `pg:"main.candidates,discard_unknown_columns"`

	ID     uuid.UUID `pg:"cand_id,pk,type:uuid" json:"cand_id"`
	UserID uuid.UUID `pg:"user_id,type:uuid" json:"user_id"`
	User   *User     `pg:"rel:has-one"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required" valid:"email"`
	Password string `json:"password" binding:"required" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
}
