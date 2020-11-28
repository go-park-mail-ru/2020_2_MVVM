package auth

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/auth"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/user"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

//Дергает юзкейсы

//Какие модели использовать и как матчить?
//Что в хендлерах, что в сервисе?
//ПОправить сессии

type AuthServer struct {
	usecase user.UseCase
	SessionBuilder common.SessionBuilder
	auth.UnimplementedAuthServer
}


func NewAuthServer(usecase user.UseCase, sessionBuilder common.SessionBuilder) *AuthServer {
	return &AuthServer{usecase: usecase, SessionBuilder: sessionBuilder}
}

func (a *AuthServer) Login(ctx context.Context, userLogin *auth.UserLogin) (*auth.User, error) {
	userModel, err := a.usecase.Login(models.UserLogin{
		Email:    userLogin.Login,
		Password: userLogin.Password,
	})
	if err != nil {
		if errMsg := err.Error(); errMsg == common.AuthErr {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//TODO
	session := a.SessionBuilder.Build(ctx)
	if userModel.UserType == common.Candidate {
		cand, err := a.usecase.GetCandidateByID(userModel.ID.String())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		session.Set(common.CandID, cand.ID.String())
		session.Set(common.EmplID, nil)

	} else if userModel.UserType == common.Employer {
		empl, err := a.usecase.GetEmployerByID(userModel.ID.String())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		session.Set("empl_id", empl.ID.String())
		session.Set("cand_id", nil)
	} else {
		return nil, status.Error(codes.InvalidArgument, common.AuthErr)
	}

	session.Set("user_id", userModel.ID.String())
	err = session.Save()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return *userModel, nil
}

//func (a *AuthServer) Check(ctx context.Context, s *api.SessionId) (*api.UserId, error) {
//	uid, err := a.usecase.Check(s.SessionId)
//	if err != nil {
//		return nil, status.Error(codes.InvalidArgument, err.Error())
//	}
//
//	return &api.UserId{UserId: uint64(uid)}, nil
//}

func (a *AuthServer) Logout(ctx context.Context, sid *auth.UserId) (*auth.Empty, error) {
	//TODO
	session := a.SessionBuilder.Build(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &auth.Empty{}, nil
}
