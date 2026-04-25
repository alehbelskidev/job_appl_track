package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

type handler struct {
	s *service
}

func newHandler(s *service) *handler {
	return &handler{s: s}
}

func (h *handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var dto *RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	tokens, err := h.s.register(dto)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var dto *LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	tokens, err := h.s.login(dto)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var dto *TokensDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	tokens, err := h.s.refreshToken(dto.RefreshToken)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
