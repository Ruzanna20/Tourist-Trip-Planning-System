package jobservice

import (
	"log/slog"
	"time"
	"travel-planning/services"
)

type HotelJob struct {
	seeder *services.DataSeeder
}

func NewHotelJob(seeder *services.DataSeeder) *HotelJob {
	return &HotelJob{
		seeder: seeder,
	}
}

func (job *HotelJob) RunJob() {
	start := time.Now()

	l := slog.With("job", "HotelJob")

	l.Info("Job started")

	if err := job.seeder.SeedHotels(); err != nil {
		l.Error("Job failed with critical error", "error", err, "duration", time.Since(start))
	} else {
		l.Info("Job completed successfully", "duration", time.Since(start))
	}
}
