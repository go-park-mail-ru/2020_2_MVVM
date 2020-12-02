package api

import (
	"context"
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	SessionBuilder "github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	CustomExperienceRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience/repository"
	CustomExperienceUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience/usecase"
	//EducationRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/education/repository"
	//EducationUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/education/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/authmicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/vacancyMicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/middlewares"
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

type Microservices struct {
	Auth common.MicroserviceConfig `yaml:"auth"`
	Vac  common.MicroserviceConfig `yaml:"vacancy"`
}

type Config struct {
	Listen  string           `yaml:"listen"`
	Db      *common.DBConfig `yaml:"db"`
	DocPath string           `yaml:"docPath"`
	Micro   Microservices    `yaml:"microservices"`
}

type App struct {
	config   Config
	log      common.Logger
	doneChan chan bool
	route    *gin.Engine
	db       *gorm.DB
}

func NewApp(config Config) *App {
	gin.Default()
	r := gin.New()

	infoLogger, err := logger.New("Info logger", 1, os.Stdout)
	errorLogger, err := logger.New("Error logger", 2, os.Stderr)

	log := common.Logger{
		Info:  infoLogger,
		Error: errorLogger,
	}
	infoLogger.SetLogLevel(logger.DebugLevel)
	if config.DocPath != "" {
		r.Static("/doc/api", config.DocPath)
	} else {
		log.Error.Warning("Document path is undefined")
	}

	credentials := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d", config.Db.User,
		config.Db.Password, config.Db.Name,
		config.Db.Host, config.Db.Port)
	db, err := gorm.Open(postgres.Open(credentials), &gorm.Config{})
	if err != nil {
		log.Error.Fatal("connection to postgres db failed...")
	}
	r.Use(middlewares.RequestLogger(log.Info))
	r.Use(middlewares.ErrorLogger(log.Error))
	r.Use(middlewares.ErrorMiddleware())
	r.Use(middlewares.Recovery(log.Error))
	r.Use(middlewares.Cors())
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	r.GET("/health", healthCheck())

	govalidator.SetFieldsRequiredByDefault(false)

	api := r.Group("/api/v1")

	authMicro, err := authmicro.NewAuthClient(config.Micro.Auth.Host, config.Micro.Auth.Port, log)
	if err != nil {
		log.Error.Fatal("connection to the auth microservice failed...")
	}
	vacMicro, err := vacancyMicro.NewVacClient(config.Micro.Vac.Host, config.Micro.Vac.Port, log)
	if err != nil {
		log.Error.Fatal("connection to the vacancy microservice failed...")
	}

	authCookieConfig := common.AuthCookieConfig{
		Key:    "session",
		Path:   "/",
		Domain: "localhost", // for postman
		//Domain:   "studhunt.ru",
		MaxAge: int((time.Hour * 12).Seconds()),
		//Secure:   true,
		Secure:   false, // for postman
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	sessionBuilder := SessionBuilder.NewSessionBuilder{}
	authMiddleware := middlewares.AuthRequired(authCookieConfig, authMicro)

	UserRep := UserRepository.NewPgRepository(db)
	userCase := UserUseCase.NewUserUseCase(log.Info, log.Error, UserRep)
	UserHandler.NewRest(api.Group("/users"), userCase, authMicro, authCookieConfig, &sessionBuilder, authMiddleware)

	vacancyRep := RepositoryVacancy.NewPgRepository(db)
	vacancy := VacancyUseCase.NewVacUseCase(log.Info, log.Error, vacancyRep)
	VacancyHandler.NewRest(api.Group("/vacancy"), &sessionBuilder, authMiddleware, vacMicro, authMicro)

	companyRep := RepositoryCompany.NewPgRepository(db)
	company := CompanyUseCase.NewCompUseCase(log.Info, log.Error, companyRep)
	CompanyHandler.NewRest(api.Group("/company"), company, authMiddleware)

	resumeRep := ResumeRepository.NewPgRepository(db)
	//educationRep := EducationRepository.NewPgRepository(db)
	customExperienceRep := CustomExperienceRepository.NewPgRepository(db)

	//education := EducationUsecase.NewUsecase(log.Info, log.Error, educationRep)
	customExperience := CustomExperienceUsecase.NewUsecase(log.Info, log.Error, customExperienceRep)
	resume := ResumeUsecase.NewUseCase(log.Info, log.Error, userCase, customExperience, resumeRep)

	ResumeHandler.NewRest(api.Group("/resume"), resume, customExperience, &sessionBuilder, authMiddleware)

	responseRep := RepositoryResponse.NewPgRepository(db, vacancyRep)
	response := ResponseUseCase.NewUsecase(log.Info, log.Error, resume, *vacancy, company, responseRep)
	ResponseHandler.NewRest(api.Group("/response"), response, &sessionBuilder, authMiddleware)

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
		a.log.Info.Infof("Start listening on %s", a.config.Listen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-a.doneChan:
	}
	a.log.Info.Info("Shutdown Server (timeout of 1 seconds) ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		mes := fmt.Sprint("Server Shutdown:", err)
		a.log.Error.Fatal(mes)
	}

	<-ctx.Done()
	a.log.Info.Info("Server exiting")
}

func (a *App) Close() {
	a.doneChan <- true
}

func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Ok")
	}
}
