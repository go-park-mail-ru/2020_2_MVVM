package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"column:user_id;default:uuid_generate_v4()" json:"id"`
	UserType      string    `gorm:"column:user_type;notnull" json:"user_type"`
	Name          string    `gorm:"column:name;notnull" json:"name"`
	Surname       string    `gorm:"column:surname" json:"surname"`
	Email         string    `gorm:"column:email;notnull" json:"email"`
	PasswordHash  []byte    `gorm:"column:password_hash;notnull" json:"-"`
	Phone         *string   `gorm:"column:phone" json:"phone"`
	SocialNetwork *string   `gorm:"column:social_network" json:"social_network"`
}

func (e Employer) TableName() string {
	return "main.employers"
}

func (u User) TableName() string {
	return "main.users"
}

type Employer struct {
	ID        uuid.UUID           `gorm:"column:empl_id;primaryKey;default:uuid_generate_v4()" json:"empl_id"`
	UserID    uuid.UUID           `gorm:"column:user_id;foreignKey:user_id" json:"user_id"`
	CompanyID uuid.UUID           `gorm:"column:comp_id;foreignKey:comp_id" json:"comp_id"`
	Favorites []FavoritesForEmpl  `gorm:"foreignKey:EmplID"`
}

type Candidate struct {
	ID     uuid.UUID `gorm:"primaryKey;column:cand_id;default:uuid_generate_v4()" json:"cand_id"`
	UserID uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	User   User
}

func (c Candidate) TableName() string {
	return "main.candidates"
}

type UserLogin struct {
	Email    string `json:"email" binding:"required" valid:"email"`
	Password string `json:"password" binding:"required" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
}
