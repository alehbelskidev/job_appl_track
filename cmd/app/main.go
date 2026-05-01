package main

import (
	"context"
	"log"
	"net/http"

	"github.com/alehbelskidev/job_appl_track/internal/applications"
	"github.com/alehbelskidev/job_appl_track/internal/auth"
	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config := shared.NewConfig()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := repo.New(pool)
	appMod := applications.NewModule(q, config)
	authMod := auth.NewModule(q, config)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))
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
