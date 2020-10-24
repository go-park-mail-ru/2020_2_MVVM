package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	CustomCompanyRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_company/repository"
	CustomCompanyUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_company/usecase"
	CustomExperienceRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience/repository"
	CustomExperienceUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience/usecase"
	EducationRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/education/repository"
	EducationUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/education/usecase"
	ResumeHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/delivery/http"
	ResumeRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/repository"
	ResumeUsecase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume/usecase"
	UserHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/delivery/http"
	UserRepository "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/repository"
	UserUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/user/usecase"
	VacancyHandler "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/delivery/http"
	RepositoryVacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/repository"
	VacancyUseCase "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/usecase"
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

	// the jwt middleware
	//identityKey := "myid"

	//authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
	//	Realm:       "my super test zone",
	//	Key:         []byte("my super secret and long secret-secret key"),
	//	Timeout:     20 * time.Hour,
	//	MaxRefresh:  20 * time.Hour,
	//	IdentityKey: identityKey,
	//
	//	SendCookie:     true,
	//	SecureCookie:   false, //non HTTPS dev environments
	//	CookieHTTPOnly: true,  // JS can't modify
	//	//CookieDomain:     "localhost",
	//	CookieDomain: "95.163.212.36",
	//
	//	Authenticator: func(c *gin.Context) (interface{}, error) {
	//		// This function should verify the user credentials given the gin context
	//		//(i.e. password matches hashed password for a given user email, and any other authentication logic).
	//		var credentials struct {
	//			Nickname string `form:"nickname" json:"nickname" binding:"required"`
	//			Email    string `form:"email" json:"email" binding:"required"`
	//			Password string `form:"password" json:"password" binding:"required"`
	//		}
	//		if err := c.ShouldBindJSON(&credentials); err != nil {
	//			return "", errors.New("missing Username, Password, or Email") // make error constant
	//		}
	//
	//		// go to the database and fetch the user
	//		var user models.User
	//		err := db.Model(&user).
	//			Where("email = ?", credentials.Email).
	//			Where("nickname = ?", credentials.Nickname).
	//			Select()
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		// compare password with the hashed one
	//		err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(credentials.Password))
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		// user is OK
	//		return &models.JWTUserData{
	//			ID:       user.ID,
	//			Nickname: user.Nickname,
	//			Email:    user.Email,
	//		}, nil
	//	},
	//	PayloadFunc: func(data interface{}) jwt.MapClaims {
	//		if v, ok := data.(*models.JWTUserData); ok {
	//			return jwt.MapClaims{
	//				identityKey: v.ID,
	//				"nickname":  v.Nickname,
	//				"email":     v.Email,
	//			}
	//		}
	//		return jwt.MapClaims{}
	//	},
	//	IdentityHandler: func(c *gin.Context) interface{} {
	//		claims := jwt.ExtractClaims(c)
	//		uid, _ := uuid.Parse(claims[identityKey].(string))
	//		return &models.JWTUserData{
	//			ID:       uid,
	//			Nickname: claims["nickname"].(string),
	//			Email:    claims["email"].(string),
	//		}
	//	},
	//	Authorizator: func(data interface{}, c *gin.Context) bool {
	//		return true
	//	},
	//	Unauthorized: func(c *gin.Context, code int, message string) {
	//		c.JSON(code, gin.H{
	//			"code":    code,
	//			"message": message,
	//		})
	//	},
	//	// TokenLookup is a string in the form of "<source>:<name>" that is used
	//	// to extract token from the request.
	//	// Optional. Default value "header:Authorization".
	//	// Possible values:
	//	// - "header:<name>"
	//	// - "query:<name>"
	//	// - "cookie:<name>"
	//	// - "param:<name>"
	//	TokenLookup: "header: Authorization, query: token, cookie: jwt",
	//	// TokenLookup: "query:token",
	//	// TokenLookup: "cookie:token",
	//
	//	// TokenHeadName is a string in the header. Default value is "Bearer"
	//	TokenHeadName: "Bearer",
	//
	//	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	//	TimeFunc:       time.Now,
	//	CookieSameSite: http.SameSiteNoneMode,
	//})

	//cookie.NewStore()

	gin.Default()
	//store := cookie.NewStore([]byte("secret"))
	store, _ := redis.NewStore(10, "tcp", "localhost:7001", "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	//if err != nil {
	//	log.ErrorLogger.Fatal("JWT Error:" + err.Error())
	//}

	//errInit := authMiddleware.MiddlewareInit()

	//if errInit != nil {
	//	log.ErrorLogger.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	//}

	api := r.Group("/api/v1")

	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://95.163.212.36"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") ||
				strings.HasPrefix(origin, "https://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	UserRep := UserRepository.NewPgRepository(db)
	userCase := UserUseCase.NewUserUseCase(log.InfoLogger, log.ErrorLogger, UserRep)
	UserHandler.NewRest(api.Group("/users"), userCase, common.AuthRequired())

	resumeRep := ResumeRepository.NewPgRepository(db)
	educationRep := EducationRepository.NewPgRepository(db)
	customCompanyRep := CustomCompanyRepository.NewPgRepository(db)
	customExperienceRep := CustomExperienceRepository.NewPgRepository(db)

	resume := ResumeUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, resumeRep)
	education := EducationUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, educationRep)
	customCompany := CustomCompanyUsecase.NewUseCase(log.InfoLogger, log.ErrorLogger, customCompanyRep)
	customExperience := CustomExperienceUsecase.NewUsecase(log.InfoLogger, log.ErrorLogger, customExperienceRep, customCompanyRep)

	ResumeHandler.NewRest(api.Group("/resume"), resume, education, customCompany, customExperience, common.AuthRequired())

	vacancyRep := RepositoryVacancy.NewPgRepository(db)
	vacancy := VacancyUseCase.NewVacUseCase(log.InfoLogger, log.ErrorLogger, vacancyRep)
	VacancyHandler.NewRest(api, vacancy)

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
