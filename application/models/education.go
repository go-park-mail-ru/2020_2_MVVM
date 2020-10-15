package models

import (
	"github.com/google/uuid"
	"time"
)

type Education struct {
	tableName struct{} `pg:"main.education,discard_unknown_columns"`

	EdId        uuid.UUID `pg:"ed_id,pk,type:uuid" json:"id"`
	CandId      uuid.UUID `pg:"cand_id,fk,type:uuid" json:"cand_id"`
	University  *string    `pg:"university" json:"university"`
	Level       *string    `pg:"level" json:"level"`
	Begin       time.Time `pg:"begin" json:"begin"`
	Finish      time.Time `pg:"finish" json:"finish"`
	Department  *string    `pg:"department" json:"department"`
	Description *string    `pg:"description" json:"description"`
}
