package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ResumeHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/delivery/http"
	ResumeRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/repository"
	ResumeUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/delivery/rest"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/storage"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/common"

	"github.com/go-pg/pg/v9"
	logger "github.com/rowdyroad/go-simple-logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
}

type Logger struct {
	InfoLogger  *logger.Logger
	ErrorLogger *logger.Logger
}

type App struct {
	config   Config
	log      *Logger
	doneChan chan bool
	route    *gin.Engine
	db       *pg.DB
}

var log *Logger

func NewApp(config Config) *App {
	log := &Logger{
		InfoLogger:  logger.New(os.Stdout, "", logger.Lshortfile|logger.LstdFlags|logger.Llevel, logger.LevelInfo),
		ErrorLogger: logger.New(os.Stderr, "", logger.Lshortfile|logger.LstdFlags|logger.Llevel, logger.LevelWarning),
	}

	r := gin.New()
	r.Use(common.RequestLogger(log.InfoLogger))
	r.Use(common.ErrorLogger(log.ErrorLogger))
	r.Use(common.ErrorMiddleware())
	r.Use(common.Recovery(log.ErrorLogger))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	r.GET("/health", healthCheck())

	if config.DocPath != "" {
		r.Static("/doc/api", config.DocPath)
	} else {
		log.ErrorLogger.Warn("Document path is undefined")
	}

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Db.Host, config.Db.Port),
		User:     config.Db.User,
		Password: config.Db.Password,
		Database: config.Db.Name,
	})

	strg := storage.NewPostgresStorage(db)

	usecase := usecase.NewUsecase(log.InfoLogger, log.ErrorLogger, strg)

	rest.NewRest(r.Group("/v1"), *usecase)


	resumeRep := ResumeRepository.NewPgRepository(db)
	resume := ResumeUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, resumeRep)
	ResumeHandler.NewRest(r.Group("/v1"), resume)

	vacancyRep := Repos


	app := App{
		config:   config,
		log:      log,
		route:    r,
		doneChan: make(chan bool, 1),
		db:       db,
	}

	return &app
}

func (a *App) Run() {
	a.route.GET("/readiness", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    a.config.Listen,
		Handler: a.route,
	}

	go func() {
		a.log.InfoLogger.Infof("Start listening on %s", a.config.Listen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.ErrorLogger.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-a.doneChan:
	}
	a.log.InfoLogger.Info("Shutdown Server (timeout of 1 seconds) ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		a.log.ErrorLogger.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	a.log.InfoLogger.Info("Server exiting")
}

func (a *App) Close() {
	a.db.Close()
	a.doneChan <- true
}

func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	}
}
