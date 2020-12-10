package models

type Sphere struct {
	Sph    int32 `gorm:"column:sphere_idx" json:"sphere_idx"`
	VacCnt int32 `gorm:"column:sphere_cnt" json:"vac_cnt"`
}

func (s Sphere) TableName() string {
	return "main.sphere"
}
