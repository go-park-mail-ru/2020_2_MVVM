package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	vacancy2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/vacancy"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pgRepository struct {
	db            *gorm.DB
	vacRepository vacancy.RepositoryVacancy
}

func NewPgRepository(db *gorm.DB, vacRepository vacancy.RepositoryVacancy) response.ResponseRepository {
	return &pgRepository{db: db, vacRepository: vacRepository}
}

func (p pgRepository) GetRecommendedVacCnt(candId uuid.UUID, startDate string) (uint, error) {
	var (
		cnt  uint = 0
		temp uint = 0
	)
	step := 2
	curSphere := 0
	preferredSphere, err := p.vacRepository.GetPreferredSpheres(candId)
	if err != nil {
		if err.Error() == common.NoRecommendation {
			return 0, err
		}
		return 0, fmt.Errorf("error in GetRecommended vacancy cnt: %w", err)
	}
	for curSphere < vacancy2.CountSpheres {
		arr := []int{preferredSphere[curSphere].SphereInd, preferredSphere[curSphere+1].SphereInd}
		err = p.db.Raw("select count(*) from main.vacancy where date_create >= ? and sphere in ?", startDate, arr).Scan(&temp).Error
		cnt += temp
		if err != nil {
			err = fmt.Errorf("error in getRecommended vacancies: %w", err)
			return 0, err
		}
		curSphere += step
	}
	return cnt, nil
}

func (p pgRepository) GetRecommendedVacancies(candId uuid.UUID, start int, limit int, startDate string) ([]models.Vacancy, error) {
	step := 2
	curSphere := 0
	preferredSphere, err := p.vacRepository.GetPreferredSpheres(candId)
	if err != nil {
		if err.Error() == common.NoRecommendation {
			return nil, err
		}
		return nil, fmt.Errorf("error in GetRecommended vacancies: %w", err)
	}
	var (
		vacList []models.Vacancy
		list    []models.Vacancy
	)
	for len(vacList) < limit && curSphere < vacancy2.CountSpheres {
		arr := []int{preferredSphere[curSphere].SphereInd, preferredSphere[curSphere+1].SphereInd}
		err := p.db.Table("main.vacancy").Where("date_create >= ? and sphere in ?", startDate, arr).
			Order("date_create desc").Limit(limit).Offset(start).Find(&list).Error
		vacList = append(vacList, list...)
		if err != nil {
			err = fmt.Errorf("error in GetRecommendation: %w", err)
			return nil, err
		}
		curSphere += step
	}
	end := limit
	if limit > len(vacList) {
		end = len(vacList)
	}
	return vacList[0:end], err
}

func (p pgRepository) GetResponsesCnt(userId uuid.UUID, userType string) (uint, error) {
	var (
		err error
		cnt uint
	)
	if userType == common.Candidate {
		err = p.db.Raw("select count(*) from main.resume "+
			"join main.response using(resume_id)"+
			"where resume.cand_id = ?", userId).Scan(&cnt).Error
	} else {
		err = p.db.Raw("select count(*) from main.vacancy "+
			"join main.response on vacancy_id=vac_id "+
			"where vacancy.empl_id = ?", userId).Scan(&cnt).Error
	}
	return cnt, err
}

func (p pgRepository) GetRespNotifications(respIds []uuid.UUID, entityId uuid.UUID, entityType int) ([]models.Response, error) {
	var (
		responses []models.Response
		err       error
	)
	if entityType == common.Resume {
		err = p.db.Table("main.response").
			Where("resume_id = ? and response_id IN ?", entityId, respIds).
			UpdateColumn("unread", false).
			Error
	} else {
		err = p.db.Table("main.response").
			Where("vacancy_id = ? and response_id IN ?", entityId, respIds).
			UpdateColumn("unread", false).
			Error
	}
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	if entityType == common.Vacancy {
		err = p.db.Where("vacancy_id = ? and unread = ?", entityId, true).
			Order("date_create desc").
			Find(&responses).Error
	} else {
		err = p.db.Where("resume_id = ? and unread = ?", entityId, true).
			Order("date_create desc").
			Find(&responses).Error
	}
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgRepository) Create(response models.Response) (*models.Response, error) {
	err := p.db.Create(&response).Error
	if err != nil {
		err = fmt.Errorf("error in inserting response: %w", err)
		return nil, err
	}
	return &response, nil
}

func (p *pgRepository) GetByID(responseID uuid.UUID) (*models.Response, error) {
	response := new(models.Response)
	err := p.db.First(&response, responseID).Error
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (p *pgRepository) UpdateStatus(response models.Response) (*models.Response, error) {
	err := p.db.Model(&response).Update("status", response.Status).Error
	if err != nil {
		err = fmt.Errorf("error in updating response with id %s, : %w", response.ID.String(), err)
		return nil, err
	}
	return &response, nil
}

func (p *pgRepository) GetResumeAllResponse(resumeID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := p.db.Order("date_create desc").
		Find(&responses, "resume_id = ?", resumeID).
		Error
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgRepository) GetVacancyAllResponse(vacancyID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := p.db.Order("date_create desc").
		Find(&responses, "vacancy_id = ?", vacancyID).
		Error
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgRepository) GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.Resume, error) {
	var resume []models.Resume
	err := p.db.Raw(`select main.resume.* 
			from main.resume 
			left join main.response on main.response.resume_id = main.resume.resume_id 
			where cand_id = ?
			group by main.resume.resume_id 
			having sum(case when vacancy_id = ? then 1 else 0 end) = 0`, candID, vacancyID).
		Scan(&resume).Error
	if err != nil {
		return nil, err
	}
	return resume, nil
}

func (p *pgRepository) GetAllVacancyWithoutResponse(emplID uuid.UUID, resumeID uuid.UUID) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	err := p.db.Raw(`select main.vacancy.*
			from main.vacancy
			left join main.response on main.response.vacancy_id = main.vacancy.vac_id
			where empl_id = ?
			group by main.vacancy.vac_id
			having sum(case when resume_id = ? then 1 else 0 end) = 0`, emplID, resumeID).Scan(&vacancies).Error
	if err != nil {
		return nil, err
	}
	return vacancies, nil
}
