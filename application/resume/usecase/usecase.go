package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"strings"
	"time"
)

type ResumeUseCase struct {
	infoLogger       *logger.Logger
	errorLogger      *logger.Logger
	userUseCase      user.UseCase
	educationUseCase education.UseCase
	customExpUseCase custom_experience.UseCase
	strg             resume.Repository
}

func NewUseCase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	userUseCase user.UseCase,
	educationUseCase education.UseCase,
	customExpUseCase custom_experience.UseCase,
	strg resume.Repository) resume.UseCase {
	usecase := ResumeUseCase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		userUseCase: userUseCase,
		educationUseCase: educationUseCase,
		customExpUseCase: customExpUseCase,
		strg:        strg,
	}
	return &usecase
}

func (u *ResumeUseCase) Create(template models.Resume) (*models.Resume, error) {
	// create resume
	template.DateCreate = time.Now()
	result, err := u.strg.Create(template)
	if err != nil {
		return nil, err
	}

	err = u.createExperienceAndEducation(template, *result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *ResumeUseCase) createExperienceAndEducation (template models.Resume, result models.Resume) error {
	// create experience
	var err error
	for i := range template.ExperienceCustomComp {
		if template.ExperienceCustomComp[i] == nil {
			continue
		}
		if *(template.ExperienceCustomComp[i].ContinueToToday){
			dateFinish := time.Now()
			template.ExperienceCustomComp[i].Finish = &dateFinish
		}
		template.ExperienceCustomComp[i].ResumeID = result.ResumeID
		template.ExperienceCustomComp[i].CandID = result.CandID

		template.ExperienceCustomComp[i], err = u.customExpUseCase.Create(*template.ExperienceCustomComp[i])
		if err != nil {
			return err
		}

	}
	// create education
	for i := range template.Education {
		if template.Education[i] == nil {
			continue
		}
		template.Education[i].ResumeId = result.ResumeID
		template.Education[i].CandID = result.CandID
		template.Education[i], err = u.educationUseCase.Create(*template.Education[i])
		if err != nil {
			return err
		}
	}
	return nil
}


func (u *ResumeUseCase) Update(resume models.Resume) (*models.Resume, error) {
	oldResume, err := u.strg.GetById(resume.ResumeID)
	if err != nil {
		err = fmt.Errorf("error in get resume by id: %w", err)
		return nil, err
	}
	if resume.CandID != oldResume.CandID {
		err = fmt.Errorf("this user cannot update this resume")
		return nil, err
	}
	resume.DateCreate = oldResume.DateCreate
	err = u.customExpUseCase.DropAllFromResume(resume.ResumeID)
	if err != nil {
		return nil, err
	}
	err = u.educationUseCase.DropAllFromResume(resume.ResumeID)
	if err != nil {
		return nil, err
	}

	result, err := u.strg.Update(resume)
	err = u.createExperienceAndEducation(resume, resume)
	return result, err
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

func (u *ResumeUseCase) Search(searchParams resume.SearchParams) ([]models.BriefResumeInfo, error) {
	if searchParams.KeyWords != nil {
		*searchParams.KeyWords = strings.ToLower(*searchParams.KeyWords)
	}
	r, err := u.strg.Search(&searchParams)
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

func (u *ResumeUseCase) GetById(id uuid.UUID) (*models.Resume, error) {
	r, err := u.strg.GetById(id)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *ResumeUseCase) List(start, limit uint) ([]models.BriefResumeInfo, error) {
	if limit >= 200 {
		return nil, fmt.Errorf("Limit is too high. ")
	}
	r, err := u.strg.List(start, limit)
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

func (u *ResumeUseCase) RemoveFavorite(favoriteForEmpl models.FavoritesForEmpl) error {
	oldFavorite, err := u.strg.GetFavoriteByID(favoriteForEmpl.FavoriteID)
	if err != nil {
		return err
	}
	if oldFavorite.EmplID != favoriteForEmpl.EmplID {
		err = fmt.Errorf("error in remove favorite: method not allowed")
		return err
	}
	return u.strg.RemoveFavorite(favoriteForEmpl.FavoriteID)
}

func (u *ResumeUseCase) GetAllEmplFavoriteResume(userID uuid.UUID) ([]models.BriefResumeInfo, error) {
	r, err := u.strg.GetAllEmplFavoriteResume(userID)
	if err != nil {
		err = fmt.Errorf("error in get list favorite resume: %w", err)
		return nil, err
	}
	return DoBriefRespResume(r)
}

func (u *ResumeUseCase) GetFavoriteByResume(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error) {
	return u.strg.GetFavoriteForResume(userID, resumeID)
}

func (u *ResumeUseCase) GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error) {
	return u.strg.GetFavoriteByID(favoriteID)
}

func DoBriefRespResume(resumes []models.Resume) ([]models.BriefResumeInfo, error) {
	var briefRespResumes []models.BriefResumeInfo
	for i := range resumes {
		brief, err := resumes[i].Brief()
		if err != nil {
			return nil, err
		}
		briefRespResumes = append(briefRespResumes, *brief)
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
