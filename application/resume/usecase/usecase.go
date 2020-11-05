package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"math"
	"strings"
)

type UseCaseResume struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        resume.ResumeRepository
}

func NewUsecase(infoLogger *logger.Logger,
				errorLogger *logger.Logger,
				strg resume.ResumeRepository) *UseCaseResume {
					usecase := UseCaseResume {
					infoLogger:  infoLogger,
					errorLogger: errorLogger,
					strg:        strg,
	}
	return &usecase
}

func (u* UseCaseResume) GetAllUserResume(userid uuid.UUID) ([]models.Resume, error) {
	return u.strg.GetAllUserResume(userid)
}

func (u *UseCaseResume) CreateResume(resume models.Resume) (*models.Resume, error) {
	r, err := u.strg.CreateResume(resume)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *UseCaseResume) UpdateResume(resume models.Resume) (*models.Resume, error) {
	oldResume, err := u.strg.GetResumeById(resume.ID.String())
	if err != nil {
		err = fmt.Errorf("error in get resume by id: %w", err)
		return nil, err
	}
	if resume.UserID != oldResume.UserID {
		err = fmt.Errorf("this user cannot update this resume")
		return nil, err
	}
	r, err := u.strg.UpdateResume(&resume)
	if err != nil {
		err = fmt.Errorf("error in update resume: %w", err)
		return nil, err
	}
	return r, nil
}

func (u *UseCaseResume)SearchResume(searchParams models.SearchResume) ([]models.Resume, error)  {
	if searchParams.SalaryMax == 0 {
		searchParams.SalaryMax = math.MaxInt64
	}
	searchParams.KeyWords = strings.ToLower(searchParams.KeyWords)
	return u.strg.SearchResume(&searchParams)
}

func (u *UseCaseResume) GetResume(id string) (*models.Resume, error) {
	r, err := u.strg.GetResumeById(id)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *UseCaseResume) GetResumePage(start, limit uint) ([]models.BriefRespResume, error) {
	if limit >= 200 {
		return nil, fmt.Errorf("Limit is too high. ")
	}
	r, err := u.strg.GetResumeArr(start, limit)
	if err != nil {
		err = fmt.Errorf("error in resume get list from %v to %v: error: %w", start, limit, err)
		return nil, err
	}

	var briefRespResumes []models.BriefRespResume
	for i := range r {
		var insert models.BriefRespResume
		err = copier.Copy(&insert, &r[i])
		if err != nil {
			err = fmt.Errorf("error in resume get list from %v to %v: error: %w", start, limit, err)
			return nil, err
		}
		insert.UserID = r[i].CandidateWithUser.UserID
		insert.Name = r[i].CandidateWithUser.User.Name
		insert.Surname = r[i].CandidateWithUser.User.Surname
		insert.Email = r[i].CandidateWithUser.User.Email
		briefRespResumes = append(briefRespResumes, insert)
	}


	return briefRespResumes, nil
}

func (u *UseCaseResume) AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error) {
	return u.strg.AddFavorite(favoriteForEmpl)
}

func (u *UseCaseResume) RemoveFavorite(favoriteForEmpl uuid.UUID) error {
	return u.strg.RemoveFavorite(favoriteForEmpl)
}

func (u *UseCaseResume) GetAllEmplFavoriteResume(userID uuid.UUID) ([]models.Resume, error) {
	return u.strg.GetAllEmplFavoriteResume(userID)
}

func (u *UseCaseResume) GetFavoriteForResume(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error) {
	return u.strg.GetFavoriteForResume(userID, resumeID )
}
