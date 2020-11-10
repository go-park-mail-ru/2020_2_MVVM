package usecase

import (
	"errors"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	CompanyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	VacancyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/usecase"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"time"
)

type UseCaseResponse struct {
	infoLogger     *logger.Logger
	errorLogger    *logger.Logger
	resumeUsecase  resume.UseCase
	vacancyUsecase VacancyUseCase.VacancyUseCase
	companyUsecase CompanyUseCase.CompanyUseCase
	strg           response.ResponseRepository
}

func NewUsecase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	resumeUsecase resume.UseCase,
	vacancyUsecase VacancyUseCase.VacancyUseCase,
	companyUsecase CompanyUseCase.CompanyUseCase,
	strg response.ResponseRepository,
) *UseCaseResponse {
	usecase := UseCaseResponse{
		infoLogger:     infoLogger,
		errorLogger:    errorLogger,
		resumeUsecase:  resumeUsecase,
		vacancyUsecase: vacancyUsecase,
		companyUsecase: companyUsecase,
		strg:           strg,
	}
	return &usecase
}

func (u *UseCaseResponse) Create(response models.Response) (*models.Response, error) {
	if response.Initial == common.Candidate {
		r, err := u.resumeUsecase.GetById(response.ResumeID)
		if r == nil || err != nil {
			err = errors.New("this user cannot send response from this resume")
			return nil, err
		}
	} else if response.Initial == common.Employer {
		r, err := u.vacancyUsecase.GetVacancy(response.VacancyID)
		if r == nil || err != nil {
			err = errors.New("this user cannot send response from this vacancy")
			return nil, err
		}
	}
	response.DateCreate = time.Now()
	response.Status = "sent"
	return u.strg.Create(response)
}

func (u *UseCaseResponse) UpdateStatus(response models.Response, userUpdate string) (*models.Response, error) {
	if response.Status == "sent" {
		return nil, nil
	}
	oldResp, err := u.strg.GetByID(response.ID)

	if oldResp.Initial == userUpdate {
		err = fmt.Errorf("this user cannot update status in response")
		return nil, err
	}

	if userUpdate == common.Candidate {
		r, err := u.resumeUsecase.GetById(oldResp.ResumeID)
		if r == nil || err != nil {
			err = errors.New("this user cannot update response from this resume")
			return nil, err
		}
	} else if userUpdate == common.Employer {
		r, err := u.vacancyUsecase.GetVacancy(oldResp.VacancyID)
		if r == nil || err != nil {
			err = errors.New("this user cannot update response from this vacancy")
			return nil, err
		}
	}

	return u.strg.UpdateStatus(response)
}

func (u *UseCaseResponse) GetAllCandidateResponses(candID uuid.UUID) ([]models.ResponseWithTitle, error) {
	resumes, err := u.resumeUsecase.GetAllUserResume(candID)
	if err != nil {
		return nil, err
	}
	var responses []models.ResponseWithTitle
	for i := range resumes {
		resp, err := u.strg.GetResumeAllResponse(resumes[i].ResumeID)
		if err != nil {
			return nil, err
		}
		for j := range resp {
			responseWithTitle := models.ResponseWithTitle{
				ResponseID: resp[j].ID,
				ResumeID:   resp[j].ResumeID,
				ResumeName: resumes[i].Title,
				VacancyID:  resp[j].VacancyID,
				Initial:    resp[j].Initial,
				Status:     resp[j].Status,
				DateCreate: resp[j].DateCreate,
			}
			responses = append(responses, responseWithTitle)
		}
	}
	for i := range responses {
		vac, err := u.vacancyUsecase.GetVacancy(responses[i].VacancyID)
		if err != nil {
			return nil, err
		}
		comp, err := u.companyUsecase.GetOfficialCompany(vac.CompID)
		if err != nil {
			return nil, err
		}
		responses[i].VacancyName = vac.Title
		responses[i].CompanyID = vac.CompID
		responses[i].CompanyName = comp.Name
	}
	return responses, nil
}

func (u *UseCaseResponse) GetAllEmployerResponses(emplID uuid.UUID) ([]models.ResponseWithTitle, error) {
	vacancyList, err := u.vacancyUsecase.GetVacancyList(0, 100, emplID, vacancy.ByEmpId)
	if err != nil {
		return nil, err
	}
	var responses []models.ResponseWithTitle
	for i := range vacancyList {
		comp, err := u.companyUsecase.GetOfficialCompany(vacancyList[i].CompID)
		resp, err := u.strg.GetVacancyAllResponse(vacancyList[i].ID)
		if err != nil {
			return nil, err
		}
		for j := range resp {
			responseWithTitle := models.ResponseWithTitle{
				ResponseID:  resp[j].ID,
				ResumeID:    resp[j].ResumeID,
				VacancyName: vacancyList[i].Title,
				VacancyID:   resp[j].VacancyID,
				CompanyID:   vacancyList[i].CompID,
				CompanyName: comp.Name,
				Initial:     resp[j].Initial,
				Status:      resp[j].Status,
				DateCreate:  resp[j].DateCreate,
			}
			responses = append(responses, responseWithTitle)
		}
	}
	for i := range responses {
		vac, err := u.resumeUsecase.GetById(responses[i].ResumeID)
		if err != nil {
			return nil, err
		}
		responses[i].VacancyName = vac.Title
	}
	return responses, nil
}

func (u *UseCaseResponse) GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.BriefResumeInfo, error) {
	r, err := u.strg.GetAllResumeWithoutResponse(candID, vacancyID)
	if err != nil {
		return nil, err
	}

	var briefRespResumes []models.BriefResumeInfo
	for i := range r {
		var insert models.BriefResumeInfo
		err = copier.Copy(&insert, &r[i])
		if err != nil {
			err = fmt.Errorf("error in copy resume for GetAllResumeWithoutResponse: %w", err)
			return nil, err
		}
		briefRespResumes = append(briefRespResumes, insert)
	}
	return briefRespResumes, nil
}
