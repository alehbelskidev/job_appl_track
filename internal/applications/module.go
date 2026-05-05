package applications

import (
	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handler
	service *service
	q       repo.Querier
	config  *shared.Config
}

func NewModule(q repo.Querier, config *shared.Config) *Module {
	s := newService(config, q)
	h := newHandler(q, s)

	return &Module{q: q, handler: h, config: config, service: s}
}

func (m *Module) Mount(r chi.Router) {
	r.Post("/applications", m.handler.createApplication)
	r.Post("/applications/ai", m.handler.createApplicationAI)
	r.Get("/applications", m.handler.getApplications)
	r.Patch("/applications/{id}/status", m.handler.updateJobApplicationStatus)
	r.Delete("/applications/{id}", m.handler.deleteApplication)
}
