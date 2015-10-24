package application

import (
	"github.com/seesawlabs/ivan-kirichenko-exercise/handler"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

func (a *app) initRoutes() {
	a.logger.Infoln("initializing routing and handlers...")
	defer a.logger.Infoln("initializing routing and handlers")

	// routes for tasks CRUD operations
	tasks := a.server.Group("/task")
	tasks.Use(handler.GetJwtAuthHandler(a.config.JwtSecret))

	tasks.Get("/:id", handler.GetGetTaskHandler(a.db))
	tasks.Post("", handler.GetCreateTaskHandler(a.db))
	tasks.Patch("/:id", handler.GetUpdateTaskHandler(a.db))
	tasks.Delete("/:id", handler.GetDeleteTaskHandler(a.db))

	// routes for auth

	conf := oauth2.Config{
		ClientID:     a.config.OAuthAppID,
		ClientSecret: a.config.OAuthSecret,
		RedirectURL:  a.config.OAuthRedirectURL,
		Scopes:       []string{},
		Endpoint:     facebook.Endpoint,
	}
	a.server.Get("/auth",
		handler.GetOAuthHandler(
			conf,
			a.config.SessionSecret,
			a.csrfStorage,
		),
	)
	a.server.Get("/auth_verify",
		handler.GetOAuthVerifyHandler(
			conf,
			a.config.JwtSecret,
			a.config.SessionSecret,
			a.csrfStorage,
		),
	)
}
