package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
)

type pgReopository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) custom_company.CustomCompanyRepository{
	return &pgReopository{db: db}
}

func (p *pgReopository) CreateCustomCompany(company models.CustomCompany) (*models.CustomCompany, error) {
	_, err := p.db.Model(&company).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting custom company: error: %w", err)
		return nil, err
	}
	return &company, nil
}

func (p *pgReopository) GetCustomCompanyById(id string) (*models.CustomCompany, error) {
	var company models.CustomCompany
	err := p.db.Model(&company).Where("custom_company_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select custom company with id: %s : error: %w", id, err)
		return nil, err
	}
	return &company, nil
}
