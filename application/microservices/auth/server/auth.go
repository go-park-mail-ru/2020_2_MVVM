package api

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/session"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/microservises/auth"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServer struct {
	usecase user.UseCase
	srepo   session.Repository
	auth.UnimplementedAuthServer
}

func convertSessionInfo(sessionID string, basic *common.BasicSession) *auth.SessionInfo {
	sinfo := auth.SessionInfo{SessionID: sessionID, UserID: basic.UserID.String()}
	if basic.EmplID != uuid.Nil {
		str := basic.EmplID.String()
		sinfo.EmplID = str
	}
	if basic.CandID != uuid.Nil {
		str := basic.CandID.String()
		sinfo.CandID = str
	}
	return &sinfo
}

func (a *authServer) Login(ctx context.Context, cred *auth.Credentials) (*auth.SessionInfo, error) {
	fmt.Print("Login")
	if cred == nil {
		return nil, errors.Errorf("Incorrect credentials format")
	}
	user, err := a.usecase.Login(models.UserLogin{
		Email:    cred.Login,
		Password: cred.Password,
	})

	if err != nil {
		if err.Error() == common.AuthErr {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// gather session information
	s := common.BasicSession{
		UserID: user.ID, EmplID: uuid.Nil, CandID: uuid.Nil,
	}
	if user.UserType == common.Candidate {
		cand, err := a.usecase.GetCandidateByID(user.ID.String())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		s.CandID = cand.ID
	} else if user.UserType == common.Employer {
		empl, err := a.usecase.GetEmployerByID(user.ID.String())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		s.EmplID = empl.ID
	} else {
		return nil, status.Error(codes.Internal, "Failed to determine user type")
	}

	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = a.srepo.Add(sessionID.String(), s)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fmt.Print("Login ok")

	return convertSessionInfo(sessionID.String(), &s), nil
}

func (a *authServer) Check(ctx context.Context, workload *auth.SessionID) (*auth.SessionInfo, error) {
	if workload == nil {
		return nil, errors.Errorf("Incorrect session id format")
	}
	sessionInfo, err := a.srepo.GetSession(workload.SessionID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return convertSessionInfo(workload.SessionID, sessionInfo), nil
}

func (a *authServer) Logout(ctx context.Context, workload *auth.SessionID) (*auth.Empty, error) {
	if workload == nil {
		return nil, errors.Errorf("Incorrect session id format")
	}
	err := a.srepo.Delete(workload.SessionID)
	return &auth.Empty{}, err
}

func NewAuthServer(usecase user.UseCase, srepo session.Repository) auth.AuthServer {
	return &authServer{usecase: usecase, srepo: srepo}
}
