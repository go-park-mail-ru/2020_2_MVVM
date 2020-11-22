package models

import (
	"github.com/google/uuid"
)

type FavoritesForEmpl struct {
	FavoriteID uuid.UUID `gorm:"column:favorite_id" json:"favorite_id"`
	EmplID     uuid.UUID `gorm:"column:empl_id" json:"empl_id"`
	ResumeID   uuid.UUID `gorm:"column:resume_id" json:"resume_id"`
	//Resume     *Resume   `gorm:"foreignKey:resume_id;references:main.resume"`
}

type FavoritesForCand struct {
	tableName struct{} `gorm:"column:main.favorite_for_cand;discard_unknown_columns"`

	ID        uuid.UUID `gorm:"column:favorite_id;pk;type:uuid" json:"favorite_id"`
	CandID    uuid.UUID `gorm:"column:cand_id; fk; type:uuid" json:"cand_id"`
	VacancyID uuid.UUID `gorm:"column:vacancy_id; fk; type:uuid" json:"vacancy_id"`
}
