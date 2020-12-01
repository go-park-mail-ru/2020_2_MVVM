package vacancy

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

//len = 29
var Spheres = []string{
	"Автомобильный бизнес",
	"Гостиницы, рестораны, общепит, кейтеринг",
	"Государственные организации",
	"Добывающая отрасль",
	"ЖКХ",
	"Информационные технологии, системная интеграция, интернет",
	"Лесная промышленность, деревообработка",
	"Медицина, фармацевтика, аптеки",
	"Металлургия, металлообработка",
	"Нефть и газ",
	"Образовательные учреждения",
	"Общественная деятельность, партии, благотворительность, НКО",
	"Перевозки, логистика, склад, ВЭД",
	"Продукты питания",
	"Промышленное оборудование, техника, станки и комплектующие",
	"Розничная торговля",
	"СМИ, маркетинг, реклама, BTL, PR, дизайн, продюсирование",
	"Сельское хозяйство",
	"Строительство, недвижимость, эксплуатация, проектирование",
	"Телекоммуникации, связь",
	"Товары народного потребления (непищевые)",
	"Тяжелое машиностроение",
	"Управление многопрофильными активами",
	"Услуги для бизнеса",
	"Услуги для населения",
	"Финансовый сектор",
	"Химическое производство, удобрения",
	"Электроника, приборостроение, бытовая техника, компьютеры и оргтехника",
	"Энергетика",
}

var CountSpheres = 29

type Pair struct {
	SphereInd int
	Score     int
}

type Resp struct {
	Vacancy *models.Vacancy `json:"vacancy"`
}

type RespList struct {
	Vacancies []models.Vacancy `json:"vacancyList"`
}

type VacRequest struct {
	Id              string `json:"vac_id,uuid" valid:"-"`
	Avatar          string `json:"avatar" valid:"-"`
	Title           string `json:"title" binding:"required" valid:"stringlength(4|128)~название вакансии должно быть от 4 до 128 символов в длину."`
	Gender          string `json:"gender" valid:"-"`
	SalaryMin       int    `json:"salary_min" valid:"-"`
	SalaryMax       int    `json:"salary_max" valid:"-"`
	Description     string `json:"description" binding:"required" valid:"-"`
	Requirements    string `json:"requirements" valid:"-"`
	Duties          string `json:"duties" valid:"-"`
	Skills          string `json:"skills" valid:"-"`
	Sphere          *int   `json:"sphere" valid:"numeric~сфера деятельности должна содержать только код"`
	Employment      string `json:"employment" valid:"-"`
	ExperienceMonth int    `json:"experience_month" valid:"-"`
	Location        string `json:"location" valid:"stringlength(4|512)~длина адреса от 4 до 512 смиволов"`
	AreaSearch      string `json:"area_search" valid:"stringlength(4|128)~длина названия региона от 4 до 128 смиволов"`
	CareerLevel     string `json:"career_level" valid:"-"`
	EducationLevel  string `json:"education_level" valid:"-"`
	EmpEmail        string `json:"email" valid:"email"`
	EmpPhone        string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(4|18)~номер телефона от 4 до 18 цифр"`
}

type VacListRequest struct {
	Start  uint   `form:"start"`
	Limit  uint   `form:"limit" binding:"required"`
	CompId string `form:"comp_id,uuid"`
}

