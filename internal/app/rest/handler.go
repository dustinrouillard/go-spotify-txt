package rest

import (
	"github.com/thetetrabyte/go-spotify-txt/internal/pkg/routes/authorization"
	"github.com/thetetrabyte/go-spotify-txt/internal/pkg/routes/misc"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// Initialize will define the routes and return a chi mux
func Initialize() *chi.Mux {
	// Create new chi Router
	Router := chi.NewRouter()

	// Initialize CORS
	CORS := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:9090", "http://localhost:9091", "http://127.0.0.1:9090", "http://127.0.0.1:9091"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Make the router use the cors handler
	Router.Use(CORS.Handler)

	// Use the authorization route
	Router.Get("/authorize", authorization.Authorize)

	// Use the callback route
	Router.Get("/callback", authorization.Callback)

	// Route not found handler
	Router.NotFound(misc.NotFound)

	// Return the router
	return Router
}
