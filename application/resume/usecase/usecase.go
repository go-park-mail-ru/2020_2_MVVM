package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"strings"
)

type ResumeUseCase struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        resume.Repository
}

func NewUsecase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	strg resume.Repository) resume.UseCase {
	usecase := ResumeUseCase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		strg:        strg,
	}
	return &usecase
}

func (u *ResumeUseCase) CreateResume(resume models.Resume) (*models.Resume, error) {
	r, err := u.strg.CreateResume(resume)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *ResumeUseCase) UpdateResume(resume models.Resume) (*models.Resume, error) {
	oldResume, err := u.strg.GetResumeById(resume.ResumeID.String())
	if err != nil {
		err = fmt.Errorf("error in get resume by id: %w", err)
		return nil, err
	}
	if resume.CandID != oldResume.CandID {
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

func (u *ResumeUseCase) GetAllUserResume(userid uuid.UUID) ([]models.BriefResumeInfo, error) {
	r, err := u.strg.GetAllUserResume(userid)
	if err != nil {
		err = fmt.Errorf("error in get my resume: %w", err)
		return nil, err
	}

	var briefRespResumes []models.BriefResumeInfo
	for i := range r {
		var insert models.BriefResumeInfo
		err = copier.Copy(&insert, &r[i])
		if err != nil {
			err = fmt.Errorf("error in copy resume for get my resume: %w", err)
			return nil, err
		}
		insert = DoBriefRespUser(insert, *r[i].Candidate)
		briefRespResumes = append(briefRespResumes, insert)
	}
	return briefRespResumes, nil
}

func (u *ResumeUseCase) SearchResume(searchParams resume.SearchParams) ([]models.BriefResumeInfo, error) {
	if searchParams.KeyWords != nil {
		*searchParams.KeyWords = strings.ToLower(*searchParams.KeyWords)
	}

	r, err := u.strg.SearchResume(&searchParams)
	if err != nil {
		err = fmt.Errorf("error in resume search: %w", err)
		return nil, err
	}

	var briefRespResumes []models.BriefResumeInfo
	for i := range r {
		var insert models.BriefResumeInfo
		err = copier.Copy(&insert, &r[i])
		if err != nil {
			err = fmt.Errorf("error in copy resume for search: %w", err)
			return nil, err
		}
		insert = DoBriefRespUser(insert, *r[i].Candidate)
		briefRespResumes = append(briefRespResumes, insert)
	}
	return briefRespResumes, nil
}

func (u *ResumeUseCase) GetResume(id string) (*models.Resume, error) {
	r, err := u.strg.GetResumeById(id)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *ResumeUseCase) GetResumePage(start, limit uint) ([]models.BriefResumeInfo, error) {
	if limit >= 200 {
		return nil, fmt.Errorf("Limit is too high. ")
	}
	r, err := u.strg.GetResumeArr(start, limit)
	if err != nil {
		err = fmt.Errorf("error in resume get list from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	briefRespResumes, err := DoBriefRespResume(r)
	if err != nil {
		err = fmt.Errorf("error in resume get list from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return briefRespResumes, nil
}

func (u *ResumeUseCase) AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error) {
	return u.strg.AddFavorite(favoriteForEmpl)
}

func (u *ResumeUseCase) RemoveFavorite(favoriteForEmpl uuid.UUID) error {
	return u.strg.RemoveFavorite(favoriteForEmpl)
}

func (u *ResumeUseCase) GetAllEmplFavoriteResume(userID uuid.UUID) ([]models.BriefResumeInfo, error) {
	r, err := u.strg.GetAllEmplFavoriteResume(userID)
	if err != nil {
		err = fmt.Errorf("error in get list favorite resume: %w", err)
		return nil, err
	}

	var briefRespResumes []models.BriefResumeInfo
	for i := range r {
		var insert models.BriefResumeInfo
		err = copier.Copy(&insert, &r[i].ResumeWithCandidate)
		if err != nil {
			err = fmt.Errorf("error in copy resumes for list my favorite: %w", err)
			return nil, err
		}
		insert = DoBriefRespUser(insert, *r[i].ResumeWithCandidate.Candidate)
		briefRespResumes = append(briefRespResumes, insert)
	}
	return briefRespResumes, nil

}

func (u *ResumeUseCase) GetFavorite(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error) {
	return u.strg.GetFavoriteForResume(userID, resumeID)
}

func DoBriefRespResume(resumes []models.Resume) ([]models.BriefResumeInfo, error) {
	var briefRespResumes []models.BriefResumeInfo
	for i := range resumes {
		var insert models.BriefResumeInfo
		err := copier.Copy(&insert, &resumes[i])
		if err != nil {
			return nil, err
		}
		insert = DoBriefRespUser(insert, *resumes[i].Candidate)
		briefRespResumes = append(briefRespResumes, insert)
	}
	return briefRespResumes, nil
}

func DoBriefRespUser(respResume models.BriefResumeInfo, user models.Candidate) models.BriefResumeInfo {
	respResume.UserID = user.UserID
	respResume.Name = user.User.Name
	respResume.Surname = user.User.Surname
	respResume.Email = user.User.Email
	return respResume
}
