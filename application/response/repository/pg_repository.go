package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pgRepository struct {
	db *gorm.DB
}

func (p pgRepository) GetRespNotifications(respIds map[uuid.UUID]bool) ([]models.Response, error) {
	var (
		responses  []models.Response
		wasReadIds []uuid.UUID
	)
	for key, isActive := range respIds {
		if !isActive {
			wasReadIds = append(wasReadIds, key)
		}
	}
	err := p.db.Table("main.response").Where("response_id IN ?", wasReadIds).UpdateColumn("unread", false).Error
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	err = p.db.Where("unread = ?", true).Find(&responses).Error
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func NewPgRepository(db *gorm.DB) response.ResponseRepository {
	return &pgRepository{db: db}
}

func (p *pgRepository) Create(response models.Response) (*models.Response, error) {
	err := p.db.Create(&response).Error
	//_, err := p.db.Model(&response).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting response: %w", err)
		return nil, err
	}
	return &response, nil
}

func (p *pgRepository) GetByID(responseID uuid.UUID) (*models.Response, error) {
	response := new(models.Response)
	err := p.db.First(&response, responseID).Error
	//err := p.db.Model(&response).Where("response_id = ?", responseID).Select()
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
	err := p.db.Find(&responses, "resume_id = ?", resumeID).Error
	//err := p.db.Model(&responses).Where("resume_id = ?", resumeID).Select()
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgRepository) GetVacancyAllResponse(vacancyID uuid.UUID) ([]models.Response, error) {
	var responses []models.Response
	err := p.db.Find(&responses, "vacancy_id = ?", vacancyID).Error
	//err := p.db.Model(&responses).Where("vacancy_id = ?", vacancyID).Select()
	if err != nil {
		err = fmt.Errorf("error in get list responses: %w", err)
		return nil, err
	}
	return responses, nil
}

func (p *pgRepository) GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.Resume, error) {
	var resume []models.Resume
	//query := fmt.Sprintf(`select resume.*
	//		from resume
	//		left join response on response.resume_id = resume.resume_id
	//		where cand_id = '%s'
	//		group by resume.resume_id
	//		having sum(case when vacancy_id = '%s' then 1 else 0 end) = 0`, candID, vacancyID)
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
