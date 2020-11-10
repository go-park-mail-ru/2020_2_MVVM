package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"

	//"github.com/google/uuid"
)

type pgReopository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) response.ResponseRepository {
	return &pgReopository{db: db}
}

func (p *pgReopository) Create(response models.Response) (*models.Response, error) {
	_, err := p.db.Model(&response).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting response: %w", err)
		return nil, err
	}
	return &response, nil
}

func (p *pgReopository) GetByID(responseID uuid.UUID) (*models.Response, error) {
	var response models.Response
	err := p.db.Model(&response).Where("response_id = ?", responseID).Select()
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (p *pgReopository) UpdateStatus(response models.Response) (*models.Response, error) {
	_, err := p.db.Model(&response).WherePK().Returning("*").UpdateNotZero()
	if err != nil {
		err = fmt.Errorf("error in updating response with id %s, : %w", response.ID.String(), err)
		return nil, err
	}
	return &response, nil
}

func (p *pgReopository)GetResumeAllResponse(resumeID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := p.db.Model(&responses).Where("resume_id = ?", resumeID).Select()
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgReopository)GetVacancyAllResponse(vacancyID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := p.db.Model(&responses).Where("vacancy_id = ?", vacancyID).Select()
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}
