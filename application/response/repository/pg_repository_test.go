package repository
//
//import (
//	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
//	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"gitlab.com/slax0rr/go-pg-wrapper/mocks"
//	ormmocks "gitlab.com/slax0rr/go-pg-wrapper/mocks/orm"
//	"testing"
//)
//
//var ID = uuid.New()
//var testResponse = models.Response{
//	ID:         ID,
//	ResumeID:   ID,
//	VacancyID:  ID,
//}
//var testResume = models.Resume{
//	ResumeID:  ID,
//	CandID:    ID,
//	Title:     "ID",
//}
//var listResp = []models.Response{testResponse}
//var listResume = []models.Resume{testResume}
//
//func mockDB() (*mocks.DB, pgRepository) {
//	db := new(mocks.DB)
//	r := pgRepository{db: db}
//	return db, r
//}
//
//func mockQueryResponse(db *mocks.DB) *ormmocks.Query {
//	query := new(ormmocks.Query)
//	mockCall := db.On("Model", mock.AnythingOfType("*models.Response")).Return(query)
//	mockCall.RunFn = func(args mock.Arguments) {
//		response := args[0].(*models.Response)
//		*response = testResponse
//	}
//	return query
//}
//
//func mockQueryResponseList(db *mocks.DB) *ormmocks.Query {
//	query := new(ormmocks.Query)
//	mockCall := db.On("Model", mock.AnythingOfType("*[]models.Response")).Return(query)
//	mockCall.RunFn = func(args mock.Arguments) {
//		respList := args[0].(*[]models.Response)
//		*respList = listResp
//	}
//	return query
//}
//
//func mockQueryResumeList(db *mocks.DB) *ormmocks.Query {
//	query := new(ormmocks.Query)
//	mockCall := db.On("Model", mock.AnythingOfType("*[]models.Resume")).Return(query)
//	mockCall.RunFn = func(args mock.Arguments) {
//		resumeList := args[0].(*[]models.Resume)
//		*resumeList = listResume
//	}
//	return query
//}
//
//func TestResponseGetUserByID(t *testing.T) {
//	db, r := mockDB()
//	query := mockQueryResponse(db)
//
//	query.On("Where", "response_id = ?", ID).Return(query)
//	query.On("Select").Return(nil)
//
//	answerCorrect, err := r.GetByID(ID)
//	assert.Nil(t, err)
//	assert.Equal(t, testResponse, *answerCorrect)
//}
//
//func TestCreateResponse(t *testing.T) {
//	db, r := mockDB()
//	query := mockQueryResponse(db)
//	mockResult := general.MockResult{}
//
//	query.On("Returning", "*").Return(query)
//	query.On("Insert").Return(mockResult, nil)
//	answerCorrect, err := r.Create(testResponse)
//	assert.Nil(t, err)
//	assert.Equal(t, testResponse, *answerCorrect)
//}
//
//func TestUpdateResponse(t *testing.T) {
//	db, r := mockDB()
//	query := mockQueryResponse(db)
//	mockResult := general.MockResult{}
//
//	query.On("WherePK").Return(query)
//	query.On("Returning", "*").Return(query)
//	query.On("UpdateNotZero").Return(mockResult, nil)
//	answerCorrect, err := r.UpdateStatus(testResponse)
//	assert.Nil(t, err)
//	assert.Equal(t, testResponse, *answerCorrect)
//}
//
//func TestGetResumeAllResponse(t *testing.T) {
//	db, r := mockDB()
//	query := mockQueryResponseList(db)
//
//	query.On("Where", "resume_id = ?", ID).Return(query)
//	query.On("Select").Return(nil)
//
//	answerCorrect, err := r.GetResumeAllResponse(ID)
//	assert.Nil(t, err)
//	assert.Equal(t, listResp, answerCorrect)
//}
//
//func TestGetVacancyAllResponse(t *testing.T) {
//	db, r := mockDB()
//	query := mockQueryResponseList(db)
//
//	query.On("Where", "vacancy_id = ?", ID).Return(query)
//	query.On("Select").Return(nil)
//
//	answerCorrect, err := r.GetVacancyAllResponse(ID)
//	assert.Nil(t, err)
//	assert.Equal(t, listResp, answerCorrect)
//}
//
////func TestGetAllResumeWithoutResponse(t *testing.T) {
////	db, r := mockDB()
////	query := mockQueryResumeList(db)
////	mockResult := general.MockResult{}
////
////	list := []models.Resume(nil)
////	queryText := fmt.Sprintf(`select main.resume.* from main.resume left join main.response on main.response.resume_id = main.resume.resume_id where cand_id = '%s' group by main.resume.resume_id having sum(case when vacancy_id = '%s' then 1 else 0 end) = 0`, ID, ID)
////	query.On("Query", &list, queryText).Return(mockResult, nil)
////
////	answerCorrect, err := r.GetAllResumeWithoutResponse(ID, ID)
////	assert.Nil(t, err)
////	assert.Equal(t, listResume, answerCorrect)
////}
