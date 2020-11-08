package models

import (
	"github.com/google/uuid"
)

type FavoritesForEmpl struct {
	tableName struct{} `pg:"main.favorite_for_empl,discard_unknown_columns"`

	FavoriteID uuid.UUID `pg:"favorite_id,pk,type:uuid" json:"favorite_id"`
	EmplID     uuid.UUID `pg:"empl_id, fk, type:uuid" json:"empl_id"`
	ResumeID   uuid.UUID `pg:"resume_id, fk, type:uuid" json:"resume_id"`
	Resume     *Resume   `pg:"has-one"`
}

type FavoritesForCand struct {
	tableName struct{} `pg:"main.favorite_for_cand,discard_unknown_columns"`

	ID        uuid.UUID `pg:"favorite_id,pk,type:uuid" json:"favorite_id"`
	CandID    uuid.UUID `pg:"cand_id, fk, type:uuid" json:"cand_id"`
	VacancyID uuid.UUID `pg:"vacancy_id, fk, type:uuid" json:"vacancy_id"`
}
