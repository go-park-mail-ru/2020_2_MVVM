package repository

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pgRepository struct {
	db *gorm.DB
}

func NewPgRepository(db *gorm.DB) custom_experience.CustomExperienceRepository{
	return &pgRepository{db: db}
}

func (p *pgRepository) Create(experience models.ExperienceCustomComp) (*models.ExperienceCustomComp, error) {
	//_, err := p.db.Model(&experience).Returning("*").Insert()
	//if err != nil {
	//	err = fmt.Errorf("error in inserting custom experience with title: error: %w", err)
	//	return nil, err
	//}
	//return &experience, nil
	return nil, nil
}

//func (p *pgRepository) GetById(id string) (*models.ExperienceCustomComp, error) {
//	var experience models.ExperienceCustomComp
//	err := p.db.Model(&experience).Where("education_id = ?", id).Select()
//	if err != nil {
//		err = fmt.Errorf("error in select custom experience with id: %s : error: %w", id, err)
//		return nil, err
//	}
//	return &experience, nil
//}
//
//func (p *pgRepository) GetAllFromResume(resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
//	var experience []models.ExperienceCustomComp
//	err := p.db.Model(&experience).Where("resume_id = ?", resumeID).Limit(5).Select()
//	if err != nil {
//		return nil, err
//	}
//	return experience, nil
//}


func (p *pgRepository) DropAllFromResume(resumeID uuid.UUID) error {
	//var experience models.ExperienceCustomComp
	//_, err := p.db.Model(&experience).Where("resume_id = ?", resumeID).Delete()
	//return err
	return nil
}
