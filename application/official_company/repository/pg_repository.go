package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
)

type pgReopository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) official_company.OfficialCompanyRepository{
	return &pgReopository{db: db}
}

func (p *pgReopository) CreateOfficialCompany(company models.OfficialCompany) (*models.OfficialCompany, error) {
	_, err := p.db.Model(&company).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting official company: error: %w", err)
		return nil, err
	}
	return &company, nil
}

func (p *pgReopository) GetOfficialCompanyById(id string) (*models.OfficialCompany, error) {
	var company models.OfficialCompany
	err := p.db.Model(&company).Where("company_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select official company with id: %s : error: %w", id, err)
		return nil, err
	}
	return &company, nil
}
