package models

import (
	"github.com/google/uuid"
)

type Recommendation struct {
	ID       uuid.UUID `gorm:"column:rec_id;primaryKey;default:uuid_generate_v4()"`
	Sphere0  int
	Sphere1  int
	Sphere2  int
	Sphere3  int
	Sphere4  int
	Sphere5  int
	Sphere6  int
	Sphere7  int
	Sphere8  int
	Sphere9  int
	Sphere10 int
	Sphere11 int
	Sphere12 int
	Sphere13 int
	Sphere14 int
	Sphere15 int
	Sphere16 int
	Sphere17 int
	Sphere18 int
	Sphere19 int
	Sphere20 int
	Sphere21 int
	Sphere22 int
	Sphere23 int
	Sphere24 int
	Sphere25 int
	Sphere26 int
	Sphere27 int
	Sphere28 int
	Sphere29 int
}

func (r Recommendation) TableName() string {
	return "main.recommendation"
}
