package main

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/api"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/server"
	RepositoryVacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/repository"
	VacancyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/usecase"
	yconfig "github.com/rowdyroad/go-yaml-config"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"
)

type Config struct {
	Listen string          `yaml:"listen"`
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
	yconfig.LoadConfig(&config, "configs/vacancy.yaml", nil)

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

	vacancyRep := RepositoryVacancy.NewPgRepository(db)
	vacancy := VacancyUseCase.NewVacUseCase(log.Info, log.Error, vacancyRep)

	vacServer := server.NewVacServer(vacancy)

	gServer := grpc.NewServer()
	api.RegisterVacancyServer(gServer, vacServer)
	err = gServer.Serve(listener)
	if err != nil {
		log.Error.Fatalf("error in listening api server: %s", err)
	}
}
