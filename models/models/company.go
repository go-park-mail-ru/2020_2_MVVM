package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OfficialCompany struct {
	ID          uuid.UUID     `gorm:"column:comp_id;default:uuid_generate_v4()" json:"id"`
	Name        string        `gorm:"column:name;notnull" json:"name"`
	Spheres     pq.Int64Array `gorm:"column:spheres" json:"sphere"`
	Description string        `gorm:"column:description;notnull" json:"description"`
	AreaSearch  string        `gorm:"column:area_search" json:"area_search"`
	Link        string        `gorm:"column:link" json:"link"`
	VacCount    int           `gorm:"column:count_vacancy" json:"vac_count"`
	Avatar      string        `gorm:"column:path_to_avatar" json:"avatar"`
}

func (v OfficialCompany) TableName() string {
	return "main.official_companies"
}

type CompanySearchParams struct {
	KeyWords   string   `json:"keywords"`
	AreaSearch []string `json:"area_search"`
	Sphere     []int    `json:"sphere"`
	OrderBy    string   `json:"order_by"`
	ByAsc      bool     `json:"byAsc"`
	VacCount   int      `json:"vac_count"`
}

type Resp struct {
	Company *OfficialCompany `json:"company"`
}

type RespList struct {
	Companies []OfficialCompany `json:"companyList"`
}

type ReqComp struct {
	Name        string `json:"name" binding:"required" valid:"stringlength(4|30)~название компании должно быть от 4 до 30 символов."`
	Description string `json:"description" binding:"required" valid:"-"`
	Spheres     []int  `json:"sphere" valid:"-"`
	AreaSearch  string `json:"area_search" valid:"stringlength(4|128)~длина названия региона от 4 до 128 смиволов"`
	Link        string `json:"link" valid:"url~неверный формат ссылки"`
	Avatar      string `json:"avatar" valid:"-"`
}

type BriefCompany struct {
	ID          uuid.UUID     `gorm:"column:comp_id;default:uuid_generate_v4()" json:"id"`
	Name        string        `gorm:"column:name;notnull" json:"name"`
}

func (v BriefCompany) TableName() string {
	return "main.official_companies"
}

//easyjson:json
type ListBriefCompany []BriefCompany
