package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	resume2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/resume"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PGRepository struct {
	db *gorm.DB
}

func NewPgRepository(db *gorm.DB) resume.Repository {
	return &PGRepository{db: db}
}

func (p *PGRepository) Create(resume models.Resume) (*models.Resume, error) {
	candidate := new(models.Candidate)
	err := p.db.
		Joins("User").
		First(&candidate, resume.CandID).
		Error
	if err != nil {
		err = fmt.Errorf("error in FK search for resume creation for user with id: %s : error: %w", resume.CandID, err)
		return nil, err
	}
	err = p.db.Create(&resume).Error
	if err != nil {
		err = fmt.Errorf("error in inserting resume with title: %s : error: %w", resume.Title, err)
		return nil, err
	}
	return &resume, nil
}

func (p *PGRepository) GetById(id uuid.UUID) (*models.Resume, error) {
	var r models.Resume
	err := p.db.
		Preload(clause.Associations).
		First(&r, id).
		Error

	if err != nil {
		err = fmt.Errorf("error in select resume with id: %s : error: %w", id, err)
		return nil, err
	}
	return &r, nil
}

func (p *PGRepository) GetByIdWithCand(id uuid.UUID) (*models.Resume, error) {
	var r models.Resume
	err := p.db.
		Preload(clause.Associations).
		Preload("Candidate").
		Preload("Candidate.User").
		First(&r, id).
		Error

	if err != nil {
		err = fmt.Errorf("error in select resume with id: %s : error: %w", id, err)
		return nil, err
	}
	return &r, nil
}

func (p *PGRepository) List(start, limit uint) ([]models.Resume, error) {
	var brief []models.Resume
	err := p.db.
		Preload("Candidate").
		Preload("Candidate.User").
		Order("date_create desc").
		Offset(int(start)).Limit(int(limit)).
		Find(&brief).
		Error
	if err != nil {
		err = fmt.Errorf("error in select resume array from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return brief, nil
}

func (p *PGRepository) GetAllUserResume(userID uuid.UUID) ([]models.Resume, error) {
	var brief []models.Resume
	err := p.db.
		Preload("Candidate").
		Preload("Candidate.User").
		Order("date_create desc").
		Find(&brief, `Resume.cand_id = ?`, userID).
		Error
	if err != nil {
		err = fmt.Errorf("GetAllUserResume: %w", err)
		return nil, err
	}
	return brief, nil
}

func (p *PGRepository) Drop(resume models.Resume) error {
	return p.db.Delete(&resume).Error
}

func (p *PGRepository) Update(resume models.Resume) (*models.Resume, error) {
	candidate := new(models.Candidate)
	err := p.db.
		Joins("User").
		First(&candidate, resume.CandID).
		Error
	if err != nil {
		err = fmt.Errorf("error in FK search for resume creation for user with id: %s : error: %w", resume.CandID, err)
		return nil, err
	}

	for i := range resume.ExperienceCustomComp {
		resume.ExperienceCustomComp[i].CandID = candidate.ID
	}
	if err := p.db.Model(&resume).Where("resume_id = ?", resume.ResumeID).Updates(resume).Error; err != nil {
		return nil, fmt.Errorf("can't update resume with id:%s, %s", resume.ResumeID, err)
	}

	return &resume, nil
}

func (p *PGRepository) Search(searchParams *resume2.SearchParams) ([]models.Resume, error) {
	var brief []models.Resume
	err := p.db.Table("main.resume").Scopes(func(q *gorm.DB) *gorm.DB {
		if len(searchParams.AreaSearch) != 0 {
			q = q.Where("area_search IN (?)", searchParams.AreaSearch)
		}
		if len(searchParams.Gender) != 0 {
			q = q.Where("gender IN (?)", searchParams.Gender)
		}
		if len(searchParams.EducationLevel) != 0 {
			q = q.Where("education_level IN (?)", searchParams.EducationLevel)
		}
		if len(searchParams.CareerLevel) != 0 {
			q = q.Where("career_level IN (?)", searchParams.CareerLevel)
		}
		if len(searchParams.ExperienceMonth) != 0 {
			q = q.Where("experience_month IN (?)", searchParams.ExperienceMonth)
		}
		if len(searchParams.Sphere) != 0 {
			q = q.Where("sphere IN (?)", searchParams.Sphere)
		}
		if searchParams.SalaryMin != nil {
			q = q.Where("salary_min >= ?", searchParams.SalaryMin)
		}
		if searchParams.SalaryMax != nil {
			q = q.Where("salary_max <= ?", searchParams.SalaryMax)
		}
		if searchParams.KeyWords != nil {
			q = q.Where("LOWER(title) LIKE ?", "%"+*searchParams.KeyWords+"%")
		}
		if searchParams.KeywordsGeo != "" {
			q = q.Where("LOWER(area_search) LIKE ?", "%"+searchParams.KeywordsGeo+"%") //for main
		}
		if searchParams.StartLimit.Start != nil {
			q = q.Offset(int(*searchParams.StartLimit.Start))
		}
		if searchParams.StartLimit.Limit != nil {
			q = q.Limit(int(*searchParams.StartLimit.Limit))
		}
		return q
	}).Preload("Candidate").
		Preload("Candidate.User").
		Order("date_create desc").
		Find(&brief).
		Error
	if err != nil {
		err = fmt.Errorf("error in resumes list selection with searchParams: %s", err)
		return nil, err
	}
	//TODO:
	//}).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
	//	if searchParams.KeyWords != nil {
	//		q = q.Where("LOWER(title) LIKE ?", "%"+*searchParams.KeyWords+"%").
	//			WhereOr("LOWER(place) LIKE ?", "%"+*searchParams.KeyWords+"%")
	//	}
	return brief, nil
}

func (p *PGRepository) AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoriteID, error) {
	favoriteID := new(models.FavoriteID)
	err := p.db.Create(&favoriteForEmpl).Error
	favoriteID.FavoriteID = &favoriteForEmpl.FavoriteID
	if err != nil {
		err = fmt.Errorf("error in inserting favorite resume: %w", err)
		return nil, err
	}
	return favoriteID, nil
}

func (p *PGRepository) RemoveFavorite(favoriteForEmpl uuid.UUID) error {
	var favorite models.FavoritesForEmpl
	err := p.db.Where("favorite_id = ?", favoriteForEmpl.String()).Delete(&favorite).Error
	//_, err := p.db.Model(&favorite).Where("favorite_id = ?", favoriteForEmpl).Delete()
	if err != nil {
		err := fmt.Errorf("error in delete favorite resume: %s", err.Error())
		return err
	}
	return nil
}

func (p *PGRepository) GetAllEmplFavoriteResume(emplID uuid.UUID) ([]models.Resume, error) {
	var buffer = models.Employer{ID: emplID}

	err := p.db.
		Preload(clause.Associations).
		Preload("Favorites.Resume").
		Preload("Favorites.Resume.Candidate").
		Preload("Favorites.Resume.Candidate.User").
		First(&buffer).
		Order("date_create desc").
		Error

	var resumes []models.Resume
	for _, item := range buffer.Favorites {
		resumes = append(resumes, item.Resume)
	}

	return resumes, err
}

func (p *PGRepository) GetFavoriteForResume(userID uuid.UUID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error) {
	var favorite models.FavoritesForEmpl

	err := p.db.Where("empl_id = ?", userID).
		Where("resume_id = ?", resumeID).
		First(&favorite).Error
	//TODO check on no rows in result
	if favorite.FavoriteID == uuid.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &favorite, nil
}

func (p *PGRepository) GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error) {
	var favorite models.FavoritesForEmpl

	err := p.db.First(&favorite, favoriteID).Error

	//err := p.db.Model(&favorite).
	//	Where("favorite_id = ?", favoriteID).
	//	Select()
	if err != nil {
		return nil, err
	}

	return &favorite, nil
}

func (p *PGRepository) Delete(resId uuid.UUID, candId uuid.UUID) error {
	err := p.db.Table("main.resume").Delete(&models.Resume{ResumeID: resId}).Where("cand_id = ?", candId).Error
	if err != nil {
		return fmt.Errorf("error in delete resume with id: %s, err: %w", resId, err)
	}
	return nil
}

//func (p *PGRepository) SelectCandidate (candID uuid.UUID) (*models.Candidate, error) {
//	//var user models.Candidate
//	//err := p.db.Model(&user).
//	//	Relation("User").
//	//	Where("cand_id = ?", candID).
//	//	Select()
//	//if err != nil {
//	//	err = fmt.Errorf("error in inserting resume with title: error: %w", err)
//	//	return nil, err
//	//}
//	//return &user, err
//	err := fmt.Errorf("implement me")
//	return nil, err
//}
