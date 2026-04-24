// Package router
package router

import (
	"net/http"

	"github.com/alehbelskidev/job_appl_track/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	r    *chi.Mux
	repo *repository.Repo
}

func NewRouter(repo *repository.Repo) *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	router := &Router{r: r, repo: repo}

	r.Route("/api", func(r chi.Router) {
		r.Post("/applications", router.createApplication)
		r.Get("/applications", router.getApplications)
		r.Patch("/applications/{id}", router.updateApplication)
	})

	return router
}

func (s *Router) Listen() error {
	return http.ListenAndServe(":3001", s.r)
}
