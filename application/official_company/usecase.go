package official_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type IUseCaseOfficialCompany interface {
	CreateOfficialCompany(models.OfficialCompany, string) (*models.OfficialCompany, error)
	GetMineCompany(string) (*models.OfficialCompany, error)
	GetOfficialCompany(string) (*models.OfficialCompany, error)
	GetCompaniesList(uint, uint) ([]models.OfficialCompany, error)
	SearchCompanies(models.CompanySearchParams) ([]models.OfficialCompany, error)
}
