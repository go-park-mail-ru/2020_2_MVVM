package response

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

type RespNotifications struct {
	UnreadResp          []models.ResponseWithTitle `json:"unread_resp"`
	UnreadRespCnt       uint                       `json:"unread_resp_cnt"`
	RecommendedVac      []models.Vacancy           `json:"recommended_vac"`
	RecommendedVacCnt   uint                       `json:"recommended_vac_cnt"`
	CountUnreadMessages *uint                      `json:"unread_messages"`
}

type ReqNotify struct {
	VacInLastNDays       *int        `json:"vac_in_last_n_days"` // notifications about recommended new vacancies, nil means from last 7 days max - month
	OnlyVacCnt           bool        `json:"only_new_vac_cnt"`   // if true -> get only count of recommended vacancies
	ListStart            uint        `json:"vac_list_start"`
	ListEnd              uint        `json:"vac_list_limit"`
	NewRespNotifications []uuid.UUID `json:"watched_responses"` // nil - all responses, for useless resp deleting put uuid in list
	OnlyRespCnt          bool        `json:"only_new_resp_cnt"` // if true -> get only count of responses notifications
	UnreadMessages       bool        `json:"unread_messages"`   // if true -> send unread messages
}
