package models

import (
	"github.com/google/uuid"
)

type OfficialCompany struct {
	tableName struct{} `pg:"main.official_company,discard_unknown_columns"`

	ID         uuid.UUID `pg:"comp_id,pk,type:uuid" json:"id"`
	Name       string    `pg:"name,notnull" json:"name"`
	Sphere     []string  `pg:"sphere,notnull" json:"sphere"`
	Location   string    `pg:"location" json:"location"`
	Link       string    `pg:"link" json:"link"`
	Phone      string    `pg:"phone" json:"phone"`
	Avatar     string    `pg:"avatar" json:"avatar"`
}

type CustomCompany struct {
	tableName struct{} `pg:"main.custom_company,discard_unknown_columns"`

	ID         uuid.UUID `pg:"company_id,pk,type:uuid" json:"id"`
	Name       string    `pg:"name,notnull" json:"name"`
	Sphere     []string  `pg:"sphere,notnull" json:"sphere"`
	Location   *string    `pg:"location" json:"location"`
}


type ReqCustomCompany struct {
	CompanyName string   `json:"name" binding:"required"`
	Location    *string   `json:"location"`
	Sphere      []string  `json:"sphere"`
}