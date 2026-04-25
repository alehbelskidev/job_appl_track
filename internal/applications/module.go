package applications

import (
	"database/sql"

	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handler
	repo    *repository
	config  *shared.Config
}

func NewModule(db *sql.DB, config *shared.Config) *Module {
	r := newRepository(db)
	h := newHandler(r)

	return &Module{repo: r, handler: h, config: config}
}

func (m *Module) Mount(r chi.Router) {
	r.Post("/applications", m.handler.createApplication)
	r.Get("/applications", m.handler.getApplications)
	r.Patch("/applications/{id}", m.handler.updateApplication)
}
