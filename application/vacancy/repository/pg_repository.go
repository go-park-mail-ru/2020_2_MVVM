package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type pgRepository struct {
	db *gorm.DB
}

func NewPgRepository(db *gorm.DB) vacancy.RepositoryVacancy {
	return &pgRepository{db: db}
}

func (p *pgRepository) CreateVacancy(vac models.Vacancy) (*models.Vacancy, error) {
	employer := new(models.Employer)
	company  := new(models.OfficialCompany)

	err := p.db.Take(employer,"empl_id = ?", vac.EmpID).Error
	if err != nil {
		err = fmt.Errorf("error in FK search for vacancy creation for user with id: %s : error: %w", vac.EmpID, err)
		return nil, err
	}
	if compId := employer.CompanyID; compId != uuid.Nil {
		//vac.ID = uuid.New()
		vac.DateCreate = time.Now().Format("2006-01-02")
		vac.CompID = compId
		company.ID = compId
	} else {
		return nil, errors.New("error: employer must have company for vacancy creation")
	}
	err = p.db.Create(&vac).Error
	if err != nil {
		err = fmt.Errorf("error in inserting vacancy with title: %s : error: %w", vac.Title, err)
		return nil, err
	}
	err = p.db.Table("main.official_companies").Where("comp_id", vac.CompID).UpdateColumn("count_vacancy", gorm.Expr("count_vacancy + ?", 1)).Error
	if err != nil {
		return nil, fmt.Errorf("error in update company with id:  %s vacCount  : error: %w", company.ID.String(), err)
	}
	return &vac, nil
}

func (p *pgRepository) GetVacancyById(id uuid.UUID) (*models.Vacancy, error) {
	vac :=  new(models.Vacancy)
	if id == uuid.Nil {
		return nil, nil
	}
	err := p.db.Take(vac, "vac_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("error in select vacancy with id: %s : error: %w", id.String(), err)
	}
	return vac, nil
}

func (p *pgRepository) UpdateVacancy(newVac models.Vacancy) (*models.Vacancy, error) {
	oldVac, err := p.GetVacancyById(newVac.ID)
	if err != nil {
		return nil, fmt.Errorf("error in select vacancy with id: %s : error: %w", newVac.ID.String(), err)
	}
	if oldVac.EmpID != newVac.EmpID {
		return nil, fmt.Errorf("this user can't update this vacancy")
	}
	if err := p.db.Model(&oldVac).Updates(newVac).Error; err != nil {
		return nil, fmt.Errorf("can't update vacancy with id:%s", newVac.ID)
	}
	return &newVac, nil
}

func (p *pgRepository) GetVacancyList(start uint, limit uint, id uuid.UUID, entityType int) ([]models.Vacancy, error) {
	var (
		vacList []models.Vacancy
		err     error
	)
	if limit <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}

	if entityType == vacancy.ByEmpId {
		err = p.db.Where("empl_id = ?", id).Limit(int(limit)).Offset(int(start)).Order("date_create").Find(&vacList).Error
	} else if entityType == vacancy.ByCompId {
		err = p.db.Where("comp_id = ?", id).Limit(int(limit)).Offset(int(start)).Order("date_create").Find(&vacList).Error
	} else {
		err = p.db.Limit(int(limit)).Offset(int(start)).Order("date_create").Find(&vacList).Error
	}
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	if len(vacList) == 0 {
		return nil, nil
	}
	return vacList, nil
}

func (p *pgRepository) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	var vacList []models.Vacancy

	err := p.db.Table("main.vacancy").Scopes(func(q *gorm.DB) *gorm.DB {
		if params.StartDate != "" {
			q = q.Where("date_create >= (?)", params.StartDate)
		}
		if len(params.Spheres) != 0 {
			q = q.Where("sphere IN (?)", params.Spheres)
		}
		if len(params.EducationLevel) != 0 {
			q = q.Where("education_level IN (?)", params.EducationLevel)
		}
		if len(params.CareerLevel) != 0 {
			q = q.Where("career_level IN (?)", params.CareerLevel)
		}
		if len(params.ExperienceMonth) != 0 {
			q = q.Where("experience_month IN (?)", params.ExperienceMonth)
		}
		if len(params.Employment) != 0 {
			q = q.Where("employment IN (?)", params.ExperienceMonth)
		}
		if len(params.AreaSearch) != 0 {
			q = q.Where("area_search IN (?)", params.AreaSearch)
		}
		if params.SalaryMin != 0 || params.SalaryMax != 0 {
			q = q.Where("salary_min >= ?", params.SalaryMin).
				Where("salary_max <= ?", params.SalaryMax)
		}
		if params.KeyWords != "" {
			q = q.Where("LOWER(title) LIKE ?", "%"+params.KeyWords+"%")
		}
		if params.OrderBy != "" {
			return q.Order(params.OrderBy)
		}
		return q
	}).Find(&vacList).Error

	if err != nil {
		err = fmt.Errorf("error in vacancies list selection with searchParams: %s", err)
		return nil, err
	}
	if len(vacList) == 0 {
		return nil, nil
	}
	return vacList, nil
}

func (p *pgRepository) AddRecommendation(userID uuid.UUID, sphere int) error {
	nameSphere := "sphere" + strconv.Itoa(sphere)
	query := fmt.Sprintf(`insert into main.recommendation (user_id, %s) values (?, 1) ON CONFLICT (user_id) DO UPDATE SET %s = main.recommendation.%s + 1`, nameSphere, nameSphere, nameSphere)

	rec := new(models.Recommendation)

	err := p.db.Raw(query, userID).
		Scan(&rec).Error

	if err != nil {
		return fmt.Errorf("error in add recommendation: %w", err)
	}
	return nil
}

func (p *pgRepository) GetRecommendation(userID uuid.UUID, start int, limit int) ([]models.Vacancy, error) {
	rec := new(models.Recommendation)
	err := p.db.Take(rec, "user_id = ?", userID).Error
	if err != nil {
		return nil, fmt.Errorf("error in get for user recommendation: %w", err)
	}
	spheres := []int{rec.Sphere0, rec.Sphere1, rec.Sphere2, rec.Sphere3, rec.Sphere4, rec.Sphere5,
		rec.Sphere6, rec.Sphere7, rec.Sphere8, rec.Sphere9, rec.Sphere10, rec.Sphere11, rec.Sphere12,
		rec.Sphere13, rec.Sphere14, rec.Sphere15, rec.Sphere16, rec.Sphere17, rec.Sphere18, rec.Sphere19,
		rec.Sphere20, rec.Sphere21, rec.Sphere22, rec.Sphere23, rec.Sphere24, rec.Sphere25, rec.Sphere26,
		rec.Sphere27, rec.Sphere28, rec.Sphere29,
	}
	recSphere := IndexesWithMax(spheres)

	type Pr struct{
		avg int
	}
	var pr Pr

	var preferredSalary *float64
	err = p.db.Raw(`select avg(main.resume.salary_max - main.resume.salary_min)
		from main.resume
		join main.candidates on resume.cand_id = candidates.cand_id
		where user_id = ? and salary_min>0 and salary_max>0`, userID).
		Scan(&pr).Error
	if err != nil {
		return nil, fmt.Errorf("error in add recommendation: %w", err)
	}

	/*найти вакансии с этими сферами*/
	var vacList []models.Vacancy

	err = p.db.Table("main.vacancy").Scopes(func(q *gorm.DB) *gorm.DB {
		if preferredSalary != nil {
			q = q.Where("salary_min >= salary_min - (?)", *preferredSalary/5).
				Where("salary_max <= salary_max + (?)", *preferredSalary/5).
				Where("salary_min = 0 and salary_max = 0")
		}
			return q
		}).Where("sphere IN (?)", recSphere).
		Limit(limit).
		Offset(start).
		Order("date_create").
		Find(&vacList).Error

	if err != nil {
		return nil, fmt.Errorf("error in add recommendation: %w", err)
	}
	return vacList, nil
}


func IndexesWithMax(array []int) []int {
	var index = 0
	var max = array[index]
	var indexes []int
	for i, value := range array {
		if max < value {
			indexes = nil
			max = value
			indexes = append(indexes, i)
		}
		if max == value {
			indexes = append(indexes, i)
		}
	}
	return indexes
}
