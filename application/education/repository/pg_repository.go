package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
)

type pgReopository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) education.Repository {
	return &pgReopository{db: db}
}

func (p *pgReopository) Create(education []*models.Education) ([]models.Education, error) {
	var dbEducation []models.Education
	for _,item := range education {
		if item == nil {
			continue
		}

		_, err := p.db.Model(item).Returning("*").Insert()
		if err != nil {
			err = fmt.Errorf("error in inserting resume with title: error: %w", err)
			return nil, err
		}
		dbEducation = append(dbEducation, *item)
	}
	return dbEducation, nil
}

func (p *pgReopository) GetById(id string) (*models.Education, error) {
	var educatuin models.Education
	err := p.db.Model(&educatuin).Where("education_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select education with id: %s : error: %w", id, err)
		return nil, err
	}
	return &educatuin, nil
}

func (p *pgReopository) GetAllFromResume(resumeID uuid.UUID) ([]models.Education, error) {
	var educations []models.Education
	err := p.db.Model(&educations).Where("resume_id = ?", resumeID).Limit(5).Select()
	if err != nil {
		return nil, err
	}
	return educations, nil
}

func (p *pgReopository) DeleteAllResumeEducation(resumeID uuid.UUID) error {
	var educations models.Education
	_, err := p.db.Model(&educations).Where("resume_id = ?", resumeID).Delete()
	if err != nil {
		return err
	}
	return nil
}