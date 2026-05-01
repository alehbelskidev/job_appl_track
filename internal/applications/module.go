package applications

import (
	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handler
	q       repo.Querier
	config  *shared.Config
}

func NewModule(q repo.Querier, config *shared.Config) *Module {
	h := newHandler(q)

	return &Module{q: q, handler: h, config: config}
}

func (m *Module) Mount(r chi.Router) {
	r.Post("/applications", m.handler.createApplication)
	r.Get("/applications", m.handler.getApplications)
	r.Patch("/applications/{id}/status", m.handler.updateJobApplicationStatus)
}
