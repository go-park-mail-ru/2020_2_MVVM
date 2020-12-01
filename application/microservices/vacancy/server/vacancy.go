package server

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/api"
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
	vacReq := models.Vacancy{Title: req.Title, SalaryMin: int(req.SalaryMin), SalaryMax: int(req.SalaryMax), AreaSearch: req.AreaSearch,
		Description: req.Description, Requirements: req.Requirements, Duties: req.Duties, Skills: req.Skills, Employment: req.Employment,
		ExperienceMonth: int(req.ExperienceMonth), Location: req.Location, CareerLevel: req.CareerLevel, EducationLevel: req.EducationLevel,
		EmpPhone: req.EmpPhone, EmpEmail: req.EmpEmail, Gender: req.Gender}
	if req.Sphere == "" {
		vacReq.Sphere = -1
	}
	newVac, err := v.vacUseCase.CreateVacancy(vacReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.Vac{Id: newVac.ID.String(), Title: newVac.Title, SalaryMin: uint32(newVac.SalaryMin), SalaryMax: uint32(newVac.SalaryMax), AreaSearch: newVac.AreaSearch,
		Description: newVac.Description, Requirements: newVac.Requirements, Duties: newVac.Duties, Skills: newVac.Skills, Employment: newVac.Employment,
		ExperienceMonth: uint32(newVac.ExperienceMonth), Location: newVac.Location, CareerLevel: newVac.CareerLevel, EducationLevel: newVac.EducationLevel,
		EmpPhone: newVac.EmpPhone, EmpEmail: newVac.EmpEmail, Gender: newVac.Gender}, nil
}

func (v *vacServer) GetVacancy(ctx context.Context, vacId *api.VacId) (*api.Empty, error) {
	id, _ := uuid.Parse(vacId.Id)
	newVac, err := v.vacUseCase.GetVacancy(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.Vac{Id: newVac.ID.String(), Title: newVac.Title, SalaryMin: uint32(newVac.SalaryMin), SalaryMax: uint32(newVac.SalaryMax), AreaSearch: newVac.AreaSearch,
		Description: newVac.Description, Requirements: newVac.Requirements, Duties: newVac.Duties, Skills: newVac.Skills, Employment: newVac.Employment,
		ExperienceMonth: uint32(newVac.ExperienceMonth), Location: newVac.Location, CareerLevel: newVac.CareerLevel, EducationLevel: newVac.EducationLevel,
		EmpPhone: newVac.EmpPhone, EmpEmail: newVac.EmpEmail, Gender: newVac.Gender}, nil
}
