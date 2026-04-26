package auth

import (
	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handler
	service *service
	q       *repo.Queries
	config  *shared.Config
}

func NewModule(q *repo.Queries, config *shared.Config) *Module {
	s := newService(q, config)
	h := newHandler(s)

	return &Module{handler: h, service: s, config: config}
}

func (m *Module) Mount(r chi.Router) {
	r.Post("/register", m.handler.registerUser)
	r.Post("/login", m.handler.loginUser)
	r.Post("/refresh", m.handler.refreshToken)
}
