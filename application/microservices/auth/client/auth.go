package client

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/auth"
	"google.golang.org/grpc"
	"os"
)

type AuthClient struct {
	client auth.AuthClient
	gConn  *grpc.ClientConn
	logger *Logger
}

type Logger struct {
	InfoLogger  *logger.Logger
	ErrorLogger *logger.Logger
}

func NewAuthClient(host, port string) (*AuthClient, error) {
	gConn, err := grpc.Dial(
		host+port,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	infoLogger, err := logger.New("Info logger", 1, os.Stdout)
	errorLogger, err := logger.New("Error logger", 2, os.Stderr)

	log := &Logger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}
	infoLogger.SetLogLevel(logger.DebugLevel)

	return &AuthClient{client: auth.NewAuthClient(gConn), gConn: gConn, logger: log}, nil
}

func (a *AuthClient) Login(login string, password string) (userID string, err error) {
	usr := &auth.UserLogin{
		Login:    login,
		Password: password,
	}
	fmt.Println(usr.Login)
	//
	//session, err := a.client.Login(context.Background(), usr)
	//if err != nil {
	//	return "", "", err
	//}

	return "1", nil
}

//func (a *AuthClient) Check(sessionId string) (userId uint, err error) {
//	sid := &api.SessionId{SessionId: sessionId}
//
//	uid, err := a.client.Check(context.Background(), sid)
//	if err != nil {
//		return 0, err
//	}
//
//	return uint(uid.UserId), err
//}

func (a *AuthClient) Logout(sessionId string) error {
	//sid := &api.UserId{UserId: 1}

	//_, err := a.client.Logout(context.Background(), sid)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (a *AuthClient) Close() {
	if err := a.gConn.Close(); err != nil {
		a.logger.ErrorLogger.Error("error while closing grpc connection")
	}
}
