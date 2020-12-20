package models

import (
	"fmt"
	"strconv"
)

type ExperienceTemplate struct {
	NameJob  string
	Position string
	Begin    string
	Finish   string
	Duties   string
}

type ResumeTemplate struct {
	Avatar          string
	Name            string
	Surname         string
	AreaSearch      string
	Phone           string
	Email           string
	Sphere          string
	Experience      string
	Education       string
	Salary          string
	Description     string
	ExperienceItems []ExperienceTemplate
	Skills          string
}

var ExperienceMonth = map[int]string{
	0:  "Не работал",
	1:  "Меньше года",
	5:  "1-5 лет",
	10: "5-10 лет",
}

var EducationLevel = map[string]string{
	"middle":                "Среднее",
	"specialized_secondary": "Среднее специальное",
	"incomplete_higher":     "Неоконченное высшее",
	"higher":                "Высшее",
	"bachelor":              "Бакалавр",
	"master":                "Магистр",
	"phD":                   "Кандидат наук",
	"doctoral":              "Доктор наук",
}

var ListSpheres = []string{
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

func СonvertToTemplate(r Resume) ResumeTemplate {
	t := ResumeTemplate{
		Avatar:      r.Avatar,
		Name:        r.CandName,
		Surname:     r.CandSurname,
		Email:       r.Candidate.User.Email,
		Description: r.Description,
		Skills:      r.Skills,
	}
	if r.AreaSearch != nil {
		t.AreaSearch = *(r.AreaSearch)
	} else {
		t.AreaSearch = ""
	}
	if r.Candidate.User.Phone != nil {
		t.Phone = *(r.Candidate.User.Phone)
	} else {
		t.Phone = ""
	}
	if r.Sphere != nil {
		t.Sphere = ListSpheres[*r.Sphere]
	} else {
		t.Sphere = ""
	}
	if r.ExperienceMonth != nil {
		t.Experience = ExperienceMonth[*r.ExperienceMonth]
	} else {
		t.Experience = ""
	}
	if r.EducationLevel != nil {
		t.Education = EducationLevel[*(r.EducationLevel)]
	} else {
		t.Education = ""
	}
	if r.SalaryMin != nil {
		if r.SalaryMax != nil {
			t.Salary = fmt.Sprintf(`%v - %v`, *r.SalaryMin, *r.SalaryMax)
		} else {
			t.Salary = strconv.Itoa(*r.SalaryMin)
		}
	} else if r.SalaryMax != nil {
		t.Salary = strconv.Itoa(*r.SalaryMax)
	} else {
		t.Salary = ""
	}

	var exp []ExperienceTemplate
	for _, val := range r.ExperienceCustomComp {
		expItem := ExperienceTemplate{NameJob: val.NameJob,
			Begin: fmt.Sprintf(val.Begin.Format("2006-01-02"))}
		if val.Position != nil {
			expItem.Position = *(val.Position)
		} else {
			expItem.Position = ""
		}
		if val.Duties != nil {
			expItem.Duties = *(val.Duties)
		} else {
			expItem.Duties = ""
		}
		if val.Finish != nil {
			expItem.Finish = fmt.Sprintf(val.Finish.Format("2006-01-02"))
		} else {
			expItem.Finish = "по настоящее время"
		}
		exp = append(exp, expItem)
	}
	t.ExperienceItems = exp

	return t
}
