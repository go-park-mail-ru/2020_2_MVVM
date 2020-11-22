package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OfficialCompany struct {
	ID          uuid.UUID     `gorm:"column:comp_id;default:uuid_generate_v4()" json:"id"`
	Name        string        `gorm:"column:name;notnull" json:"name"`
	Spheres     pq.Int64Array `gorm:"column:spheres" json:"spheres"`
	Description string        `gorm:"column:description;notnull" json:"description"`
	AreaSearch  string        `gorm:"column:area_search" json:"area_search"`
	Link        string        `gorm:"column:link" json:"link"`
	VacCount    int           `gorm:"column:count_vacancy" json:"vac_count"`
}


func (c OfficialCompany) TableName() string {
	return "main.official_companies"
}

type CompanySearchParams struct {
	KeyWords   string   `json:"keywords"`
	AreaSearch []string `json:"area_search"`
	Spheres    []int    `json:"spheres"`
	OrderBy    string   `json:"order_by"`
	ByAsc      bool     `json:"byAsc"`
	VacCount   int      `json:"vac_count"`
}

//dont work
type CustomCompany struct {
}
type ReqCustomCompany struct {
}
