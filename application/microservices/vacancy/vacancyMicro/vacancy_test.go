package vacancyMicro

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/microservises/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToDbModel(t *testing.T) {
	res1 := ConvertToDbModel(nil)
	assert.Nil(t, res1)
	pbModel := vacancy.Vac{Title: "test", Sphere: "1"}
	res2 := ConvertToDbModel(&pbModel)
	assert.Equal(t, *res2, models.Vacancy{Title: "test", Sphere: 1})
}

func TestConvertToPbModel(t *testing.T) {
	res1 := ConvertToPbModel(nil)
	assert.Equal(t, &vacancy.Vac{}, res1)
	dbModel := models.Vacancy{Title: "test", Sphere: 0}
	res2 := ConvertToPbModel(&dbModel)
	assert.Equal(t, res2.Title, dbModel.Title)
}

func TestConvertToListPbModels(t *testing.T) {
	res1 := ConvertToListPbModels(nil)
	assert.Nil(t, res1)
	var pbModels = new(vacancy.VacList)
	pbModels.List = make([]*vacancy.Vac, 1)
	pbModels.List[0] = &vacancy.Vac{Title: "test"}
	res2 := ConvertToListPbModels(pbModels)
	assert.Equal(t, res2, []models.Vacancy{{Title: "test"}})
}

func TestConvertToListDbModels(t *testing.T) {
	res1 := ConvertToListDbModels(nil)
	assert.Equal(t, res1, &vacancy.VacList{})
	var pbModels = new(vacancy.VacList)
	pbModels.List = make([]*vacancy.Vac, 1)
	pbModels.List[0] = &vacancy.Vac{Title: "test"}
	res2 := ConvertToListDbModels([]models.Vacancy{{Title: "test"}})
	assert.Equal(t, pbModels.List[0].Title, res2.List[0].Title)
}

func TestConvertSphToPbModels(t *testing.T) {
	res1 := ConvertSphToPbModels(nil)
	assert.Equal(t, &vacancy.SphereList{}, res1)
	var pbModels = new(vacancy.SphereList)
	pbModels.List = make([]*vacancy.Sphere, 1)
	pbModels.List[0] = &vacancy.Sphere{SphereIdx: 0, VacCnt: 1}
	res2 := ConvertSphToPbModels([]models.Sphere{{0,1}})
	assert.Equal(t, res2.List, pbModels.List)
}

func TestConvertSphToDbModels(t *testing.T) {
	res1 := ConvertSphToDbModels(nil)
	assert.Nil(t, res1)
	var pbModels = new(vacancy.SphereList)
	pbModels.List = make([]*vacancy.Sphere, 1)
	pbModels.List[0] = &vacancy.Sphere{SphereIdx: 0, VacCnt: 1}
	res2 := ConvertSphToDbModels(pbModels)
	assert.Equal(t, models.Sphere{Sph: 0, VacCnt: 1}, res2[0])
}
