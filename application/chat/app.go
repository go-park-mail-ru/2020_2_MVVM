package chat

import (
	"context"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat/api/delivery/rest"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat/api/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Listen        string `yaml:"listen"`
	GinModuleMode string `yaml:"ginModuleMode"`
	DocPath       string `yaml:"docPath"`
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
}

func NewApp(config interface{}) *App {
	cfg := config.(Config)

	infoLogger, err := logger.New("test", 1, os.Stdout)
	errorLogger, err := logger.New("test", 2, os.Stderr)
	if err != nil {
		fmt.Println("logger init error")
	}

	log := &Logger{
		//InfoLogger:  logger.New(os.Stdout, "", logger.Lshortfile|logger.LstdFlags|logger.Llevel, logger.LevelInfo),
		//ErrorLogger: logger.New(os.Stderr, "", logger.Lshortfile|logger.LstdFlags|logger.Llevel, logger.LevelWarning),
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}

	infoLogger.SetLogLevel(logger.DebugLevel)

	r := gin.New()
	r.Use(common.RequestLogger(log.InfoLogger))
	r.Use(common.ErrorLogger(log.ErrorLogger))
	r.Use(common.ErrorMiddleware())
	r.Use(common.Recovery(log.ErrorLogger))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Sign", "X-Passenger"}
	r.Use(cors.New(corsConfig))

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("Incorrect request"))
	})

	usecase := usecase.NewUsecase(log.InfoLogger, log.ErrorLogger)
	rest.NewRest(r.Group("/v1"), usecase)

	a := &App{
		config:   cfg,
		log:      log,
		doneChan: make(chan bool, 1),
		route:    r,
	}

	return a
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
		a.log.ErrorLogger.Fatal(fmt.Sprintf("Server Shutdown: %s", err.Error()))
	}

	<-ctx.Done()
	a.log.InfoLogger.Info("Server exiting")
}

func (a *App) Close() {
	a.doneChan <- true
}
