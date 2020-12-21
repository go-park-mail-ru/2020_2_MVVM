package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	mCompany "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/official_company"
	mResponse "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/response"
	mResume "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/resume"
	mUser "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/user"
	mVacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/vacancy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

var ID, _ = uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
var testResponse = models.Response{
	ID:        ID,
	ResumeID:  ID,
	VacancyID: ID,
	Status:    "sent",
}
var testUser = models.User{
	ID:       ID,
	UserType: common.Candidate,
	Name:     "ID",
	Surname:  "ID",
	Email:    "ID",
}
var candidate = models.Candidate{
	ID:     ID,
	UserID: ID,
	User:   testUser,
}
var testResume = models.Resume{
	ResumeID:  ID,
	CandID:    ID,
	Candidate: candidate,
	Title:     "ID",
}
var briefResume = models.BriefResumeInfo{
	ResumeID: ID,
	CandID:   ID,
	UserID:   ID,
	Title:    "ID",
	Name:     "ID",
	Surname:  "ID",
	Email:    "ID",
}
var testVacacy = models.Vacancy{
	ID:     ID,
	EmpID:  ID,
	CompID: ID,
}
var testCompany = models.OfficialCompany{
	ID: ID,
}
var respWithTitle = models.ResponseWithTitle{
	ResponseID:  ID,
	ResumeID:    ID,
	ResumeName:  "ID",
	CandName:    "ID",
	CandSurname: "ID",
	VacancyID:   ID,
	VacancyName: "",
	CompanyID:   ID,
	CompanyName: "",
	Initial:     "",
	Status:      "sent",
}

func beforeTest(t *testing.T) (*mResponse.ResponseRepository, *mResume.UseCase, *mVacancy.IUseCaseVacancy,
	*mCompany.IUseCaseOfficialCompany, *mUser.UseCase, UseCaseResponse) {
	infoLogger, _ := logger.New(os.Stdout)
	errorLogger, _ := logger.New(os.Stderr)
	mockRepo := new(mResponse.ResponseRepository)
	mockResumeUS := new(mResume.UseCase)
	mockVacancyUS := new(mVacancy.IUseCaseVacancy)
	mockCompanyUS := new(mCompany.IUseCaseOfficialCompany)
	mockUserUS := new(mUser.UseCase)
	usecase := UseCaseResponse{
		infoLogger:     infoLogger,
		errorLogger:    errorLogger,
		resumeUsecase:  mockResumeUS,
		vacancyUsecase: mockVacancyUS,
		companyUsecase: mockCompanyUS,
		strg:           mockRepo,
	}
	return mockRepo, mockResumeUS, mockVacancyUS, mockCompanyUS, mockUserUS, usecase
}

//CreateChatAndTechMes(models.Response) (*models.Response, error)
func TestResponseCreate(t *testing.T) {
	mockRepo, mockResumeUS, mockVacancyUS, _, _, usecase := beforeTest(t)
	var testResponse = models.Response{
		ID:        ID,
		ResumeID:  ID,
		VacancyID: ID,
		Status:    "sent",
	}

	testResponse.Initial = common.Candidate
	mockResumeUS.On("GetById", ID).Return(&testResume, nil)
	mockRepo.On("Create", mock.Anything).Return(&testResponse, nil)
	answerCorrect, errNil := usecase.Create(testResponse)

	assert.Nil(t, errNil)
	assert.Equal(t, *answerCorrect, testResponse)

	mockResumeUS.On("GetById", uuid.Nil).Return(nil, assert.AnError)
	testResponse.ResumeID = uuid.Nil
	answerWrong, errNotNil := usecase.Create(testResponse)
	assert.Nil(t, answerWrong)
	assert.Error(t, errNotNil)

	testResponse.Initial = common.Employer
	mockVacancyUS.On("GetVacancy", ID).Return(&testVacacy, nil)
	answerCorrect2, errNil2 := usecase.Create(testResponse)

	assert.Nil(t, errNil2)
	assert.Equal(t, *answerCorrect2, testResponse)

	mockVacancyUS.On("GetVacancy", uuid.Nil).Return(nil, assert.AnError)
	testResponse.VacancyID = uuid.Nil
	answerWrong2, errNotNil2 := usecase.Create(testResponse)
	assert.Nil(t, answerWrong2)
	assert.Error(t, errNotNil2)
}

//UpdateStatus(response models.Response, userUpdate string) (*models.Response, error)
func TestResponseUpdateStatus(t *testing.T) {
	mockRepo, mockResumeUS, mockVacancyUS, _, _, usecase := beforeTest(t)
	var testResponse = models.Response{
		ID:        ID,
		ResumeID:  ID,
		VacancyID: ID,
		Status:    "sent",
	}
	answerWrong, err := usecase.UpdateStatus(testResponse, "sent")
	assert.Nil(t, answerWrong, err)

	testResponse.Status = "accept"
	testResponse.Status = common.Employer

	mockRepo.On("GetByID", ID).Return(&testResponse, nil)
	mockRepo.On("UpdateStatus", testResponse).Return(&testResponse, nil)
	mockResumeUS.On("GetById", ID).Return(&testResume, nil)
	answerCorrect, errNil := usecase.UpdateStatus(testResponse, common.Candidate)
	assert.Nil(t, errNil)
	assert.Equal(t, *answerCorrect, testResponse)

	testResponse.ResumeID = uuid.Nil
	mockResumeUS.On("GetById", uuid.Nil).Return(nil, assert.AnError)
	answerWromg2, errNotNil2 := usecase.UpdateStatus(testResponse, common.Candidate)
	assert.Nil(t, answerWromg2)
	assert.Error(t, errNotNil2)

	testResponse.ResumeID = ID
	testResponse.Status = "accept"
	testResponse.Status = common.Candidate

	mockRepo.On("UpdateStatus", testResponse).Return(&testResponse, nil)
	mockVacancyUS.On("GetVacancy", ID).Return(&testVacacy, nil)
	answerCorrect2, errNil2 := usecase.UpdateStatus(testResponse, common.Employer)
	assert.Nil(t, errNil2)
	assert.Equal(t, *answerCorrect2, testResponse)

	testResponse.VacancyID = uuid.Nil
	mockVacancyUS.On("GetVacancy", uuid.Nil).Return(nil, assert.AnError)
	answerWromg3, errNotNil3 := usecase.UpdateStatus(testResponse, common.Employer)
	assert.Nil(t, answerWromg3)
	assert.Error(t, errNotNil3)

	testResponse.VacancyID = ID
}

//GetAllCandidateResponses(candID uuid.UUID) ([]models.ResponseWithTitle, error)
func TestGetAllCandidateResponses(t *testing.T) {
	mockRepo, mockResumeUS, mockVacancyUS, mockCompanyUS, _, usecase := beforeTest(t)
	var testResponse = models.Response{
		ID:        ID,
		ResumeID:  ID,
		VacancyID: ID,
		Status:    "sent",
	}

	listResume := []models.BriefResumeInfo{briefResume}
	listResponse := []models.Response{testResponse}
	mockResumeUS.On("GetAllUserResume", ID).Return(listResume, nil)
	mockRepo.On("GetResumeAllResponse", ID).Return(listResponse, nil)
	mockVacancyUS.On("GetVacancy", ID).Return(&testVacacy, nil)
	mockCompanyUS.On("GetOfficialCompany", ID).Return(&testCompany, nil)

	listRespWithTitle := []models.ResponseWithTitle{respWithTitle}
	answerCorrect, errNill := usecase.GetAllCandidateResponses(ID, []uuid.UUID(nil))
	assert.Nil(t, errNill)
	assert.Equal(t, answerCorrect, listRespWithTitle)
}

func TestGetAllEmployerResponses(t *testing.T) {
	mockRepo, mockResumeUS, mockVacancyUS, mockCompanyUS, _, usecase := beforeTest(t)
	var testResponse = models.Response{
		ID:        ID,
		ResumeID:  ID,
		VacancyID: ID,
		Status:    "sent",
	}

	listResponse := []models.Response{testResponse}
	listVacancy := []models.Vacancy{testVacacy}
	mockVacancyUS.On("GetVacancyList", uint(0), uint(100), ID, vacancy.ByEmpId).Return(listVacancy, nil)
	mockCompanyUS.On("GetOfficialCompany", ID).Return(&testCompany, nil)
	mockRepo.On("GetVacancyAllResponse", ID).Return(listResponse, nil)
	mockResumeUS.On("GetById", ID).Return(&testResume, nil)

	listRespWithTitle := []models.ResponseWithTitle{respWithTitle}
	answerCorrect, errNill := usecase.GetAllEmployerResponses(ID, []uuid.UUID(nil))
	assert.Nil(t, errNill)
	assert.Equal(t, answerCorrect, listRespWithTitle)
}

func TestGetAllResumeWithoutResponse(t *testing.T) {
	briefResume.UserID = uuid.Nil
	briefResume.Email = ""
	briefResume.Name = ""
	briefResume.Surname = ""

	mockRepo, _, _, _, _, usecase := beforeTest(t)
	listResume := []models.Resume{testResume}
	listBriefResume := []models.BriefResumeInfo{briefResume}

	mockRepo.On("GetAllResumeWithoutResponse", ID, ID).Return(listResume, nil)

	answerCorrect, errNill := usecase.GetAllResumeWithoutResponse(ID, ID)
	assert.Nil(t, errNill)
	assert.Equal(t, answerCorrect, listBriefResume)
}

func TestGetAllVacancyWithoutResponse(t *testing.T) {
	mockRepo, _, _, _, _, usecase := beforeTest(t)
	listVacancy := []models.Vacancy{testVacacy}

	mockRepo.On("GetAllVacancyWithoutResponse", ID, ID).Return(listVacancy, nil)

	answerCorrect, errNill := usecase.GetAllVacancyWithoutResponse(ID, ID)
	assert.Nil(t, errNill)
	assert.Equal(t, answerCorrect, listVacancy)
}

func TestNewUsecase(t *testing.T) {
	useCase := NewUsecase(nil, nil, nil, nil, nil, nil)
	assert.Equal(t, useCase, &UseCaseResponse{nil, nil, nil, nil, nil, nil})
}
