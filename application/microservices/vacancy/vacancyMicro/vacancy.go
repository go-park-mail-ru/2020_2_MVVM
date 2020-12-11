package vacancyMicro

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/microservises/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"strconv"
)

type gRPCVacClient struct {
	client vacancy.VacancyClient
	gConn  *grpc.ClientConn
	logger common.Logger
	ctx    context.Context
}

func ConvertToDbModel(pbModel *vacancy.Vac) *models.Vacancy {
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

func ConvertToPbModel(dbModel *models.Vacancy) *vacancy.Vac {
	if dbModel == nil {
		return &vacancy.Vac{}
	}
	return &vacancy.Vac{ID: dbModel.ID.String(), EmpID: dbModel.EmpID.String(), CompID: dbModel.CompID.String(), Title: dbModel.Title,
		SalaryMin: uint32(dbModel.SalaryMin), SalaryMax: uint32(dbModel.SalaryMax), AreaSearch: dbModel.AreaSearch, Description: dbModel.Description, Requirements: dbModel.Requirements,
		Duties: dbModel.Duties, Skills: dbModel.Skills, Employment: dbModel.Employment, ExperienceMonth: uint32(dbModel.ExperienceMonth), Location: dbModel.Location,
		CareerLevel: dbModel.CareerLevel, EducationLevel: dbModel.EducationLevel, EmpPhone: dbModel.EmpPhone, EmpEmail: dbModel.EmpEmail, Gender: dbModel.Gender,
		Sphere: strconv.Itoa(dbModel.Sphere), Avatar: dbModel.Avatar, DateCreate: dbModel.DateCreate}
}

func ConvertToListPbModels(pbModels *vacancy.VacList) []models.Vacancy {
	var vacList []models.Vacancy
	if pbModels == nil {
		return nil
	}
	vacList = make([]models.Vacancy, len(pbModels.List))
	for i, pbModel := range pbModels.List {
		vacList[i] = *ConvertToDbModel(pbModel)
	}
	return vacList
}

func ConvertToListDbModels(dbModels []models.Vacancy) *vacancy.VacList {
	if dbModels == nil {
		return &vacancy.VacList{}
	}
	var pbModels = new(vacancy.VacList)
	pbModels.List = make([]*vacancy.Vac, len(dbModels))
	for i, dbModel := range dbModels {
		pbModels.List[i] = ConvertToPbModel(&dbModel)
	}
	return pbModels
}

func ConvertSphToPbModels(dbModels []models.Sphere) *vacancy.SphereList {
	if dbModels == nil {
		return &vacancy.SphereList{}
	}
	pbList := new(vacancy.SphereList)
	pbList.List = make([]*vacancy.Sphere, len(dbModels))
	for i, e := range dbModels {
		pbList.List[i] = &vacancy.Sphere{SphereIdx: e.Sph, VacCnt: e.VacCnt}
	}
	return pbList
}

func ConvertSphToDbModels(pbModels *vacancy.SphereList) []models.Sphere {
	if pbModels == nil {
		return nil
	}
	dbList := make([]models.Sphere, len(pbModels.List))
	for i, e := range pbModels.List {
		dbList[i] = models.Sphere{Sph: e.SphereIdx, VacCnt: e.VacCnt}
	}
	return dbList
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
	newVac, err := g.client.GetVacancy(g.ctx, &vacancy.VacId{Id: vacId.String()})
	return ConvertToDbModel(newVac), err
}

func (g *gRPCVacClient) GetVacancyList(start uint, limit uint, entityId uuid.UUID, entityType int) ([]models.Vacancy, error) {
	vacList, err := g.client.GetVacancyList(g.ctx, &vacancy.VacListParams{Start: uint32(start), Limit: uint32(limit),
		EntityId: entityId.String(), EntityType: int32(entityType)})
	return ConvertToListPbModels(vacList), err
}

func ConvertToStringListPbModels(pbModels *vacancy.StringArr) []string {
	if pbModels == nil {
		return nil
	}
	strList := make([]string, len(pbModels.Elem))
	for i, pbModel := range pbModels.Elem {
		strList[i] = pbModel
	}
	return strList
}
func ConvertToIntListPbModels(pbModels *vacancy.IntArr) []int {
	if pbModels == nil {
		return nil
	}
	intList := make([]int, len(pbModels.Elem))
	for i, pbModel := range pbModels.Elem {
		intList[i] = int(pbModel)
	}
	return intList
}

func ConvertIntListToPbModel(slice []int) *vacancy.IntArr {
	if slice == nil {
		return nil
	}
	intArr := new(vacancy.IntArr)
	intArr.Elem = make([]int32, len(slice))
	for i, e := range slice {
		intArr.Elem[i] = int32(e)
	}
	return intArr
}

func ConvertStringListToPbModel(slice []string) *vacancy.StringArr {
	if slice == nil {
		return nil
	}
	stringArr := new(vacancy.StringArr)
	stringArr.Elem = make([]string, len(slice))
	for i, s := range slice {
		stringArr.Elem[i] = s
	}
	return stringArr
}

func (g *gRPCVacClient) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	searchParams := &vacancy.SearchParams{
		AreaList:    ConvertStringListToPbModel(params.AreaSearch),
		ExpList:     ConvertIntListToPbModel(params.ExperienceMonth),
		EdList:      ConvertStringListToPbModel(params.EducationLevel),
		EmpList:     ConvertStringListToPbModel(params.Employment),
		CarList:     ConvertStringListToPbModel(params.CareerLevel),
		SpheresList: ConvertIntListToPbModel(params.Sphere),
		DaysFromNow: int32(params.DaysFromNow),
		KeyWords:    params.KeyWords,
		SalaryMax:   int32(params.SalaryMax),
		SalaryMin:   int32(params.SalaryMin),
		Gender:      params.Gender,
		OrderBy:     params.OrderBy,
		ByAsc:       params.ByAsc,
		KeyWordsGeo: params.KeywordsGeo,
	}
	vacList, err := g.client.SearchVacancies(g.ctx, searchParams)

	return ConvertToListPbModels(vacList), err
}

func (g *gRPCVacClient) AddRecommendation(userId uuid.UUID, sphere int) error {
	_, err := g.client.AddRecommendation(g.ctx, &vacancy.AddRecParams{UserId: userId.String(), Sphere: int32(sphere)})
	return err
}

func (g *gRPCVacClient) GetRecommendation(userId uuid.UUID, start int, limit int) ([]models.Vacancy, error) {
	vacList, err := g.client.GetRecommendation(g.ctx, &vacancy.GetRecParams{UserId: userId.String(), Start: int32(start), Limit: int32(limit)})
	return ConvertToListPbModels(vacList), err
}

func (g *gRPCVacClient) GetVacancyTopSpheres(sphereCnt int32) ([]models.Sphere, *models.VacTopCnt, error) {
	topInfo, err := g.client.GetVacancyTopSpheres(g.ctx, &vacancy.SphereCnt{Cnt: sphereCnt})
	if err != nil {
		return nil, nil, err
	}
	vacInfo := models.VacTopCnt{}
	if topInfo.VacInfo != nil {
		vacInfo.NewVacCnt = topInfo.VacInfo.NewVacCnt
		vacInfo.AllVacCnt = topInfo.VacInfo.AllVacCnt
	}
	return ConvertSphToDbModels(topInfo.SphereInfo), &vacInfo, err
}

func NewVacClient(host string, port int, logger common.Logger) (VacClient, error) {
	gConn, err := grpc.Dial(
		host+":"+strconv.Itoa(port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &gRPCVacClient{client: vacancy.NewVacancyClient(gConn), gConn: gConn, logger: logger, ctx: context.Background()}, nil
}
