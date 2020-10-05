package resume

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type ResumeRepository interface {
	CreateResume(resume models.Resume) (models.Resume, error)
	//UpdateResume(id uint) (models.Resume, error)
	GetResumeById(id string) (models.Resume, error)
	GetResumeByName(name string) (models.Resume, error)
}

