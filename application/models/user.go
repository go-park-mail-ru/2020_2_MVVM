package models

import (
	"github.com/google/uuid"
)

type User struct {
	tableName struct{} `pg:"main.users,discard_unknown_columns"`

	ID            uuid.UUID `pg:"user_id,pk,type:uuid" json:"id"`
	UserType      string    `pg:"user_type,notnull" json:"user_type"`
	Nickname      string    `pg:"nickname,notnull" json:"nickname"`
	Name          string    `pg:"name,notnull" json:"name"`
	Surname       string    `pg:"surname" json:"surname"`
	Email         string    `pg:"email,notnull" json:"email"`
	PasswordHash  []byte    `pg:"password_hash,notnull" json:"-"`
	Phone         *string   `pg:"phone" json:"phone"`
	SocialNetwork *string   `pg:"social_network" json:"social_network"`
}

type Employer struct {
	tableName struct{} `pg:"main.employers,discard_unknown_columns"`

	ID        uuid.UUID `pg:"empl_id,pk,type:uuid" json:"empl_id"`
	UserID    uuid.UUID `pg:"user_id,type:uuid" json:"user_id"`
	CompanyID uuid.UUID `pg:"comp_id,type:uuid" json:"comp_id"`
}

type Candidate struct {
	tableName struct{} `pg:"main.candidates,discard_unknown_columns"`

	ID     uuid.UUID `pg:"cand_id,pk,type:uuid" json:"cand_id"`
	UserID uuid.UUID `pg:"user_id,type:uuid" json:"user_id"`
}

type UserLogin struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CandidateWithUser struct {
	tableName struct{} `pg:"main.candidates,discard_unknown_columns"`

	CandID uuid.UUID `pg:"cand_id,pk,type:uuid" json:"cand_id"`
	UserID uuid.UUID `pg:"user_id,type:uuid" json:"user_id"`
	User   *User     `pg:"rel:has-one"`
}
