package main

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	server "github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/server"
	sessionrepo "github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/session/repository"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/repository"
	UserUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/microservises/auth"
	"github.com/go-redis/redis/v8"
	yconfig "github.com/rowdyroad/go-yaml-config"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"
)

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

type Config struct {
	Listen string          `yaml:"listen"`
	Redis  RedisConfig     `yaml:"redis"`
	DB     common.DBConfig `yaml:"db"`
}

func main() {
	infoLogger, err := logger.New("Info logger", 1, os.Stdout)
	errorLogger, err := logger.New("Error logger", 2, os.Stderr)

	log := common.Logger{
		Info:  infoLogger,
		Error: errorLogger,
	}
	infoLogger.SetLogLevel(logger.DebugLevel)

	var config Config
	yconfig.LoadConfig(&config, "configs/auth.yaml", nil)

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d",
		config.DB.User, config.DB.Password, config.DB.Name,
		config.DB.Host, config.DB.Port)), &gorm.Config{})
	if err != nil {
		log.Error.Fatal("connection to postgres db failed...")
	}

	listener, err := net.Listen("tcp", config.Listen)
	if err != nil {
		log.Error.Fatal(err.Error())
	}

	pg := repository.NewPgRepository(db)
	ucase := UserUseCase.NewUserUseCase(log.Info, log.Error, pg)

	redis := sessionrepo.NewRedisSessionRepository(redis.NewClient(&redis.Options{
		Addr: config.Redis.Address,
		Password: config.Redis.Password,
		DB: config.Redis.DB,
	}), log)
	server := server.NewAuthServer(ucase, redis)

	gServer := grpc.NewServer()
	auth.RegisterAuthServer(gServer, server)
	err = gServer.Serve(listener)
	if err != nil {
		log.Error.Fatalf("error in listening api server: %s", err)
	}
}
