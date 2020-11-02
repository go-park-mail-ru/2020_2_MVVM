package official_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type OfficialCompanyRepository interface {
	CreateOfficialCompany(resume models.OfficialCompany) (*models.OfficialCompany, error)
	GetCompaniesList(uint, uint) ([]models.OfficialCompany, error)
	GetMineCompany(empId uuid.UUID) (*models.OfficialCompany, error)
	GetOfficialCompany(compId uuid.UUID) (*models.OfficialCompany, error)
}
