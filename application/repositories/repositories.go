package repositories

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type EducationRepository interface {
	CreateEducation(education models.Education) ()
}
