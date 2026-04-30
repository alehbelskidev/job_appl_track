package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type handler struct {
	s *service
}

func newHandler(s *service) *handler {
	return &handler{s: s}
}

type userResponseDto struct {
	Email  string    `json:"email"`
	Tokens TokensDTO `json:"tokens"`
}

func (h *handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var dto *RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	tokens, err := h.s.register(r.Context(), dto)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response := userResponseDto{
		Email:  dto.Email,
		Tokens: *tokens,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var dto *LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	tokens, err := h.s.login(r.Context(), dto)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	response := userResponseDto{
		Email:  dto.Email,
		Tokens: *tokens,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var dto *RefreshTokenDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Print(err)
		return
	}

	tokens, err := h.s.refreshToken(r.Context(), dto.RefreshToken)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
