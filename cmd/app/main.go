package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alehbelskidev/job_appl_track/internal/applications"
	"github.com/alehbelskidev/job_appl_track/internal/auth"
	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config := shared.NewConfig()
	db, err := sql.Open("pgx", config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	q := repo.New(db)
	appMod := applications.NewModule(q, config)
	authMod := auth.NewModule(q, config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		authMod.Mount(r)
	})
	r.Route("/api", func(r chi.Router) {
		r.Use(shared.AuthMiddleware(config.JwtSecret))
		appMod.Mount(r)
	})

	err = http.ListenAndServe(":3001", r)
	if err != nil {
		log.Fatal(err)
	}
}
