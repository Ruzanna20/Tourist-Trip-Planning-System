package jobservice

import (
	"log/slog"
	"time"
	"travel-planning/services"
)

type CityJob struct {
	seeder *services.DataSeeder
}

func NewCityJob(seeder *services.DataSeeder) *CityJob {
	return &CityJob{
		seeder: seeder,
	}
}

func (job *CityJob) RunJob() {
	start := time.Now()

	l := slog.With("job", "CityJob")

	l.Info("Job started")

	if err := job.seeder.SeedCities(); err != nil {
		l.Error("Job failed with critical error", "error", err, "duration", time.Since(start))
	} else {
		l.Info("Job completed successfully", "duration", time.Since(start))
	}
}
