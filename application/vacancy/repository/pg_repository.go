package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	vacancy2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/vacancy"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"path"
	"sort"
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
	company := new(models.OfficialCompany)

	err := p.db.Take(employer, "empl_id = ?", vac.EmpID).Error
	if err != nil {
		err = fmt.Errorf("error in FK search for vacancy creation for user with id: %s : error: %w", vac.EmpID, err)
		return nil, err
	}
	if compId := employer.CompanyID; compId != uuid.Nil {
		fileDir, _ := os.Getwd()
		avatarName := path.Join(common.ImgDir, "company", compId.String())
		imgPath := path.Join(fileDir, avatarName)
		if _, err = os.Stat(imgPath); err == nil {
			vac.Avatar = common.DOMAIN + avatarName
		}
		vac.DateCreate = time.Now().Format(time.RFC3339)
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
	//err = p.db.Table("main.official_companies").Where("comp_id", vac.CompID).UpdateColumn("count_vacancy", gorm.Expr("count_vacancy + ?", 1)).Error
	//if err != nil {
	//	return nil, fmt.Errorf("error in update company with id:  %s vacCount  : error: %w", company.ID.String(), err)
	//}
	return &vac, nil
}

func (p *pgRepository) GetVacancyById(id uuid.UUID) (*models.Vacancy, error) {
	vac := new(models.Vacancy)
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
	var empId string
	if err := p.db.Raw("select empl_id from main.vacancy where vac_id = ?", newVac.ID).Scan(&empId).Error; err != nil {
		return nil, fmt.Errorf("error in select vacancy with id: %s : error: %w", newVac.ID.String(), err)
	}
	if empId != newVac.EmpID.String() {
		return nil, fmt.Errorf("this user can't update this vacancy")
	}
	if err := p.db.Model(&newVac).Updates(newVac).Error; err != nil {
		return nil, fmt.Errorf("can't update vacancy with id:%s", newVac.ID)
	}
	return p.GetVacancyById(newVac.ID)
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
		err = p.db.Where("empl_id = ?", id).Limit(int(limit)).Offset(int(start)).Order("date_create desc").Find(&vacList).Error
	} else if entityType == vacancy.ByCompId {
		err = p.db.Where("comp_id = ?", id).Limit(int(limit)).Offset(int(start)).Order("date_create desc").Find(&vacList).Error
	} else {
		err = p.db.Limit(int(limit)).Offset(int(start)).Order("date_create desc").Find(&vacList).Error
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
			q = q.Where("date(date_create) >= (?)", params.StartDate)
		}
		if len(params.Sphere) != 0 {
			q = q.Where("sphere IN (?)", params.Sphere)
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
			q = q.Where("employment IN (?)", params.Employment)
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
		if params.KeywordsGeo != "" {
			q = q.Where("LOWER(area_search) LIKE ?", "%"+params.KeywordsGeo+"%") //for main
		}
		if params.OrderBy != "" {
			return q.Order(params.OrderBy)
		}
		return q.Order("date_create desc")
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

func (p *pgRepository) GetPreferredSpheres(userID uuid.UUID) ([]vacancy2.Pair, error) {
	rec := new(models.Recommendation)
	err := p.db.Take(rec, "user_id = ?", userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(common.NoRecommendation)
		}
		err = fmt.Errorf("error in GetPreferredSpheres: %w", err)
		return nil, err
	}

	spheres := []vacancy2.Pair{{0, rec.Sphere0}, {1, rec.Sphere1},
		{2, rec.Sphere2}, {3, rec.Sphere3}, {4, rec.Sphere4},
		{5, rec.Sphere5}, {6, rec.Sphere6}, {7, rec.Sphere7},
		{8, rec.Sphere8}, {9, rec.Sphere9}, {10, rec.Sphere10},
		{11, rec.Sphere11}, {12, rec.Sphere12}, {13, rec.Sphere13},
		{14, rec.Sphere14}, {15, rec.Sphere15}, {16, rec.Sphere16},
		{17, rec.Sphere17}, {18, rec.Sphere18}, {19, rec.Sphere19},
		{20, rec.Sphere20}, {21, rec.Sphere21}, {22, rec.Sphere22},
		{23, rec.Sphere23}, {24, rec.Sphere24}, {25, rec.Sphere25},
		{26, rec.Sphere26}, {27, rec.Sphere27}, {28, rec.Sphere28},
		{29, rec.Sphere29}}

	sort.Slice(spheres, func(i, j int) bool {
		return spheres[i].Score >= spheres[j].Score
	})
	return spheres, nil
}

func (p *pgRepository) GetPreferredSalary(userID uuid.UUID) (*float64, error) {
	preferredSalary := new(float64)
	//sal := &preferredSalary
	err := p.db.Raw(`select avg(main.resume.salary_max - main.resume.salary_min) as avg
		from main.resume
		join main.candidates on resume.cand_id = candidates.cand_id
		where user_id = ? and salary_min>0 and salary_max>0`, userID).
		Scan(preferredSalary).Error
	if err != nil {
		return nil, nil
	}
	return preferredSalary, nil
}

func (p *pgRepository) GetRecommendation(start int, limit int, salary float64, spheres []int) ([]models.Vacancy, error) {
	var vacList []models.Vacancy
	err := p.db.Table("main.vacancy").
		//
		//Scopes(func(q *gorm.DB) *gorm.DB {
		//if salary != nil {
		//	q = q.Where("salary_min >= salary_min - ? and salary_max <= salary_max + ?", int(*salary/5), int(*salary/5)).
		//		Or("salary_min = 0 and salary_max = 0")
		//}
		//	return q
		//}).

		Where("sphere IN (?)", spheres).
		Limit(limit).
		Offset(start).
		Order("date_create desc").
		Find(&vacList).Error

	if err != nil {
		return nil, fmt.Errorf("error in add recommendation: %w", err)
	}
	return vacList, nil
}

func (p *pgRepository) GetVacancyTopSpheres(topSpheresCnt int32) ([]models.Sphere, *models.VacTopCnt, error) {
	var (
		topList   []models.Sphere
		vacInfo   models.VacTopCnt
		allVacCnt uint64
		newVacCnt uint64
		notEmpty  = true
		err       error
	)

	check := p.db.Raw("select exists(select 1 top from main.vacancy)").Row()
	_ = check.Scan(&notEmpty)
	if !notEmpty {
		return topList, &vacInfo, nil
	}
	if topSpheresCnt == -1 {
		err = p.db.Raw("select * from main.sphere order by sphere_cnt desc").Scan(&topList).Error
	} else {
		err = p.db.Raw("select * from main.sphere order by sphere_cnt desc limit ?", topSpheresCnt).Scan(&topList).Error
	}
	currentDate := time.Now().Format("2006-01-02")
	if err == nil {
		err = p.db.Raw("select sum(sphere_cnt) from main.sphere").Scan(&allVacCnt).Error
	}
	if err == nil {
		err = p.db.Raw("select count(*) from main.vacancy where date(date_create) = ?", currentDate).Scan(&newVacCnt).Error
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error in add recommendation: %w", err)
	}
	vacInfo.NewVacCnt = newVacCnt
	vacInfo.AllVacCnt = allVacCnt
	return topList, &vacInfo, nil
}

func (p *pgRepository) DeleteVacancy(vacId uuid.UUID, empId uuid.UUID) error {
	err := p.db.Table("main.vacancy").Delete(&models.Vacancy{ID: vacId}).Where("empl_id = ?", empId).Error
	if err != nil {
		return fmt.Errorf("error in delete vacancy with id: %s, err: %w", vacId, err)
	}
	return nil
}