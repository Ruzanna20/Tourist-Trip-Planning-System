package jobservice

import (
	"log/slog"
	"time"
	"travel-planning/services"
)

type AttractionJob struct {
	seeder *services.DataSeeder
}

func NewAttractionJob(seeder *services.DataSeeder) *AttractionJob {
	return &AttractionJob{
		seeder: seeder,
	}
}

func (job *AttractionJob) RunJob() {
	start := time.Now()
	l := slog.With("job", "AttractionJob")

	l.Info("Job started")

	if err := job.seeder.SeedAttractions(); err != nil {
		l.Error("Job failed with critical error", "error", err, "duration", time.Since(start))
	} else {
		l.Info("Job completed successfully", "duration", time.Since(start))
	}
}
