package jobservice

import (
	"log/slog"
	"time"
	"travel-planning/services"
)

type RestaurantJob struct {
	seeder *services.DataSeeder
}

func NewRestaurantJob(seeder *services.DataSeeder) *RestaurantJob {
	return &RestaurantJob{
		seeder: seeder,
	}
}

func (job *RestaurantJob) RunJob() {
	start := time.Now()
	l := slog.With("job", "RestaurantJob")

	l.Info("Job started")

	if err := job.seeder.SeedRestaurants(); err != nil {
		l.Error("Job failed with critical error", "error", err, "duration", time.Since(start))
	} else {
		l.Info("Job completed successfully", "duration", time.Since(start))
	}
}
