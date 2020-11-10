package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	CompanyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	VacancyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/usecase"
	"github.com/google/uuid"
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
	response.DateCreate = time.Now()
	response.Status = "sent"
	return u.strg.Create(response)
}

func (u *UseCaseResponse) UpdateStatus(response models.Response) (*models.Response, error) {
	if response.Status == "sent" {
		return nil, nil
	}
	return u.strg.UpdateStatus(response)
}

func (u *UseCaseResponse) GetAllUserResponses(candID uuid.UUID) ([]models.ResponseWithTitle, error) {
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
