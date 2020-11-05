package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
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

func (p *pgReopository) GetResumeArr(start, limit uint) ([]models.ResumeWithCandidate, error) {
	var brief []models.ResumeWithCandidate
	err := p.db.Model(&brief).
		Relation("CandidateWithUser").
		Relation("CandidateWithUser.User").
		Offset(int(start)).Limit(int(limit)).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume array from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return brief, nil
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


func (p *pgReopository) SearchResume(searchParams *models.SearchResume) ([]models.ResumeWithCandidate, error) {
	var brief []models.ResumeWithCandidate
	err := p.db.Model(&brief).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		if len(searchParams.AreaSearch) != 0 {
			q = q.Where("area_search IN (?)", pg.In(searchParams.AreaSearch))
		}
		if len(searchParams.Gender) != 0 {
			q = q.Where("gender IN (?)", pg.In(searchParams.Gender))
		}
		if len(searchParams.EducationLevel) != 0 {
			q = q.Where("education_level IN (?)", pg.In(searchParams.EducationLevel))
		}
		if len(searchParams.CareerLevel) != 0 {
			q = q.Where("career_level IN (?)", pg.In(searchParams.CareerLevel))
		}
		if len(searchParams.ExperienceMonth) != 0 {
			q = q.Where("experience_month IN (?)", pg.In(searchParams.ExperienceMonth))
		}
		q = q.Where("salary_min >= ?", searchParams.SalaryMin).
			Where("salary_max <= ?", searchParams.SalaryMax)
		return q, nil
	}).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		q = q.Where("LOWER(title) LIKE ?", "%" + searchParams.KeyWords+ "%").
			WhereOr("LOWER(place) LIKE ?", "%" + searchParams.KeyWords + "%")
		return q, nil
	}).Relation("CandidateWithUser").
		Relation("CandidateWithUser.User").
		Select()
	if err != nil {
		return nil, err
	}
	return brief, nil
}

func (p *pgReopository) AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error) {
	_, err := p.db.Model(&favoriteForEmpl).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting favorite resume: %w", err)
		return nil, err
	}
	return &favoriteForEmpl, nil
}

func (p *pgReopository) RemoveFavorite(favoriteForEmpl uuid.UUID) error {
	var favorite models.FavoritesForEmpl
	_, err := p.db.Model(&favorite).Where("favorite_id = ?", favoriteForEmpl).Delete()
	if err != nil {
		err = fmt.Errorf("error in delete favorite resume: %w", err)
		return err
	}
	return nil
}

func (p *pgReopository) GetAllEmplFavoriteResume(empl_id uuid.UUID) ([]models.FavoritesForEmplWithResume, error) {
	var brief []models.FavoritesForEmplWithResume
	err := p.db.Model(&brief).
		Relation("ResumeWithCandidate").
		Relation("ResumeWithCandidate.CandidateWithUser").
		Relation("ResumeWithCandidate.CandidateWithUser.User").
		Where("empl_id = ?", empl_id).Select()
	if err != nil {
		err = fmt.Errorf("error in get my favorite resume: %w", err)
		return nil, err
	}
	return brief, nil
}

func (p *pgReopository) GetFavoriteForResume(userID uuid.UUID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error) {
	var favorite models.FavoritesForEmpl
	err := p.db.Model(&favorite).
		Where("empl_id = ?", userID).
		Where("resume_id = ?", resumeID).
		Select()

	//TODO check on no rows in result
	if favorite.ID == uuid.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &favorite, nil
}
