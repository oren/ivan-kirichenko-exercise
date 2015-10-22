package application

import "github.com/seesawlabs/ivan-kirichenko-exercise/handler"

func (a *app) initRoutes() {
	a.logger.Infoln("initializing routing and handlers...")
	defer a.logger.Infoln("initializing routing and handlers")

	// task CRUD
	a.server.Get("/task/:id", handler.GetGetTaskHandler(a.db))
	a.server.Post("/task", handler.GetCreateTaskHandler(a.db))
	a.server.Put("/task/:id", handler.GetUpdateTaskHandler(a.db))
	a.server.Delete("/task/:id", handler.GetDeleteTaskHandler(a.db))
}
