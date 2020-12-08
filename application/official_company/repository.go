package official_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

type OfficialCompanyRepository interface {
	CreateOfficialCompany(models.OfficialCompany, uuid.UUID) (*models.OfficialCompany, error)
	UpdateOfficialCompany(models.OfficialCompany, uuid.UUID) (*models.OfficialCompany, error)
	DeleteOfficialCompany(uuid.UUID, uuid.UUID) error
	GetCompaniesList(uint, uint) ([]models.OfficialCompany, error)
	GetMineCompany(uuid.UUID) (*models.OfficialCompany, error)
	GetOfficialCompany(uuid.UUID) (*models.OfficialCompany, error)
	SearchCompanies(models.CompanySearchParams) ([]models.OfficialCompany, error)
	GetAllCompaniesNames () ([]models.BriefCompany, error)
}
