package auth

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/auth"
)

type AuthServer struct {
	//usecase user.UseCase
	auth.UnimplementedAuthServer
}



func NewAuthServer() *AuthServer {
	return &AuthServer{}
}

func (a *AuthServer) Login(ctx context.Context, usr *auth.UserLogin) (*auth.UserId, error) {
	//sessionId, csrfToken, err := a.usecase.Login(usr.Login, usr.Password)
	//if err != nil {
	//	return nil, status.Error(codes.InvalidArgument, err.Error())
	//}

	return &auth.UserId{
		UserId: 1,
	}, nil
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
	//err := a.usecase.Logout(sid.SessionId)
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	return &auth.Empty{}, nil
}
