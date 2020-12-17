package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
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
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"chat_unique\" (SQLSTATE 23505)" {
			return nil, nil
		}
		err = fmt.Errorf("error in inserting chat: %w", err)
		return nil, err
	}

	mes := models.Message{
		MessageID:  uuid.UUID{},
		ChatID:     chat.ChatID,
		Sender:     "technical",
		Message:    "Костыль",
		IsRead:     false,
	}
	mes.DateCreate = time.Now()
	err = p.db.Create(&mes).Error
	if err != nil {
		err = fmt.Errorf("error in inserting message: %w", err)
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


func (p *pgRepository) ListChats(userID uuid.UUID, userType string) ([]models.BriefChat, error) {
	var listChats []models.BriefChat

	userForWhere := "user_id_"
	userForJoin := "user_id_"
	if userType == common.Candidate {
		userForWhere += "cand = ?"
		userForJoin += "empl"
	} else {
		userForWhere += "empl = ?"
		userForJoin += "cand"
	}

	query := fmt.Sprintf(`select DISTINCT ON(chat.chat_id) m.chat_id, message, 
								name, surname, path_to_avatar, sender, date_create
			from main.chat
			join main.message m on chat.chat_id = m.chat_id
			join main.users u on u.user_id = chat.%s
			where %s
			order by chat.chat_id, date_create asc`, userForJoin, userForWhere)

	err := p.db.Raw(query, userID).Scan(&listChats).Error
	if err != nil {
		err = fmt.Errorf("error in select messages: %w", err)
		return nil, err
	}
	return listChats, nil
}