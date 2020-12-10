package server

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/vacancyMicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	vacancy2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/microservises/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type vacServer struct {
	vacUseCase vacancy.IUseCaseVacancy
	vacancy2.UnimplementedVacancyServer
}

func NewVacServer(vacUseCase vacancy.IUseCaseVacancy) vacancy2.VacancyServer {
	return &vacServer{vacUseCase: vacUseCase}
}

func (v *vacServer) GetVacancyTopSpheres(ctx context.Context, sphereCnt *vacancy2.SphereCnt) (*vacancy2.SphereList, error) {
	topSpheres, err := v.vacUseCase.GetVacancyTopSpheres(sphereCnt.Cnt)
	if err != nil {
		return nil, err
	}
	return vacancyMicro.ConvertSphToPbModels(topSpheres), err
}

func (v *vacServer) CreateVacancy(ctx context.Context, req *vacancy2.Vac) (*vacancy2.Vac, error) {
	reqVac := vacancyMicro.ConvertToDbModel(req)
	if req.Sphere == "" {
		reqVac.Sphere = -1
	}
	newVac, err := v.vacUseCase.CreateVacancy(*reqVac)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return vacancyMicro.ConvertToPbModel(newVac), nil
}

func (v *vacServer) UpdateVacancy(ctx context.Context, req *vacancy2.Vac) (*vacancy2.Vac, error) {
	reqVac := vacancyMicro.ConvertToDbModel(req)
	if req.Sphere == "" {
		reqVac.Sphere = -1
	}
	newVac, err := v.vacUseCase.UpdateVacancy(*reqVac)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return vacancyMicro.ConvertToPbModel(newVac), nil
}

func (v *vacServer) GetVacancy(ctx context.Context, vacId *vacancy2.VacId) (*vacancy2.Vac, error) {
	id, _ := uuid.Parse(vacId.Id)
	newVac, err := v.vacUseCase.GetVacancy(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return vacancyMicro.ConvertToPbModel(newVac), nil
}

func (v *vacServer) GetVacancyList(ctx context.Context, params *vacancy2.VacListParams) (*vacancy2.VacList, error) {
	entityId, _ := uuid.Parse(params.EntityId)
	vacList, err := v.vacUseCase.GetVacancyList(uint(params.Start), uint(params.Limit), entityId, int(params.EntityType))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return vacancyMicro.ConvertToListDbModels(vacList), err
}

func (v *vacServer) AddRecommendation(ctx context.Context, params *vacancy2.AddRecParams) (*vacancy2.Empty, error) {
	userId, _ := uuid.Parse(params.UserId)
	err := v.vacUseCase.AddRecommendation(userId, int(params.Sphere))
	return &vacancy2.Empty{}, err
}

func (v *vacServer) GetRecommendation(ctx context.Context, params *vacancy2.GetRecParams) (*vacancy2.VacList, error) {
	userId, _ := uuid.Parse(params.UserId)
	vacList, err := v.vacUseCase.GetRecommendation(userId, int(params.Start), int(params.Limit))
	return vacancyMicro.ConvertToListDbModels(vacList), err
}

func (v *vacServer) SearchVacancies(ctx context.Context, params *vacancy2.SearchParams) (*vacancy2.VacList, error) {
	searchParams := models.VacancySearchParams{
		KeyWords:        params.KeyWords,
		SalaryMin:       int(params.SalaryMin),
		SalaryMax:       int(params.SalaryMax),
		Gender:          params.Gender,
		ExperienceMonth: vacancyMicro.ConvertToIntListPbModels(params.ExpList),
		Employment:      vacancyMicro.ConvertToStringListPbModels(params.EmpList),
		EducationLevel:  vacancyMicro.ConvertToStringListPbModels(params.EdList),
		CareerLevel:     vacancyMicro.ConvertToStringListPbModels(params.CarList),
		Sphere:          vacancyMicro.ConvertToIntListPbModels(params.SpheresList),
		AreaSearch:      vacancyMicro.ConvertToStringListPbModels(params.AreaList),
		OrderBy:         params.OrderBy,
		ByAsc:           params.ByAsc,
		DaysFromNow:     int(params.DaysFromNow),
		StartDate:       params.StartDate,
	}
	vacList, err := v.vacUseCase.SearchVacancies(searchParams)
	return vacancyMicro.ConvertToListDbModels(vacList), err
}
