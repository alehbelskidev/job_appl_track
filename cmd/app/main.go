package main

import (
	"log"

	"github.com/alehbelskidev/job_appl_track/internal/repository"
	"github.com/alehbelskidev/job_appl_track/internal/router"
)

func main() {
	repo, err := repository.NewRepo()
	if err != nil {
		log.Fatal(err)
	}
	err = repo.InitTables()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	internalRouter := router.NewRouter(repo)
	if err := internalRouter.Listen(); err != nil {
		log.Fatal(err)
	}
}
