package models

import (
	"github.com/google/uuid"
	"time"
)

type Education struct {
	tableName struct{} `pg:"main.education,discard_unknown_columns"`

	EdId        uuid.UUID  `pg:"ed_id,pk,type:uuid" json:"id" form:"id"`
	CandID      uuid.UUID  `pg:"cand_id,fk,type:uuid" json:"cand_id" form:"cand_id"`
	ResumeId    uuid.UUID  `pg:"resume_id,fk,type:uuid" json:"resume_id" form:"resume_id"`
	University  string     `pg:"university, notnull" json:"university" form:"university"`
	Level       *string    `pg:"level" json:"level" form:"level"`
	Begin       *time.Time `pg:"begin" json:"begin" form:"begin"`
	Finish      time.Time  `pg:"finish, notnull" json:"finish" form:"finish"`
	Department  *string    `pg:"department" json:"department" form:"department"`
	Description *string    `pg:"description" json:"description" form:"description"`
}
