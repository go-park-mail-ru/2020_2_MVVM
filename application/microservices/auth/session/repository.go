package session

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
)

type Repository interface {
	Add(sessionID string, session common.BasicSession) error
	GetSession(sessionID string) (session *common.BasicSession, err error)
	Delete(sessionID string) error
}
