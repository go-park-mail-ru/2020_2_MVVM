package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/auth"

	/*SessionBuilder "github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"*/
	/*UserHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/delivery/http"
	UserRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/repository"
	UserUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/usecase"*/
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
)

type Server struct {
	port string
	Auth *AuthServer
}

func NewServer(port string, db *gorm.DB, store sessions.Store) *Server {
	//sessions := session.NewSessionDatabase(rd, logger)
	//users := userRepository.NewUserDatabase(db, logger)
	//user := usecase.NewUser(sessions, users, logger)

	return &Server{
		port: port,
		Auth: NewAuthServer(),
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	gServer := grpc.NewServer()
	auth.RegisterAuthServer(gServer, s.Auth)

	err = gServer.Serve(listener)
	if err != nil {
		return nil
	}

	return nil
}


