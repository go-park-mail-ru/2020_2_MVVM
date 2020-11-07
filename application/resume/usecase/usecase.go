package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	UserUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/usecase"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"strings"
	"time"
)

type ResumeUseCase struct {
	infoLogger       *logger.Logger
	errorLogger      *logger.Logger
	userUseCase      UserUseCase.UserUseCase
	educationUseCase education.UseCase
	customExpUseCase custom_experience.UseCase
	strg             resume.Repository
}

func NewUseCase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	userUseCase UserUseCase.UserUseCase,
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
	template.DateCreate = time.Now()
	result, err := u.strg.Create(template)

	for i := range template.ExperienceCustomComp {
		template.ExperienceCustomComp[i].ResumeID = result.ResumeID
		template.ExperienceCustomComp[i].CandID = result.CandID
	}
	for i := range template.Education {
		template.Education[i].ResumeId = result.ResumeID
		template.Education[i].CandID = result.CandID
	}

	result.Education, _ = u.educationUseCase.Create(template.Education)
	result.ExperienceCustomComp, _ = u.customExpUseCase.Create(template.ExperienceCustomComp)


	//user := *resume.Candidate.User
	//resume.Candidate = nil
	//
	//var additionParam models.AdditionInResume
	//if err := ctx.ShouldBindBodyWith(&additionParam, binding.JSON); err != nil {
	//	ctx.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	//
	//pEducations, err := r.createEducation(additionParam.Education, candID,  resume.ResumeID)
	//if err != nil {
	//	ctx.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	//
	//var customExperience []models.ExperienceCustomComp
	//for i := range additionParam.CustomExperience {
	//	item := additionParam.CustomExperience[i]
	//	dateBegin, err := time.Parse(time.RFC3339, item.Begin+"T00:00:00Z")
	//	if err != nil {
	//		ctx.AbortWithError(http.StatusBadRequest, err)
	//		return
	//	}
	//	var dateFinish time.Time
	//	if !item.ContinueToToday {
	//		dateFinish, err = time.Parse(time.RFC3339, *item.Finish+"T00:00:00Z")
	//		if err != nil {
	//			ctx.AbortWithError(http.StatusBadRequest, err)
	//			return
	//		}
	//	} else {
	//		dateFinish = time.Now()
	//	}
	//	//dateBegin := time.Now()
	//	//dateFinish := time.Now()
	//
	//	insertExp := models.ExperienceCustomComp{
	//		NameJob:         item.NameJob,
	//		Position:        item.Position,
	//		Begin:           dateBegin,
	//		Finish:          &dateFinish,
	//		Duties:          item.Duties,
	//		ContinueToToday: &item.ContinueToToday,
	//	}
	//	customExperience = append(customExperience, insertExp)
	//}
	//
	//pCustomExperience, err := r.createCustomExperience(customExperience, candID,  resume.ResumeID)
	//if err != nil {
	//	ctx.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}

	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
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
	r, err := u.strg.Update(&resume)
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
