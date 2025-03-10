package application

func (app *Application) configureRoutes() {
	router := app.http

	router.Route("/health", app.handlers.health.Register)

	apiv1 := router.Group("/api/v1")
	apiv1.Route("/events", app.handlers.events.Register)
	apiv1.Route("/venues", app.handlers.venues.Register)
}
