package application

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/deoxxa/echo-logrus"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pmylund/go-cache"
	"github.com/rs/cors"
	"github.com/seesawlabs/ivan-kirichenko-exercise/model"
)

// Config defines application config
type Config struct {
	ListenAddress    string `yaml:"listen"`
	DbFile           string `yaml:"db_file"`
	JwtSecret        string `yaml:"jwt_secret"`
	OAuthAppID       string `yaml:"oauth_appid"`
	OAuthSecret      string `yaml:"oauth_secret"`
	OAuthRedirectURL string `yaml:"oauth_redirect"`
	SessionSecret    string `yaml:"session_secret"`
}

// Runnable defines an interface that can run
type Runnable interface {
	Run()
}

// Migratable defines an interface that can run data migration
type Migratable interface {
	Migrate() error
}

type app struct {
	config       *Config
	logger       *logrus.Logger
	server       *echo.Echo
	db           *gorm.DB
	csrfStorage  *cache.Cache
	tokenStorage *cache.Cache
}

// NewApp instantiates and initializes new application
func NewApp(config *Config, logger *logrus.Logger) (Runnable, error) {
	a := &app{}
	a.config = config
	a.logger = logger

	a.server = echo.New()
	a.server.Use(echologrus.NewWithNameAndLogger("web", a.logger))
	a.server.Use(mw.Recover())
	a.server.Use(cors.Default().Handler)

	a.csrfStorage = cache.New(5*time.Minute, 30*time.Second)
	a.tokenStorage = cache.New(5*time.Minute, 30*time.Second)

	if err := a.initDb(); err != nil {
		return nil, err
	}
	a.initRoutes()

	return a, nil
}

// NewMigratable builds new instance of the app which can not run, but
// can migrate database
func NewMigratable(config *Config, logger *logrus.Logger) (Migratable, error) {
	a := &app{}
	a.config = config
	a.logger = logger
	err := a.initDb()
	return a, err
}

// Run tries to start the application. Panics in case of error
func (a *app) Run() {
	a.server.Run(a.config.ListenAddress)
}

// Migrate runs basic migrations for this simple application
func (a *app) Migrate() error {
	return a.db.AutoMigrate(&model.Task{}).Error
}

func (a *app) initDb() error {
	a.logger.Infoln("initializing database...")
	defer a.logger.Infoln("initializing database finished")

	db, err := gorm.Open("sqlite3", a.config.DbFile)
	if err != nil {
		return err
	}
	if err := db.DB().Ping(); err != nil {
		return err
	}
	a.db = &db

	return nil
}
