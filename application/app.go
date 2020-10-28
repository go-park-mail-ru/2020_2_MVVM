package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	UserHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/delivery/http"
	UserRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/repository"
	UserUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/usecase"
	VacancyHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/delivery/http"
	RepositoryVacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/repository"
	VacancyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/usecase"
	CompanyHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/delivery/http"
	RepositoryCompany "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/repository"
	CompanyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company/usecase"
	"github.com/go-pg/pg/v9"
	logger "github.com/rowdyroad/go-simple-logger"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	db       *pg.DB
}

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

	// Only for requests WITHOUT credentials, the literal value "*" can be specified
	corsMiddleware := cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://127.0.0.1") ||
				strings.HasPrefix(origin, "https://localhost") ||
				strings.HasPrefix(origin, "http://studhunt") ||
				strings.HasPrefix(origin, "https://studhunt")
		},
		MaxAge: time.Hour,
	})

	r.Use(corsMiddleware)

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


	gin.Default()
	store, err := redis.NewStore(10, "tcp", config.Redis, "", []byte("secret"))
	if err != nil {
		log.ErrorLogger.Fatal("connection to redis db failed...")
	}

	store.Options(sessions.Options{
		//Domain:   "studhunt.ru",
		Domain:   "127.0.0.1",
		MaxAge:   int((12 * time.Hour).Seconds()),
		Secure:   true,
		HttpOnly: false,
		Path: "/",
		SameSite: http.SameSiteNoneMode,
	})
	api := r.Group("/api/v1")

	sessionsMiddleware := sessions.Sessions("127.0.0.1", store)

	r.Use(sessionsMiddleware)
	UserRep := UserRepository.NewPgRepository(db)
	userCase := UserUseCase.NewUserUseCase(log.InfoLogger, log.ErrorLogger, UserRep)
	UserHandler.NewRest(api.Group("/users"), userCase, common.AuthRequired())

	/*resumeRep := ResumeRepository.NewPgRepository(db)
	educationRep := EducationRepository.NewPgRepository(db)
	customCompanyRep := CustomCompanyRepository.NewPgRepository(db)
	customExperienceRep := CustomExperienceRepository.NewPgRepository(db)

	resume := ResumeUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, resumeRep)
	education := EducationUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, educationRep)
	customCompany := CustomCompanyUsecase.NewUseCase(log.InfoLogger, log.ErrorLogger, customCompanyRep)
	customExperience := CustomExperienceUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, customExperienceRep, customCompanyRep)

	ResumeHandler.NewRest(api.Group("/resume"), resume, education, customCompany, customExperience, common.AuthRequired())*/

	companyRep := RepositoryCompany.NewPgRepository(db)
	company := CompanyUseCase.NewCompUseCase(log.InfoLogger, log.ErrorLogger, companyRep)
	CompanyHandler.NewRest(api.Group("/company"), company)

	vacancyRep := RepositoryVacancy.NewPgRepository(db)
	vacancy := VacancyUseCase.NewVacUseCase(log.InfoLogger, log.ErrorLogger, vacancyRep)
	VacancyHandler.NewRest(api.Group("/vacancy"), vacancy)

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
