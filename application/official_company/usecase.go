package official_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)



type IUseCaseOfficialCompany interface {
	CreateOfficialCompany(company models.OfficialCompany) (*models.OfficialCompany, error)
	//UpdateResume(resume models.Resume) (*models.Resume, error)
	//GetOfficialCompanyByEmployerId(id string) (*models.OfficialCompany, error)
	GetOfficialCompany(id string) (*models.OfficialCompany, error)
}
