package client

import (
	"context"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/auth"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"os"
)

//Вызывается в хендлерах и в местном сервере!

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

func (a *AuthClient) Login(login string, password string) (*models.User, error) {
	usr := &auth.UserLogin{
		Login:    login,
		Password: password,
	}
	fmt.Println(usr.Login)
	//
	answer, err := a.client.Login(context.Background(), usr)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(answer.ID)
	if err != nil {
		return nil, err
	}
	user := models.User{
		ID:            userID,
		UserType:      answer.UserType,
		Name:          answer.Name,
		Surname:       answer.Surname,
		Email:         answer.Email,
		PasswordHash:  nil,
		Phone:         &answer.Phone,
		SocialNetwork: &answer.SocialNetwork,
	}
	return &user, nil
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

func (a *AuthClient) Logout(userID uuid.UUID) error {
	userIDStruct := &auth.UserId{UserId: userID.String()}
	_, err := a.client.Logout(context.Background(), userIDStruct)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthClient) Close() {
	if err := a.gConn.Close(); err != nil {
		a.logger.ErrorLogger.Error("error while closing grpc connection")
	}
}
