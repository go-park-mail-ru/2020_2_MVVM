package official_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseOfficialCompany interface {
	CreateOfficialCompany(models.OfficialCompany, uuid.UUID) (*models.OfficialCompany, error)
	GetMineCompany(uuid.UUID) (*models.OfficialCompany, error)
	GetOfficialCompany(uuid.UUID) (*models.OfficialCompany, error)
	GetCompaniesList(uint, uint) ([]models.OfficialCompany, error)
	SearchCompanies(models.CompanySearchParams) ([]models.OfficialCompany, error)
}
