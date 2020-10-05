package resume

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type IUseCaseResume interface {
	CreateResume(resume models.Resume) (*models.Resume, error)
	UpdateResume(resume models.Resume) (*models.Resume, error)
	GetResume(id string) (*models.Resume, error)
	GetResumeList(begin, end uint) ([]models.Resume, error)
}
