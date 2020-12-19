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
	db *gorm.DB
}

func NewPgRepository(db *gorm.DB) chat.ChatRepository {
	return &pgRepository{db: db}
}

func (p *pgRepository) MessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.MessageBrief, error) {
	var messages []models.MessageBrief
	query := p.db.Table("main.message").
		Where("chat_id = ?", chatID).
		Select("sender", "message", "is_read", "date_create").
		Order("date_create")
	if from != nil {
		query = query.Where("date_create > ?", from)
	}
	if to != nil {
		query = query.Where("date_create < ?", to)
	}
	if offset != nil {
		query = query.Offset(int(*offset))
	}
	if limit != nil {
		query = query.Limit(int(*limit))
	}

	err := query.Scan(&messages).Error

	return &messages, err
}

func (p *pgRepository) MarkMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error {
	var fromCond, toCond, offsetTemplate, limitTemplate string
	params := make([]interface{}, 0)
	params = append(params, chatID)

	if from != nil {
		fromCond = "and date_create > ?"
		params = append(params, *from)
	}
	if to != nil {
		toCond = "and date_create < ?"
		params = append(params, *to)
	}
	if offset != nil {
		offsetTemplate = fmt.Sprintf("offset %d", *offset)
	}
	if limit != nil {
		limitTemplate = fmt.Sprintf("limit %d", *limit)
	}

	template := fmt.Sprintf(`with matched as (
			select message_id
			from main.message 
			where chat_id = ? %s %s
			order by date_create %s %s
		)
		update main.message as m
		set is_read = true
		from matched
		where m.message_id = matched.message_id and m.sender <> ?;`, fromCond, toCond, offsetTemplate, limitTemplate)

	params = append(params, utype)
	return p.db.Exec(template, params...).Error
}

func (p *pgRepository) TechnicalMessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.TechMessageBrief, error) {
	// fetch tech messages
	var techMessages []models.TechMessageBrief

	query := p.db.Table("main.tech_message as m").
		Select("m.date_create", "m.response_id",
			"r.initial as response_initial", "m.response_status as response_status",
			"v.title as vacancy_title", "v.vac_id as vacancy_id",
			"r2.resume_id", "r2.title as resume_title",
			"o.comp_id as company_id", "o.name as company_name").
		Joins("inner join main.response r on r.response_id = m.response_id").
		Joins("inner join main.vacancy v on v.vac_id = r.vacancy_id").
		Joins("inner join main.resume r2 on r2.resume_id = r.resume_id").
		Joins("inner join main.official_companies o on o.comp_id = v.comp_id").
		Where("chat_id = ?", chatID).
		Order("m.date_create")
	if from != nil {
		query = query.Where("m.date_create > ?", from)
	}
	if to != nil {
		query = query.Where("m.date_create < ?", to)
	}
	if offset != nil {
		query = query.Offset(int(*offset))
	}
	if limit != nil {
		query = query.Limit(int(*limit))
	}

	err := query.Scan(&techMessages).Error
	return &techMessages, err
}

func (p *pgRepository) MarkTechnicalMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error {
	var shortcut string
	if utype == common.Candidate {
		shortcut = "cand"
	} else {
		shortcut = "empl"
	}

	var fromCond, toCond, offsetTemplate, limitTemplate string
	params := make([]interface{}, 0)
	params = append(params, chatID)

	if from != nil {
		fromCond = "and date_create > ?"
		params = append(params, *from)
	}
	if to != nil {
		toCond = "and date_create < ?"
		params = append(params, *to)
	}
	if offset != nil {
		offsetTemplate = fmt.Sprintf("offset %d", *offset)
	}
	if limit != nil {
		limitTemplate = fmt.Sprintf("limit %d", *limit)
	}

	sql := fmt.Sprintf(`with matched as (
			select message_id
			from main.tech_message 
			where chat_id = ? %s %s
			order by date_create %s %s
		)
		update main.tech_message as m
		set read_by_%s = true
		from matched
		where m.message_id = matched.message_id;`, fromCond, toCond, offsetTemplate, limitTemplate, shortcut)

	params = append(params, utype)
	return p.db.Exec(sql, params...).Error
}

func (p *pgRepository) CreateTechMesToUpdate(response models.Response) (*models.Chat, error) {
	var chat models.Chat
	var empl models.Employer
	err := p.db.Raw(`
			select e.user_id
			from main.response
         	join main.vacancy v on response.vacancy_id = v.vac_id
         	join main.employers e on e.empl_id = v.empl_id
         	join main.users u on u.user_id = e.user_id
			where response_id = ?;
			`, response.ID).Scan(&empl).Error
	if err != nil {
		return nil, err
	}
	chat.EmplID = empl.UserID

	var cand models.Candidate
	err = p.db.Raw(`
			select c.user_id
			from main.response
			join main.resume r using (resume_id)
			join main.candidates c on c.cand_id = r.cand_id
			join main.users u on u.user_id = c.user_id
			where response_id = ?;
			`, response.ID).Scan(&cand).Error
	if err != nil {
		return nil, err
	}
	chat.CandID = cand.UserID

	err = p.db.First(&chat).Error

	tech := models.TechMessage{
		ChatID:     chat.ChatID,
		ResponseID: response.ID,
		DateCreate: time.Now(),
		ResponseStatus: response.Status,
	}
	err = p.db.Create(&tech).Error
	if err != nil {
		err = fmt.Errorf("error in inserting message: %w", err)
		return nil, err
	}

	return &chat, nil
}

func (p *pgRepository) CreateChatAndTechMes(response models.Response) (*models.Chat, error) {
	var chat models.Chat
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
			err = p.db.Where("user_id_cand = ? and user_id_empl = ?", cand.UserID, empl.UserID).
				First(&chat).Error
			if err != nil {
				return nil, err
			}
		} else {
			err = fmt.Errorf("error in inserting chat: %w", err)
			return nil, err
		}
	}
	
	tech := models.TechMessage{
		ChatID:     chat.ChatID,
		ResponseID: response.ID,
		DateCreate: time.Now(),
		ResponseStatus: response.Status,

	}
	err = p.db.Create(&tech).Error
	if err != nil {
		err = fmt.Errorf("error in inserting message: %w", err)
		return nil, err
	}

	return &chat, nil
}

func (p *pgRepository) CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error) {
	chat := models.Chat{
		ChatID: mes.ChatID,
	}

	err := p.db.First(&chat).Error
	if err != nil {
		return nil, err
	}
	if chat.CandID != sender && chat.EmplID != sender {
		err := fmt.Errorf("there is no chat with this user")
		return nil, err
	}

	err = p.db.Create(&mes).Error
	if err != nil {
		err = fmt.Errorf("error in inserting message: %w", err)
		return nil, err
	}
	return &mes, nil
}

func (p *pgRepository) ListChats(userID uuid.UUID, utype string) ([]models.ChatSummary, error) {
	var toPrefix, fromPrefix, sender string
	if utype == common.Candidate {
		fromPrefix = "cand"
		toPrefix = "empl"
		sender = common.Employer
	} else {
		fromPrefix = "empl"
		toPrefix = "cand"
		sender = common.Candidate
	}

	// last dialog messages
	sql := fmt.Sprintf(`
		select distinct on (message.chat_id) message.chat_id, sender, message.is_read, message, date_create, users.name, 
											 users.surname, total_unread, users.path_to_avatar, 'message' as type
		from main.message
		inner join main.chat on message.chat_id = main.chat.chat_id
		inner join main.users on main.chat.user_id_%s = main.users.user_id
		inner join (select message.chat_id,
						   SUM(CASE WHEN message.is_read = False and message.sender = '%s' THEN 1 ELSE 0 END) AS total_unread
					from main.message
					group by message.chat_id) cte on cte.chat_id = main.message.chat_id
		where main.chat.user_id_%s = ?
		order by message.chat_id, date_create;
		`, toPrefix, sender, fromPrefix)

	type dialog struct {
		models.ChatSummary
		models.MessageBrief
	}

	var mes []dialog
	err := p.db.Raw(sql, userID).Scan(&mes).Error
	if err != nil {
		return nil, err
	}

	// last dialog tech messages
	sql = fmt.Sprintf(`
		select distinct on (tm.chat_id) tm.chat_id, tm.response_id, tm.date_create,
                                r.initial, tm.response_status as response_status,
                                v.title as vacancy_title, v.vac_id as vacancy_id,
                                r2.resume_id, r2.title as resume_title,
                                r.initial as response_initial, r.status as response_status,
                                o.comp_id as company_id, o.name as company_name,
								u.name, u.surname, u.path_to_avatar, total_unread, 'technical' as type
		from main.tech_message as tm
			inner join main.response r on r.response_id = tm.response_id
			inner join main.vacancy v on v.vac_id = r.vacancy_id
			inner join main.resume r2 on r2.resume_id = r.resume_id
			inner join main.official_companies o on o.comp_id = v.comp_id
			inner join main.chat c on tm.chat_id = c.chat_id
			inner join main.users u on c.user_id_%s = u.user_id
			inner join (select tech_message.chat_id,
							   SUM(CASE WHEN tech_message.read_by_%s = False THEN 1 ELSE 0 END) AS total_unread
						from main.tech_message
						group by tech_message.chat_id) cte on cte.chat_id = tm.chat_id
		where c.user_id_%s = ?
		order by tm.chat_id, date_create;`, toPrefix, fromPrefix, fromPrefix)

	type technical struct {
		models.ChatSummary
		models.TechMessageBrief
	}

	var techmes []technical
	err = p.db.Raw(sql, userID).Scan(&techmes).Error
	if err != nil {
		return nil, err
	}

	// pool them together
	result := make(map[uuid.UUID]*models.ChatSummary)
	for _, m := range mes {
		result[m.ChatID] = &models.ChatSummary{
			ChatID:      m.ChatID,
			TotalUnread: m.TotalUnread,
			Name:        m.Name,
			Surname:     m.Surname,
			Avatar:      m.Avatar,
			Type:        m.Type,
			Message:     m.MessageBrief,
		}
	}

	for _, m := range techmes {
		summary := models.ChatSummary{
			ChatID:      m.ChatID,
			TotalUnread: m.TotalUnread,
			Name:        m.Name,
			Surname:     m.Surname,
			Avatar:      m.Avatar,
			Type:        m.Type,
			Message:     m.TechMessageBrief,
		}

		val, ok := result[m.ChatID]
		if ok {
			created := val.Message.(models.MessageBrief).DateCreate
			totalUnread := result[m.ChatID].TotalUnread + summary.TotalUnread
			if created.Before(m.DateCreate) {
				result[m.ChatID] = &summary
			}
			result[m.ChatID].TotalUnread = totalUnread
		} else {
			result[m.ChatID] = &summary
		}
	}

	summaries := make([]models.ChatSummary, 0, len(result))
	for _, val := range result {
		summaries = append(summaries, *val)
	}
	return summaries, err
}


func (p *pgRepository) GetTotalUnreadMes(userID uuid.UUID, userType string) (*uint, error) {
	var typePrefix, sender string
	if userType == common.EmplID {
		sender = common.Candidate
		typePrefix = "empl"
	} else {
		sender = common.Employer
		typePrefix = "cand"
	}

	queryToMes := fmt.Sprintf(`
				select 
			SUM(CASE WHEN message.is_read = False and message.sender = '%s' THEN 1 ELSE 0 END) 
			AS user_mes
				from main.message
				join main.chat c on c.chat_id = message.chat_id
				where user_id_%s = ?`, sender, typePrefix)
	queryToTechMes := fmt.Sprintf(`
				select SUM(CASE WHEN tech_message.read_by_%s = False THEN 1 ELSE 0 END) AS tech_mes
				from main.tech_message
         		join main.chat c on c.chat_id = tech_message.chat_id
				where user_id_%s = ?`, typePrefix, typePrefix)

	var totalUnread struct {
		UserMes uint
		TechMes uint
	}

	err := p.db.Raw(queryToMes, userID).Scan(&totalUnread).Error
	if err != nil {
		return nil, err
	}
	err = p.db.Raw(queryToTechMes, userID).Scan(&totalUnread).Error
	if err != nil {
		return nil, err
	}

	total := totalUnread.UserMes + totalUnread.TechMes
	return &total, nil
}

func (p *pgRepository) OnlyUnreadMessagesForChat(chatID uuid.UUID, userType string) (*[]models.MessageBrief, error) {
	var messages []models.MessageBrief
	var sender string
	if userType == common.Candidate {
		sender = common.Employer
	} else {
		sender = common.Candidate
	}
	query := p.db.Table("main.message").
		Where("chat_id = ? and sender = ? and is_read = false", chatID, sender).
		Select("sender", "message", "is_read", "date_create").
		Order("date_create")
	err := query.Scan(&messages).Error

	return &messages, err
}

func (p *pgRepository) OnlyUnreadTechnicalMessagesForChat(chatID uuid.UUID, userType string) (*[]models.TechMessageBrief, error) {
	var techMessages []models.TechMessageBrief
	var typeToTable string
	if userType == common.Candidate {
		typeToTable = "cand"
	} else {
		typeToTable = "empl"
	}

	whereQuery := fmt.Sprintf("chat_id = ? and read_by_%s = False", typeToTable)

	query := p.db.Table("main.tech_message as m").
		Select("m.date_create", "m.response_id",
			"r.initial as response_initial", "m.response_status as response_status",
			"v.title as vacancy_title", "v.vac_id as vacancy_id",
			"r2.resume_id", "r2.title as resume_title",
			"o.comp_id as company_id", "o.name as company_name").
		Joins("inner join main.response r on r.response_id = m.response_id").
		Joins("inner join main.vacancy v on v.vac_id = r.vacancy_id").
		Joins("inner join main.resume r2 on r2.resume_id = r.resume_id").
		Joins("inner join main.official_companies o on o.comp_id = v.comp_id").
		Where(whereQuery, chatID).
		Order("m.date_create")

	err := query.Scan(&techMessages).Error
	return &techMessages, err
}
