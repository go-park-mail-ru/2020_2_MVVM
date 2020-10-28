package official_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type OfficialCompanyRepository interface {
	CreateOfficialCompany(resume models.OfficialCompany) (*models.OfficialCompany, error)
	//UpdateEducation(id uuid.UUID, updResume *models.Resume) (*models.Resume, error)
	//GetOfficialCompanyByEmployerId(id string) (*models.OfficialCompany, error)
	GetOfficialCompanyById(id string) (*models.OfficialCompany, error)
}
