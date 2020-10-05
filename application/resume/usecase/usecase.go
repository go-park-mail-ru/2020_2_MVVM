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

func (u *UseCaseResume) CreateResume(resume models.Resume) (models.Resume, error) {
	r, err := u.strg.CreateResume(resume)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return models.Resume{}, err
	}
	return r, nil
}
//UpdateResume(id uint) (models.Resume, bool)
func (u *UseCaseResume) GetResume(id string) (models.Resume, error) {
	r, err := u.strg.GetResumeById(id)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return models.Resume{}, err
	}
	return r, nil
}

//func (u *UseCaseResume) GetResumeList(begin, end uint) (models.Resume, error) {
//	r, ok := u.strg.GetResume(id)
//	if !ok {
//		return models.Films{}, false
//	}
//	return *films, true
//}


