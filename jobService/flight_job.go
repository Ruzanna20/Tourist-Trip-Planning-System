package jobservice

import (
	"log/slog"
	"time"
	"travel-planning/services"
)

type FlightJob struct {
	seeder *services.DataSeeder
}

func NewFlightJob(seeder *services.DataSeeder) *FlightJob {
	return &FlightJob{
		seeder: seeder,
	}
}

func (job *FlightJob) RunJob() {
	start := time.Now()

	l := slog.With("job", "FlightJob")

	l.Info("Job started")

	if err := job.seeder.SeedFlights(); err != nil {
		l.Error("Job failed with critical error", "error", err, "duration", time.Since(start))
	} else {
		l.Info("Job completed successfully", "duration", time.Since(start))
	}
}
