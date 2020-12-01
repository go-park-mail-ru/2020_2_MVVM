package server

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/api"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/vacancyMicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type vacServer struct {
	vacUseCase vacancy.IUseCaseVacancy
	api.UnimplementedVacancyServer
}

func NewVacServer(vacUseCase vacancy.IUseCaseVacancy) api.VacancyServer {
	return &vacServer{vacUseCase: vacUseCase}
}

func (v *vacServer) CreateVacancy(ctx context.Context, req *api.Vac) (*api.Vac, error) {
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

func (v *vacServer) GetVacancy(ctx context.Context, vacId *api.VacId) (*api.Vac, error) {
	id, _ := uuid.Parse(vacId.Id)
	newVac, err := v.vacUseCase.GetVacancy(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return vacancyMicro.ConvertToPbModel(newVac), nil
}

func (v *vacServer) GetVacancyList(ctx context.Context, params *api.VacListParams) (*api.VacList, error) {
	entityId, _ := uuid.Parse(params.EntityId)
	vacList, err := v.vacUseCase.GetVacancyList(uint(params.Start), uint(params.Limit), entityId, int(params.EntityType))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return vacancyMicro.ConvertToListDbModels(vacList), err
}

func (v *vacServer) AddRecommendation(ctx context.Context, params *api.AddRecParams) (*api.Empty, error) {
	userId, _ := uuid.Parse(params.UserId)
	err := v.vacUseCase.AddRecommendation(userId, int(params.Sphere))
	return &api.Empty{}, err
}

func (v *vacServer) GetRecommendation(ctx context.Context, params *api.GetRecParams) (*api.VacList, error) {
	userId, _ := uuid.Parse(params.UserId)
	vacList, err := v.vacUseCase.GetRecommendation(userId, int(params.Start), int(params.Limit))
	return vacancyMicro.ConvertToListDbModels(vacList), err
}

func (v *vacServer) SearchVacancies(ctx context.Context, params *api.SearchParams) (*api.VacList, error) {
	searchParams := models.VacancySearchParams{
		KeyWords:        params.KeyWords,
		SalaryMin:       int(params.SalaryMin),
		SalaryMax:       int(params.SalaryMax),
		Gender:          params.Gender,
		ExperienceMonth: vacancyMicro.ConvertToIntListPbModels(params.ExpList),
		Employment:      vacancyMicro.ConvertToStringListPbModels(params.EmpList),
		EducationLevel:  vacancyMicro.ConvertToStringListPbModels(params.EdList),
		CareerLevel:     vacancyMicro.ConvertToStringListPbModels(params.CarList),
		Spheres:         vacancyMicro.ConvertToIntListPbModels(params.SpheresList),
		AreaSearch:      vacancyMicro.ConvertToStringListPbModels(params.AreaList),
		OrderBy:         params.OrderBy,
		ByAsc:           params.ByAsc,
		DaysFromNow:     int(params.DaysFromNow),
		StartDate:       params.StartDate,
	}
	vacList, err := v.vacUseCase.SearchVacancies(searchParams)
	return vacancyMicro.ConvertToListDbModels(vacList), err
}
