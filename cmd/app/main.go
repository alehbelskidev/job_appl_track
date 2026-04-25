package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alehbelskidev/job_appl_track/internal/applications"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	config := shared.NewConfig()
	appMod := applications.NewModule(db, config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		appMod.Mount(r)
	})

	err = http.ListenAndServe(":3001", r)
	if err != nil {
		log.Fatal(err)
	}
}
