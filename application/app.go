package api

import (
	"context"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	SessionBuilder "github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	CustomExperienceRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience/repository"
	CustomExperienceUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience/usecase"
	EducationRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/education/repository"
	EducationUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/education/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/client"
	CompanyHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/delivery/http"
	RepositoryCompany "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/repository"
	CompanyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/usecase"
	ResponseHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/response/delivery/http"
	RepositoryResponse "github.com/go-park-mail-ru/2020_2_MVVM.git/application/response/repository"
	ResponseUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/response/usecase"
	ResumeHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/delivery/http"
	ResumeRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/repository"
	ResumeUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/usecase"
	UserHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/delivery/http"
	UserRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/repository"
	UserUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/usecase"
	VacancyHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/delivery/http"
	RepositoryVacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/repository"
	VacancyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	Redis   string    `yaml:"redis_address"`
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
	db       *gorm.DB
}

func NewApp(config Config) *App {
	gin.Default()
	r := gin.New()

	infoLogger, err := logger.New("Info logger", 1, os.Stdout)
	errorLogger, err := logger.New("Error logger", 2, os.Stderr)

	log := &Logger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}
	infoLogger.SetLogLevel(logger.DebugLevel)
	if config.DocPath != "" {
		r.Static("/doc/api", config.DocPath)
	} else {
		log.ErrorLogger.Warning("Document path is undefined")
	}

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

	r.Use(common.RequestLogger(log.InfoLogger))
	r.Use(common.ErrorLogger(log.ErrorLogger))
	r.Use(common.ErrorMiddleware())
	r.Use(common.Recovery(log.ErrorLogger))
	r.Use(common.Cors())
	r.Use(common.Sessions(store))
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	r.GET("/health", healthCheck())

	govalidator.SetFieldsRequiredByDefault(false)

	api := r.Group("/api/v1")

	rpcAuth, err := client.NewAuthClient("http://127.0.0.1", config.Listen)
	if err != nil {
		log.ErrorLogger.Fatal("connection to client microservice auth failed...")
	}

	UserRep := UserRepository.NewPgRepository(db)
	userCase := UserUseCase.NewUserUseCase(log.InfoLogger, log.ErrorLogger, UserRep)

	sessionBuilder := SessionBuilder.NewSessionBuilder{}

	UserHandler.NewRest(api.Group("/users"), userCase, &sessionBuilder, common.AuthRequired(), rpcAuth)

	vacancyRep := RepositoryVacancy.NewPgRepository(db)
	vacancy := VacancyUseCase.NewVacUseCase(log.InfoLogger, log.ErrorLogger, vacancyRep)
	VacancyHandler.NewRest(api.Group("/vacancy"), vacancy, &sessionBuilder, common.AuthRequired())

	companyRep := RepositoryCompany.NewPgRepository(db)
	company := CompanyUseCase.NewCompUseCase(log.InfoLogger, log.ErrorLogger, companyRep)
	CompanyHandler.NewRest(api.Group("/company"), company, common.AuthRequired())

	resumeRep := ResumeRepository.NewPgRepository(db)
	educationRep := EducationRepository.NewPgRepository(db)
	customExperienceRep := CustomExperienceRepository.NewPgRepository(db)

	education := EducationUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, educationRep)
	customExperience := CustomExperienceUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, customExperienceRep)
	resume := ResumeUsecase.NewUseCase(log.InfoLogger, log.ErrorLogger, userCase, education, customExperience, resumeRep)

	ResumeHandler.NewRest(api.Group("/resume"), resume, education, customExperience, &sessionBuilder, common.AuthRequired())

	responseRep := RepositoryResponse.NewPgRepository(db)
	response := ResponseUseCase.NewUsecase(log.InfoLogger, log.ErrorLogger, resume, *vacancy, company, responseRep)
	ResponseHandler.NewRest(api.Group("/response"), response, &sessionBuilder, common.AuthRequired())

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
		mes := fmt.Sprint("Server Shutdown:", err)
		a.log.ErrorLogger.Fatal(mes)
	}

	<-ctx.Done()
	a.log.InfoLogger.Info("Server exiting")
}

func (a *App) Close() {
	a.doneChan <- true
}

func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	}
}
