package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pgRepository struct {
	db            *gorm.DB
}

func NewPgRepository(db *gorm.DB) chat.ChatRepository {
	return &pgRepository{db: db}
}


func (p *pgRepository) Create(response models.Response) (*models.Chat, error) {
	var chat models.Chat
	chat.ResponseID = response.ID
	var empl models.Employer
	err := p.db.Raw(`select user_id
				from main.employers
				join main.vacancy on vacancy.empl_id = employers.empl_id
				where vacancy.vac_id = ?`, response.VacancyID).Scan(&empl).Error
	if err != nil {
		return nil, err
	}
	chat.EmplID = empl.UserID

	var cand models.Candidate
	err = p.db.Raw(`select user_id
				from main.candidates
				join main.resume on resume.cand_id = candidates.cand_id
				where resume.resume_id = ?`, response.ResumeID).Scan(&cand).Error
	if err != nil {
		return nil, err
	}
	chat.CandID = cand.UserID

	err = p.db.Create(&chat).Error
	if err.Error() == "ERROR: duplicate key value violates unique constraint \"chat_unique\" (SQLSTATE 23505)" {
		return nil, nil
	}
	if err != nil {
		err = fmt.Errorf("error in inserting chat: %w", err)
		return nil, err
	}
	return &chat, nil
}


func (p *pgRepository) GetById(chatID uuid.UUID, start uint, limit uint) ([]models.Message, error) {
	var mes []models.Message

	err := p.db.Raw(`with updated as (
			UPDATE main.message SET is_read = true WHERE chat_id = ? RETURNING *
		)
		select u.* from updated as u 
		order by date_create desc offset ? limit ?;`, chatID, start, limit).
		Scan(&mes).Error

	if err != nil {
		err = fmt.Errorf("error in select messages array from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return mes, nil
}

func (p *pgRepository) CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error) {
	var chat models.Chat
	err := p.db.First(&chat).Error
	if err != nil {
		return nil, err
	}
	if chat.CandID != sender && chat.EmplID != sender {
		err := fmt.Errorf("not found chat with this user")
		return nil, err
	}

	err = p.db.Create(&mes).Error
	if err != nil {
		err = fmt.Errorf("error in inserting message: %w", err)
		return nil, err
	}
	return &mes, nil
}
