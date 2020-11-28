package client

type IAuthClient interface {
	Login(login string, password string) (userID string, err error)
	//Check(sessionId string) (userId uint, err error)
	Logout(userID string) error
}
