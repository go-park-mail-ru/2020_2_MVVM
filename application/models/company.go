package models

import (
	"github.com/google/uuid"
)

type OfficialCompany struct {
	tableName struct{} `pg:"main.official_companies,discard_unknown_columns"`

	ID          uuid.UUID `pg:"comp_id,pk,type:uuid" json:"id"`
	Name        string    `pg:"name,notnull" json:"name"`
	Sphere      []string  `pg:"sphere,notnull" json:"sphere"`
	Description string    `pg:"description,notnull" json:"description"`
	Location    string    `pg:"location" json:"location"`
	Link        string    `pg:"link" json:"link"`
	VacCount    int       `pg:"count_vacancy" json:"vac_count"`
}

type CompanySearchParams struct {
	KeyWords  string   `json:"keywords"`
	Spheres   []string `json:"spheres"`
	OrderBy   string   `json:"order_by"`
	ByAsc     bool     `json:"byAsc"`
	VacCount  int      `json:"vac_count"`
}

//dont work
type CustomCompany struct {
}
type ReqCustomCompany struct {
}
