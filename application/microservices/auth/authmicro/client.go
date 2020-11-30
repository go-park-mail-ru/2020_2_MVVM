package authmicro

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
)

type AuthClient interface {
	Login(login string, password string) (common.Session, error)
	Check(sessionID string) (common.Session, error)
	Logout(sessionID string) error
}
