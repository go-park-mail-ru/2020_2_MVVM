package custom_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type IUseCaseCustomCompany interface {
	CreateCustomCompany(company models.CustomCompany) (*models.CustomCompany, error)
	//UpdateResume(resume models.Resume) (*models.Resume, error)
	GetCustomCompany(id string) (*models.CustomCompany, error)
}
