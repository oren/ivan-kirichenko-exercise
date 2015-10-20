package exercise

import (
	"github.com/Sirupsen/logrus"
	"github.com/deoxxa/echo-logrus"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// Config defines application config
type Config struct {
	ListenAddress string `yaml:"listen"`
	DbFile        string `yaml:"db_file"`
}

type app struct {
	config *Config
	logger *logrus.Logger
	server *echo.Echo
}

// NewApp instantiates and initializes new application
func NewApp(config *Config, logger *logrus.Logger) (*app, error) {
	a := &app{}
	a.server = echo.New()
	a.logger = logger

	a.initMiddleware()

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
}
