package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-pg/pg/v9"
)

type pgReopository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) resume.ResumeRepository {
	return &pgReopository{db: db}
}

func (p *pgReopository) CreateResume(resume models.Resume) (models.Resume, error) {
	_, err := p.db.Model(&resume).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting resume with title: %s : error: %w", resume.Title, err)
		return models.Resume{}, err
	}
	return resume, nil
}

//func (p *pgReopository) UpdateResume(resume models.Resume) (models.Resume, error) {
//	newResume := new(models.Resume)
//	err := p.db.Model(&resume).Where("resume_id = ?", resume.ID).First()
//	if err != nil {
//		err = fmt.Errorf("error in create resume with id: %s : error: %w", resume.ID, err)
//		return models.Resume{}, err
//	}
//	if p.Name != "" {
//		savedPerson.Name = p.Name
//	}
//	if p.Occupation != "" {
//		savedPerson.Occupation = p.Occupation
//	}
//	if p.BirthDate != "" {
//		savedPerson.BirthDate = p.BirthDate
//	}
//	if p.BirthPlace != "" {
//		savedPerson.BirthPlace = p.BirthPlace
//	}
//	if p.Image != "" {
//		savedPerson.Image = p.Image
//	}
//}

func (p *pgReopository) GetResumeById(id string) (models.Resume, error) {
	var r models.Resume
	err := p.db.Model(&r).Where("resume_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume with id: %s : error: %w", id, err)
		return models.Resume{}, err
	}
	return r, nil
}

func (p *pgReopository) GetResumeByName(name string) (models.Resume, error) {
	var r models.Resume
	err := p.db.Model(&r).Where("title = ?", name).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume with title: %s : error: %w", name, err)
		return models.Resume{}, err
	}
	return r, nil
}





