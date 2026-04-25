package applications

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	repo *repository
}

func newHandler(repo *repository) *handler {
	return &handler{repo: repo}
}

func (h *handler) createApplication(w http.ResponseWriter, r *http.Request) {
	var app JobApplication

	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	err := h.repo.createJobApplication(&app)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (h *handler) getApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := h.repo.getJobApplications()
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

func (h *handler) updateApplication(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		log.Print(err)
		return
	}

	payload := &UpdateApplicationDto{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	if err := h.repo.updateApplication(id, payload); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
