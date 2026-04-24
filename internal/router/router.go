// Package router
package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alehbelskidev/job_appl_track/internal/dto"
	"github.com/alehbelskidev/job_appl_track/internal/models"
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
	r.Patch("/api/applications/:id", router.updateApplication)

	return router
}

func (s *Router) Listen() error {
	return http.ListenAndServe(":3001", s.r)
}

func (s *Router) createApplication(w http.ResponseWriter, r *http.Request) {
	var app models.JobApplication

	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			log.Print(err)
			return
		}
	}()

	err := s.repo.CreateJobApplication(&app)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *Router) getApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := s.repo.GetJobApplications()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(apps)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *Router) updateApplication(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		log.Print(err)
		return
	}

	payload := &dto.UpdateApplicationDto{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	if err := s.repo.UpdateApplication(id, payload); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
