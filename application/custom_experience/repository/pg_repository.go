package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
)

type pgReopository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) custom_experience.CustomExperienceRepository{
	return &pgReopository{db: db}
}

func (p *pgReopository) CreateCustomExperience(experience []models.ExperienceCustomComp) ([]models.ExperienceCustomComp, error) {
	var dbExperience []models.ExperienceCustomComp
	for _,item := range experience {
		_, err := p.db.Model(&item).Returning("*").Insert()
		if err != nil {
			err = fmt.Errorf("error in inserting custom experience with title: error: %w", err)
			return nil, err
		}
		dbExperience = append(dbExperience, item)
	}
	return dbExperience, nil
}

func (p *pgReopository) GetCustomExperienceById(id string) (*models.ExperienceCustomComp, error) {
	var experience models.ExperienceCustomComp
	err := p.db.Model(&experience).Where("education_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select custom experience with id: %s : error: %w", id, err)
		return nil, err
	}
	return &experience, nil
}

func (p *pgReopository) GetAllResumeCustomExperience(resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	var experience []models.ExperienceCustomComp
	err := p.db.Model(&experience).Where("resume_id = ?", resumeID).Limit(5).Select()
	if err != nil {
		return nil, err
	}
	return experience, nil
}


func (p *pgReopository) DeleteAllResumeCustomExperience(resumeID uuid.UUID) error {
	var experience models.ExperienceCustomComp
	_, err := p.db.Model(&experience).Where("resume_id = ?", resumeID).Delete()
	if err != nil {
		return err
	}
	return nil
}
