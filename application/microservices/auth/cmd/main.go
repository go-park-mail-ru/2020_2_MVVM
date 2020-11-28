package main

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	api "github.com/go-park-mail-ru/2020_2_MVVM.git/application"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/server"
	yconfig "github.com/rowdyroad/go-yaml-config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type dbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Listen  string    `yaml:"listen"`
	Db      *dbConfig `yaml:"db"`
	DocPath string    `yaml:"docPath"`
	Redis   string    `yaml:"redis_address"`
}

type Logger struct {
	InfoLogger  *logger.Logger
	ErrorLogger *logger.Logger
}

func main() {
	infoLogger, err := logger.New("Info logger", 1, os.Stdout)
	errorLogger, err := logger.New("Error logger", 2, os.Stderr)

	log := &Logger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}
	infoLogger.SetLogLevel(logger.DebugLevel)

	var config api.Config
	yconfig.LoadConfig(&config, "configs/config.yaml", nil)

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d", config.Db.User,
		config.Db.Password, config.Db.Name,
		config.Db.Host, config.Db.Port)), &gorm.Config{})
	if err != nil {
		log.ErrorLogger.Fatal("connection to postgres db failed...")
	}
	store, err := redis.NewStore(10, "tcp", config.Redis, "", []byte("secret"))
	if err != nil {
		log.ErrorLogger.Fatal("connection to redis db failed...")
	}
	store.Options(sessions.Options{
		Domain: "studhunt.ru",
		//Domain:   "localhost", // for postman
		MaxAge:   int((12 * time.Hour).Seconds()),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		//SameSite: http.SameSiteStrictMode, // prevent csrf attack
	})

	srv1 := auth.NewServer(":8081", db, store)
	if err = srv1.ListenAndServe(); err != nil {
		log.ErrorLogger.FatalF("error in listening auth server: %s", err)
	}
}
