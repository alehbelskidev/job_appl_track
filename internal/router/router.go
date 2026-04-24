// Package router
package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alehbelskidev/job_appl_track/internal/repository"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	r    *chi.Mux
	repo *repository.Repo
}

func NewRouter(repo *repository.Repo) *Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Hello, world"); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	})

	router := &Router{r: r, repo: repo}

	r.Post("/api/applications", router.createApplication)
	r.Get("/api/applications", router.getApplications)

	return router
}

func (s *Router) Listen() error {
	return http.ListenAndServe(":3001", s.r)
}

func (s *Router) createApplication(w http.ResponseWriter, r *http.Request) {
	var app repository.JobApplication

	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
		}
	}()

	err := s.repo.CreateJobApplication(&app)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Router) getApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := s.repo.GetJobApplications()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(apps)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
	}
}
