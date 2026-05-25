package applications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	q repo.Querier
	s *service
}

func newHandler(q repo.Querier, s *service) *handler {
	return &handler{q: q, s: s}
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

	log.Printf("dto: %+v", dto)
	params := repo.CreateJobApplicationParams{
		Company:     dto.Company,
		Role:        dto.Role,
		DateApplied: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Status:      int32(applied),
		OwnerID:     *ownerID,
		Description: pgtype.Text{String: dto.Description, Valid: dto.Description != ""},
		Url:         pgtype.Text{String: dto.Url, Valid: dto.Url != ""},
		Notes:       pgtype.Text{String: dto.Notes, Valid: dto.Notes != ""},
	}
	log.Printf("params: %+v", params)
	app, err := h.q.CreateJobApplication(r.Context(), params)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response := createApplicationResponseDto{Data: app}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (h *handler) importJobApplications(w http.ResponseWriter, r *http.Request) {
	ownerID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var dtos []importJobApplicationsDto
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	params := make([]repo.ImportJobApplicationsParams, len(dtos))
	for i, d := range dtos {
		params[i] = repo.ImportJobApplicationsParams{
			Company:     d.Company,
			Role:        d.Role,
			DateApplied: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			Status:      int32(d.Status),
			OwnerID:     *ownerID,
			Description: pgtype.Text{String: d.Description, Valid: d.Description != ""},
			Url:         pgtype.Text{String: d.Url, Valid: d.Url != ""},
			Notes:       pgtype.Text{String: d.Notes, Valid: d.Notes != ""},
		}
	}

	batchResults := h.q.ImportJobApplications(r.Context(), params)

	var savedApps []repo.JobApplication
	var batchErr error

	batchResults.QueryRow(func(i int, app repo.JobApplication, err error) {
		if err != nil {
			batchErr = err
			return
		}
		savedApps = append(savedApps, app)
	})

	if batchErr != nil {
		log.Printf("Batch execution error: %v", batchErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(importApplicationsResultDto{Data: savedApps})
}

func (h *handler) createApplicationAI(w http.ResponseWriter, r *http.Request) {
	ownerID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("createApplicationAI :: Error parsing ownerID. [[%s]]", err)
		return
	}

	var dto createApplicationAIDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Cannot decode request body.", http.StatusBadRequest)
		log.Printf("createApplicationAI :: Error decoding body. [[%s]]", err)
		return
	}

	pageBody, err := h.s.parseHtmlPage(dto.Url)
	if err != nil {
		http.Error(w, "Cannot parse given url, try adding manually.", http.StatusInternalServerError)
		log.Printf("createApplicationAI :: Error parsing url. [[%s]]", err)
		return
	}

	prompt := fmt.Sprintf(
		"You're job is to parse given job application into specific format. JOB application CURL bod result: %s. Response in following format `{company: string, role: string, description: string}`. Pls use some html tags in description to structure content. page url: %s. JUST RETURN JSON NOT ```json{}``` RETURN PLAIN JSON STRING",
		*pageBody, dto.Url,
	)

	app, err := h.s.createApplicationFromPrompt(*ownerID, r.Context(), prompt, dto.Notes, dto.Url)
	if err != nil {
		http.Error(w, "Cannot create job application from AI prompt", http.StatusInternalServerError)
		log.Printf("createApplicationAI :: Error returned from service. [[%s]]", err)
		return
	}

	response := createApplicationResponseDto{Data: *app}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Cannot encode give application for client", http.StatusInternalServerError)
		log.Printf("createApplicationAI :: Create application error. [[%s]]", err)
		return
	}
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

func (h *handler) updateJobApplicationStatus(w http.ResponseWriter, r *http.Request) {
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
	payload := repo.UpdateJobApplicationStatusParams{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	payload.ID = id
	payload.OwnerID = *ownerID

	app, err := h.q.UpdateJobApplicationStatus(r.Context(), payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response := createApplicationResponseDto{Data: app}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (h *handler) deleteApplication(w http.ResponseWriter, r *http.Request) {
	_, err := h.getUserIDFromContext(r.Context())
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

	err = h.q.DeleteJobApplication(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response := deleteApplicationResponseDto{Data: true}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

// func (h *handler) updateApplication(w http.ResponseWriter, r *http.Request) {
// 	ownerID, err := h.getUserIDFromContext(r.Context())
// 	if err != nil {
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}
//
// 	idStr := chi.URLParam(r, "id")
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusBadRequest)
// 		log.Print(err)
// 		return
// 	}
//
// 	payload := repo.UpdateJobApplicationParams{}
// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		http.Error(w, "bad request", http.StatusBadRequest)
// 		log.Print(err)
// 		return
// 	}
//
// 	payload.ID = id
// 	payload.OwnerID = *ownerID
//
// 	app, err := h.q.UpdateJobApplication(r.Context(), payload)
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		log.Print(err)
// 		return
// 	}
//
// 	response := createApplicationResponseDto{Data: app}
// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		log.Print(err)
// 		return
// 	}
// }
