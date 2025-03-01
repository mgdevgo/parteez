package server

func (server *Server) configureRoutes() {
	router := server.http

	router.Route("/health", server.handlers.health.Register)

	apiv1 := router.Group("/api/v1")
	apiv1.Route("/events", server.handlers.events.Register)
	apiv1.Route("/venues", server.handlers.venues.Register)
}
