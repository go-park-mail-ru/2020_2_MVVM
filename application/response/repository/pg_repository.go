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

func (p *pgReopository) GetResumeAllResponse(resumeID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := p.db.Model(&responses).Where("resume_id = ?", resumeID).Select()
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgReopository) GetVacancyAllResponse(vacancyID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := p.db.Model(&responses).Where("vacancy_id = ?", vacancyID).Select()
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgReopository) GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.Resume, error) {
	var resume []models.Resume
	query := fmt.Sprintf(`select main.resume.*
			from main.resume
			left join main.response on main.response.resume_id = main.resume.resume_id
			where cand_id = '%s'
			group by main.resume.resume_id
			having sum(case when vacancy_id = '%s' then 1 else 0 end) = 0`, candID, vacancyID)
	_, err := p.db.Query(&resume, query)
	if err != nil {
		return nil, err
	}
	return resume, nil
}

func (p *pgReopository) GetAllVacancyWithoutResponse(emplID uuid.UUID, resumeID uuid.UUID) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	query := fmt.Sprintf(`select main.vacancy.*
			from main.vacancy
			left join main.response on main.response.vacancy_id = main.vacancy.vac_id
			where empl_id = '%s'
			group by main.vacancy.vac_id
			having sum(case when resume_id = '%s' then 1 else 0 end) = 0`, emplID, resumeID)
	_, err := p.db.Query(&vacancies, query)
	if err != nil {
		return nil, err
	}
	return vacancies, nil
}