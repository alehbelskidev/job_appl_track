package applications

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type handler struct {
	q repo.Querier
}

func newHandler(q repo.Querier) *handler {
	return &handler{q: q}
}

func (h *handler) getUserIDFromContext(ctx context.Context) (*uuid.UUID, error) {
	userID, ok := ctx.Value(shared.UserIDKey).(string)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return &ownerID, nil
}

func (h *handler) createApplication(w http.ResponseWriter, r *http.Request) {
	ownerID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var dto createJobApplicationDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	app, err := h.q.CreateJobApplication(r.Context(), repo.CreateJobApplicationParams{
		Company:     dto.Company,
		Role:        dto.Role,
		DateApplied: time.Now(),
		Status:      int32(applied),
		OwnerID:     *ownerID,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(app)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

type getJobApplicationsResponseDto struct {
	Data []repo.GetJobApplicationsRow `json:"data"`
}

func (h *handler) getApplications(w http.ResponseWriter, r *http.Request) {
	ownerID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	apps, err := h.q.GetJobApplications(r.Context(), *ownerID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	var response getJobApplicationsResponseDto
	if apps == nil {
		response.Data = make([]repo.GetJobApplicationsRow, 0)
	} else {
		response.Data = apps
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (h *handler) updateApplication(w http.ResponseWriter, r *http.Request) {
	ownerID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		log.Print(err)
		return
	}

	payload := repo.UpdateJobApplicationParams{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	payload.ID = id
	payload.OwnerID = *ownerID

	app, err := h.q.UpdateJobApplication(r.Context(), payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(app)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}
