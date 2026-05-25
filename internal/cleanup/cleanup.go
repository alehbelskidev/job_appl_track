package cleanup

import (
	"context"
	"log"

	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/robfig/cron/v3"
)

func RunCleanup(q repo.Querier) *cron.Cron {
	c := cron.New()

	log.Println("[[Registering cleanup]]")

	_, err := c.AddFunc("*/5 * * * *", func() {
		log.Println("[[Starting cleanup]]")

		err := q.CleanupGhostedApplicatrions(context.Background())
		if err != nil {
			log.Printf("[[Error cleanup]]: %s", err)
		}

		log.Println("[[Ending cleanup]]")
	})

	if err != nil {
		log.Printf("[[Error registering cleanup]]: %s", err)
	}

	c.Start()

	return c
}
