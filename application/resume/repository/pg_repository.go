package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"

	//"github.com/google/uuid"
)

type pgReopository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) resume.ResumeRepository {
	return &pgReopository{db: db}
}

func (p *pgReopository) CreateResume(resume models.Resume) (*models.Resume, error) {
	_, err := p.db.Model(&resume).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting resume with title: error: %w", err)
		return nil, err
	}
	return &resume, nil
}

func (p *pgReopository) GetResumeById(id string) (*models.Resume, error) {
	var r models.Resume
	err := p.db.Model(&r).Where("resume_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume with id: %s : error: %w", id, err)
		return nil, err
	}
	return &r, nil
}

func (p *pgReopository) GetResumeByName(name string) (*models.Resume, error) {
	var r models.Resume
	err := p.db.Model(&r).Where("title = ?", name).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume with title: %s : error: %w", name, err)
		return nil, err
	}
	return &r, nil
}

func (p *pgReopository) GetResumeArr(start, limit uint) ([]models.Resume, error) {
	var r []models.Resume
	err := p.db.Model(&r).Offset(int(start)).Limit(int(limit)).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume array from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return r, nil
}

func (p *pgReopository) GetAllUserResume(userID uuid.UUID) ([]models.Resume, error) {
	var r[]models.Resume

	err := p.db.Model(&r).Where("cand_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (p *pgReopository) UpdateResume(newResume *models.Resume) (*models.Resume, error) {
	_, err := p.db.Model(newResume).WherePK().Returning("*").Update()
	if err != nil {
		err = fmt.Errorf("error in updating resume with id %s, : %w", newResume.ID.String(), err)
		return nil, err
	}
	return newResume, nil
}
