package dataprocessor

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type CityFileRecord struct {
	Name        string
	CountryCode string
	Population  int
	Latitude    float64
	Longitude   float64
}

func FetchAllCitiesFromFile(filePath string, topN int) ([]CityFileRecord, error) {
	log.Printf("Reading and processing cities data from: %s",filePath)

	file,err := os.Open(filePath)
	if err != nil {
		return  nil,fmt.Errorf("failed to open cities file %s: %w",filePath,err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records,err := reader.ReadAll()
	if err != nil {
		return  nil,fmt.Errorf("failed to read CSV records: %w",err)
	}

	citiesByCountry := make(map[string][]CityFileRecord)

	const (
		colCity      = 0
		colLat       = 2
		colLng       = 3
		colCountry   = 5 
		colPopulation = 9
	)

	for i,record := range records {
		if i == 0 {
			continue
		}

		if len(record) < colPopulation +1 {
			log.Printf("Skipping row %d: Not enough fields.", i)
			continue
		}

		countryCode := record[colCountry]

		lat,_ := strconv.ParseFloat(record[colLat],64)
		lon,_ := strconv.ParseFloat(record[colLng],64)
		pop,_ := strconv.Atoi(record[colPopulation])

		if countryCode == "" || lat == 0 || lon == 0 {
			continue
		}

		city := CityFileRecord {
			Name: record[colCity],
			CountryCode: countryCode,
			Population: pop,
			Latitude: lat,
			Longitude: lon,
		}
		citiesByCountry[countryCode] = append(citiesByCountry[countryCode], city)
	}

	var topCities []CityFileRecord

	for _,cityList := range citiesByCountry {
		sort.Slice(cityList,func(i, j int) bool {
			return cityList[i].Population > cityList[j].Population
		})

		limit := topN
		if len(cityList) < limit {
			limit = len(cityList)
		}

		topCities = append(topCities, cityList[:limit]...)
	}

	log.Printf("Finished processing. Total %d cities selected for insertion.", len(topCities))
	return topCities, nil
}