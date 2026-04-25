package auth

import (
	"database/sql"

	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handler
	service *service
	repo    *repository
	config  *shared.Config
}

func NewModule(db *sql.DB, config *shared.Config) *Module {
	r := newRepository(db)
	s := newService(r, config)
	h := newHandler(s)

	return &Module{repo: r, handler: h, service: s, config: config}
}

func (m *Module) Mount(r chi.Router) {
	r.Post("/register", m.handler.registerUser)
	r.Post("/login", m.handler.loginUser)
	r.Post("/refresh", m.handler.refreshToken)
}
