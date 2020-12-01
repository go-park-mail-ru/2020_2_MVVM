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

func (g *gRPCVacClient) CreateVacancy(vacancy models.Vacancy) (*models.Vacancy, error) {
	newVac, err := g.client.CreateVacancy(g.ctx, &api.Vac{Id: vacancy.ID.String(), Title: vacancy.Title, SalaryMin: uint32(vacancy.SalaryMin), SalaryMax: uint32(vacancy.SalaryMax), AreaSearch: vacancy.AreaSearch,
		Description: vacancy.Description, Requirements: vacancy.Requirements, Duties: vacancy.Duties, Skills: vacancy.Skills, Employment: vacancy.Employment,
		ExperienceMonth: uint32(vacancy.ExperienceMonth), Location: vacancy.Location, CareerLevel: vacancy.CareerLevel, EducationLevel: vacancy.EducationLevel,
		EmpPhone: vacancy.EmpPhone, EmpEmail: vacancy.EmpEmail, Gender: vacancy.Gender})
	if newVac == nil {
		return nil, err
	}
	vacResp := models.Vacancy{Title: newVac.Title, SalaryMin: int(newVac.SalaryMin), SalaryMax: int(newVac.SalaryMax), AreaSearch: newVac.AreaSearch,
		Description: newVac.Description, Requirements: newVac.Requirements, Duties: newVac.Duties, Skills: newVac.Skills, Employment: newVac.Employment,
		ExperienceMonth: int(newVac.ExperienceMonth), Location: newVac.Location, CareerLevel: newVac.CareerLevel, EducationLevel: newVac.EducationLevel,
		EmpPhone: newVac.EmpPhone, EmpEmail: newVac.EmpEmail, Gender: newVac.Gender}
	return &vacResp, err
}

func (g *gRPCVacClient) UpdateVacancy(vacancy models.Vacancy) (*models.Vacancy, error) {
	//panic("implement me")
	return nil, nil
}

func (g *gRPCVacClient) GetVacancy(vacId uuid.UUID) (*models.Vacancy, error) {
	id := api.VacId{Id: vacId.String()}
	newVac, err := g.client.GetVacancy(g.ctx, &id)
	if newVac == nil {
		return nil, err
	}
	vacResp := models.Vacancy{Title: newVac.Title, SalaryMin: int(newVac.SalaryMin), SalaryMax: int(newVac.SalaryMax), AreaSearch: newVac.AreaSearch,
		Description: newVac.Description, Requirements: newVac.Requirements, Duties: newVac.Duties, Skills: newVac.Skills, Employment: newVac.Employment,
		ExperienceMonth: int(newVac.ExperienceMonth), Location: newVac.Location, CareerLevel: newVac.CareerLevel, EducationLevel: newVac.EducationLevel,
		EmpPhone: newVac.EmpPhone, EmpEmail: newVac.EmpEmail, Gender: newVac.Gender}
	return nil, err
}

func (g *gRPCVacClient) GetVacancyList(u uint, u2 uint, uuid uuid.UUID, i int) ([]models.Vacancy, error) {
	return nil, nil
}

func (g *gRPCVacClient) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	return nil, nil
}

func (g *gRPCVacClient) AddRecommendation(uuid uuid.UUID, i int) error {
	return nil
}

func (g gRPCVacClient) GetRecommendation(uuid uuid.UUID, i int, i2 int) ([]models.Vacancy, error) {
	return nil, nil
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
