package authmicro

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/microservises/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"strconv"
)

type gRPCAuthClient struct {
	client auth.AuthClient
	gConn  *grpc.ClientConn
	logger common.Logger
	ctx    context.Context
}

func buildSession(answer *auth.SessionInfo) common.BasicSession {
	s := common.BasicSession{SessionID: answer.SessionID, UserID: uuid.Nil, EmplID: uuid.Nil, CandID: uuid.Nil}
	s.UserID, _ = uuid.Parse(answer.UserID)
	if answer.CandID != "" {
		s.CandID, _ = uuid.Parse(answer.CandID)
	}
	if answer.EmplID != "" {
		s.EmplID, _ = uuid.Parse(answer.EmplID)
	}
	return s
}

func NewAuthClient(host string, port int, logger common.Logger) (AuthClient, error) {
	gConn, err := grpc.Dial(
		host+":"+strconv.Itoa(port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &gRPCAuthClient{client: auth.NewAuthClient(gConn), gConn: gConn, logger: logger, ctx: context.Background()}, nil
}

func (a *gRPCAuthClient) Login(login string, password string) (common.Session, error) {
	usr := &auth.Credentials{
		Login:    login,
		Password: password,
	}
	answer, err := a.client.Login(a.ctx, usr)
	if err != nil {
		return nil, err
	}

	s := buildSession(answer)
	return &s, nil
}

func (a *gRPCAuthClient) Check(sessionID string) (common.Session, error) {
	sid := &auth.SessionID{SessionID: sessionID}

	answer, err := a.client.Check(a.ctx, sid)
	if err != nil {
		return nil, err
	}
	s := buildSession(answer)
	return &s, err
}

func (a *gRPCAuthClient) Logout(sessionID string) error {
	workload := &auth.SessionID{SessionID: sessionID}
	_, err := a.client.Logout(a.ctx, workload)
	if err != nil {
		return err
	}
	return nil
}

func (a *gRPCAuthClient) Close() {
	if err := a.gConn.Close(); err != nil {
		a.logger.Error.Error("error while closing grpc connection")
	}
}
