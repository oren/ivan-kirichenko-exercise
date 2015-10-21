package application

import (
	"github.com/Sirupsen/logrus"
	"github.com/deoxxa/echo-logrus"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// Config defines application config
type Config struct {
	ListenAddress string `yaml:"listen"`
	DbFile        string `yaml:"db_file"`
	JwtSecret     string `yaml:"jwt_secret"`
}

type app struct {
	config *Config
	logger *logrus.Logger
	server *echo.Echo
	db     *gorm.DB
}

// NewApp instantiates and initializes new application
func NewApp(config *Config, logger *logrus.Logger) (*app, error) {
	a := &app{}
	a.server = echo.New()
	a.logger = logger

	a.initMiddleware()
	if err := a.initDb(); err != nil {
		return nil, err
	}

	return a, nil
}

// Run tries to start the application. Panics in case of error
func (a *app) Run() {
	a.server.Run(a.config.ListenAddress)
}

func (a *app) initMiddleware() {
	// Middleware
	a.server.Use(echologrus.NewWithNameAndLogger("web", a.logger))
	a.server.Use(mw.Recover())
	a.server.Use(getJwtAuthMiddleware(a.config.JwtSecret))
}

func (a *app) initDb() error {
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		return err
	}
	if err := db.DB().Ping(); err != nil {
		return err
	}
	a.db = db

	return nil
}
