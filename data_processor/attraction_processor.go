package dataprocessor

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type AttractionFileRecord struct {
	CityName    string
	Name        string
	Category    string
	Latitude    float64
	Longitude   float64
	Rating      float64
	Description string
	EntryFee    float64
}

func FetchAttractionFromEuropeanTour(filePath string) ([]AttractionFileRecord, error) {
	log.Printf("Reading and processing Attractions data from: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open attractions file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	const (
		colName        = 0
		colCategory    = 3
		colLat         = 4
		colLon         = 5
		colDescription = 15
		colRating      = 16
		colEntryFee    = 17
	)

	var allAttractions []AttractionFileRecord

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < colDescription+1 {
			continue
		}
		name := strings.TrimSpace(record[colName])
		category := strings.TrimSpace(record[colCategory])
		description := strings.TrimSpace(record[colDescription])
		lat, _ := strconv.ParseFloat(record[colLat], 64)
		lon, _ := strconv.ParseFloat(record[colLon], 64)

		rating := 0.0
		if len(record) > colRating {
			if r, err := strconv.ParseFloat(record[colRating], 64); err == nil && r > 0 {
				rating = r
			}
		}

		entryFee := 0.0
		if len(record) > colEntryFee {
			if f, err := strconv.ParseFloat(record[colEntryFee], 64); err == nil {
				entryFee = f
			}
		}

		if name == "" || lat == 0 || lon == 0 {
			continue
		}

		newAttraction := AttractionFileRecord{
			Name:        name,
			Category:    category,
			Latitude:    lat,
			Longitude:   lon,
			Rating:      rating,
			Description: description,
			EntryFee:    entryFee,
		}

		allAttractions = append(allAttractions, newAttraction)
	}
	return allAttractions, nil
}
