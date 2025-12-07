package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/gl1n0m3c/AKSP-KR/services/feemous/internal/store"
)

const adminID = "30f6afc4-b8eb-4e6b-b0bd-158fcef0dc28"

type Server struct {
	store *store.Store
}

func New(store *store.Store) *Server {
	return &Server{store: store}
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/health", s.handleHealth)

	r.Route("/api", func(api chi.Router) {
		// Public reads
		api.Get("/teams", s.handleListTeams)
		api.Get("/positions", s.handleListPositions)
		api.Get("/units", s.handleListUnits)
		api.Get("/users", s.handleListUsers)

		api.Post("/register", s.handleRegister)
		api.Post("/login", s.handleLogin)

		api.Group(func(pr chi.Router) {
			pr.Use(s.authMiddleware)
			pr.Get("/me", s.handleMe)
			pr.Get("/events", s.handleListEvents)
			pr.Post("/events", s.handleCreateEvent)
			pr.Put("/events/{id}", s.handleUpdateEvent)
			pr.Delete("/events/{id}", s.handleDeleteEvent)

			// Admin-only
			pr.Group(func(admin chi.Router) {
				admin.Use(s.adminMiddleware)
				admin.Post("/teams", s.handleCreateTeam)
				admin.Delete("/teams/{id}", s.handleDeleteTeam)

				admin.Post("/positions", s.handleCreatePosition)
				admin.Delete("/positions/{id}", s.handleDeletePosition)

				admin.Post("/units", s.handleCreateUnit)
				admin.Delete("/units/{id}", s.handleDeleteUnit)
			})
		})
	})

	// Minimal admin HTML
	r.Get("/admin", s.handleAdminPage)

	return r
}
