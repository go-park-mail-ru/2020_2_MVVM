package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	logger "github.com/rowdyroad/go-simple-logger"
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

func (u *UseCaseResume) CreateResume(resume models.Resume) (*models.Resume, error) {
	r, err := u.strg.CreateResume(resume)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *UseCaseResume) UpdateResume(resume models.Resume) (*models.Resume, error) {
	//if resume.ID == uuid.Nil {
	//	err := fmt.Errorf("error in update resume: resume does not exist")
	//	return nil, err
	//}

	// ID from Session
	r, err := u.strg.UpdateResume(resume.ID, &resume)
	if err != nil {
		err = fmt.Errorf("error in update resume: %w", err)
		return nil, err
	}
	return r, nil
}

func (u *UseCaseResume) GetResume(id string) (*models.Resume, error) {
	r, err := u.strg.GetResumeById(id)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *UseCaseResume) GetResumeList(begin, end uint) ([]models.Resume, error) {
	r, err := u.strg.GetResumeArr(begin, end)
	if err != nil {
		err = fmt.Errorf("USE error in resume get list from %v to %v: error: %w", begin, end, err)
		return nil, err
	}
	return r, nil
}
