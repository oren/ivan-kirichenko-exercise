package application

import "github.com/seesawlabs/ivan-kirichenko-exercise/handler"

func (a *app) initRoutes() {
	a.logger.Infoln("initializing routing and handlers...")
	defer a.logger.Infoln("initializing routing and handlers")

	// routes for tasks CRUD operations
	tasks := a.server.Group("/task")
	tasks.Use(getJwtAuthMiddleware(a.config.JwtSecret, a.tokenStorage))

	tasks.Get("/:id", handler.GetGetTaskHandler(a.db))
	tasks.Post("", handler.GetCreateTaskHandler(a.db))
	tasks.Patch("/:id", handler.GetUpdateTaskHandler(a.db))
	tasks.Delete("/:id", handler.GetDeleteTaskHandler(a.db))

	// routes for auth
	a.server.Get("/auth",
		handler.GetOauthHandler(
			a.config.OauthAppId,
			a.config.OauthSecret,
			a.config.OauthRedirectUrl,
			a.csrfStorage,
		),
	)
	a.server.Get("/auth_verify",
		handler.GetOauthVerifyHandler(
			a.config.OauthAppId,
			a.config.OauthSecret,
			a.config.OauthRedirectUrl,
			a.csrfStorage,
			a.tokenStorage,
		),
	)
}
