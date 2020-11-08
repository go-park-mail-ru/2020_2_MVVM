package custom_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type CustomCompanyRepository interface {
	CreateCustomCompany(resume models.CustomCompany) (*models.CustomCompany, error)
	//Drop(id uuid.UUID, updResume *models.Resume) (*models.Resume, error)
	GetCustomCompanyById(id string) (*models.CustomCompany, error)
}
