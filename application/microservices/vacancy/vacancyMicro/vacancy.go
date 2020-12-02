package vacancyMicro

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/api"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"strconv"
)

type gRPCVacClient struct {
	client api.VacancyClient
	gConn  *grpc.ClientConn
	logger common.Logger
	ctx    context.Context
}

func ConvertToDbModel(pbModel *api.Vac) *models.Vacancy {
	if pbModel == nil {
		return nil
	}
	vacId, _ := uuid.Parse(pbModel.ID)
	empId, _ := uuid.Parse(pbModel.EmpID)
	compId, _ := uuid.Parse(pbModel.CompID)
	vacSphere := 0
	if pbModel.Sphere != "" {
		vacSphere, _ = strconv.Atoi(pbModel.Sphere)
	}
	vacResp := models.Vacancy{ID: vacId, EmpID: empId, CompID: compId, Title: pbModel.Title, SalaryMin: int(pbModel.SalaryMin), SalaryMax: int(pbModel.SalaryMax),
		AreaSearch: pbModel.AreaSearch, Description: pbModel.Description, Requirements: pbModel.Requirements, Duties: pbModel.Duties, Skills: pbModel.Skills,
		Employment: pbModel.Employment, ExperienceMonth: int(pbModel.ExperienceMonth), Location: pbModel.Location, CareerLevel: pbModel.CareerLevel,
		EducationLevel: pbModel.EducationLevel, EmpPhone: pbModel.EmpPhone, EmpEmail: pbModel.EmpEmail, Gender: pbModel.Gender, Sphere: vacSphere, Avatar: pbModel.Avatar,
		DateCreate: pbModel.DateCreate}
	return &vacResp
}

func ConvertToPbModel(dbModel *models.Vacancy) *api.Vac {
	if dbModel == nil {
		return &api.Vac{}
	}
	return &api.Vac{ID: dbModel.ID.String(), EmpID: dbModel.EmpID.String(), CompID: dbModel.CompID.String(), Title: dbModel.Title,
		SalaryMin: uint32(dbModel.SalaryMin), SalaryMax: uint32(dbModel.SalaryMax), AreaSearch: dbModel.AreaSearch, Description: dbModel.Description, Requirements: dbModel.Requirements,
		Duties: dbModel.Duties, Skills: dbModel.Skills, Employment: dbModel.Employment, ExperienceMonth: uint32(dbModel.ExperienceMonth), Location: dbModel.Location,
		CareerLevel: dbModel.CareerLevel, EducationLevel: dbModel.EducationLevel, EmpPhone: dbModel.EmpPhone, EmpEmail: dbModel.EmpEmail, Gender: dbModel.Gender,
		Sphere: strconv.Itoa(dbModel.Sphere), Avatar: dbModel.Avatar, DateCreate: dbModel.DateCreate}
}

func ConvertToListPbModels(pbModels *api.VacList) []models.Vacancy {
	var vacList []models.Vacancy
	if pbModels == nil {
		return nil
	}
	for _, pbModel := range pbModels.List {
		vacList = append(vacList, *ConvertToDbModel(pbModel))
	}
	return vacList
}

func ConvertToListDbModels(dbModels []models.Vacancy) *api.VacList {
	if dbModels == nil {
		return &api.VacList{}
	}
	var pbModels = new(api.VacList)
	for _, dbModel := range dbModels {
		pbModels.List = append(pbModels.List, ConvertToPbModel(&dbModel))
	}
	return pbModels
}

func (g *gRPCVacClient) CreateVacancy(vacancy models.Vacancy) (*models.Vacancy, error) {
	newVac, err := g.client.CreateVacancy(g.ctx, ConvertToPbModel(&vacancy))
	return ConvertToDbModel(newVac), err
}

func (g *gRPCVacClient) UpdateVacancy(vacancy models.Vacancy) (*models.Vacancy, error) {
	newVac, err := g.client.UpdateVacancy(g.ctx, ConvertToPbModel(&vacancy))
	return ConvertToDbModel(newVac), err
}

func (g *gRPCVacClient) GetVacancy(vacId uuid.UUID) (*models.Vacancy, error) {
	newVac, err := g.client.GetVacancy(g.ctx, &api.VacId{Id: vacId.String()})
	return ConvertToDbModel(newVac), err
}

func (g *gRPCVacClient) GetVacancyList(start uint, limit uint, entityId uuid.UUID, entityType int) ([]models.Vacancy, error) {
	vacList, err := g.client.GetVacancyList(g.ctx, &api.VacListParams{Start: uint32(start), Limit: uint32(limit),
		EntityId: entityId.String(), EntityType: int32(entityType)})
	return ConvertToListPbModels(vacList), err
}

func ConvertToStringListPbModels(pbModels *api.StringArr) []string {
	if pbModels == nil {
		return nil
	}
	strList := make([]string, len(pbModels.Elem))
	for _, pbModel := range pbModels.Elem {
		strList = append(strList, pbModel)
	}
	return strList
}
func ConvertToIntListPbModels(pbModels *api.IntArr) []int {
	if pbModels == nil {
		return nil
	}
	intList := make([]int, len(pbModels.Elem))
	for _, pbModel := range pbModels.Elem {
		intList = append(intList, int(pbModel))
	}
	return intList
}

func ConvertIntListToPbModel(slice []int) *api.IntArr {
	if slice == nil {
		return nil
	}
	intArr := new(api.IntArr)
	for _, e := range slice {
		intArr.Elem = append(intArr.Elem, int32(e))
	}
	return intArr
}

func ConvertStringListToPbModel(slice []string) *api.StringArr {
	if slice == nil {
		return nil
	}
	stringArr := new(api.StringArr)
	for _, s := range slice {
		stringArr.Elem = append(stringArr.Elem, s)
	}
	return stringArr
}

func (g *gRPCVacClient) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	searchParams := &api.SearchParams{
		AreaList:    ConvertStringListToPbModel(params.AreaSearch),
		ExpList:     ConvertIntListToPbModel(params.ExperienceMonth),
		EdList:      ConvertStringListToPbModel(params.EducationLevel),
		EmpList:     ConvertStringListToPbModel(params.Employment),
		CarList:     ConvertStringListToPbModel(params.CareerLevel),
		SpheresList: ConvertIntListToPbModel(params.Spheres),
		DaysFromNow: int32(params.DaysFromNow),
		KeyWords:    params.KeyWords,
		SalaryMax:   int32(params.SalaryMax),
		SalaryMin:   int32(params.SalaryMin),
		Gender:      params.Gender,
		OrderBy:     params.OrderBy,
		ByAsc:       params.ByAsc,
	}
	vacList, err := g.client.SearchVacancies(g.ctx, searchParams)

	return ConvertToListPbModels(vacList), err
}

func (g *gRPCVacClient) AddRecommendation(userId uuid.UUID, sphere int) error {
	_, err := g.client.AddRecommendation(g.ctx, &api.AddRecParams{UserId: userId.String(), Sphere: int32(sphere)})
	return err
}

func (g gRPCVacClient) GetRecommendation(userId uuid.UUID, start int, limit int) ([]models.Vacancy, error) {
	vacList, err := g.client.GetRecommendation(g.ctx, &api.GetRecParams{UserId: userId.String(), Start: int32(start), Limit: int32(limit)})
	return ConvertToListPbModels(vacList), err
}

func NewVacClient(host string, port int, logger common.Logger) (VacClient, error) {
	gConn, err := grpc.Dial(
		host+":"+strconv.Itoa(port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &gRPCVacClient{client: api.NewVacancyClient(gConn), gConn: gConn, logger: logger, ctx: context.Background()}, nil
}
