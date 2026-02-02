package jobservice

import (
	"log/slog"
	"time"
	"travel-planning/services"
)

type CountryJob struct {
	seeder *services.DataSeeder
}

func NewCountryJob(seeder *services.DataSeeder) *CountryJob {
	return &CountryJob{
		seeder: seeder,
	}
}

func (job *CountryJob) RunJob() {
	start := time.Now()

	l := slog.With("job", "CountryJob")

	l.Info("Job started")

	if err := job.seeder.SeedCountries(); err != nil {
		l.Error("Job failed with critical error", "error", err, "duration", time.Since(start))
	} else {
		l.Info("Job completed successfully", "duration", time.Since(start))
	}
}
