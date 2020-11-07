package custom_company

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type UseCase interface {
	CreateCustomCompany(company models.CustomCompany) (*models.CustomCompany, error)
	//Update(resume models.Resume) (*models.Resume, error)
	GetCustomCompany(id string) (*models.CustomCompany, error)
}
