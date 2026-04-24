package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alehbelskidev/job_appl_track/internal/dto"
	"github.com/alehbelskidev/job_appl_track/internal/models"
	"github.com/go-chi/chi/v5"
)

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
