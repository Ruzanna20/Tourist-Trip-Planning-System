package services

import (
	"fmt"
	"log"
	"math"
	dataprocessor "travel-planning/data_processor"
	"travel-planning/models"
	"travel-planning/repository"
)

type DataSeeder struct {
	countryRepo *repository.CountryRepository
	cityRepo *repository.CityRepository
	attractionRepo *repository.AttractionRepository
}

func NewDataSeeder(countryRepo *repository.CountryRepository,cityRepo *repository.CityRepository,attractionRepo *repository.AttractionRepository) *DataSeeder {
	return &DataSeeder{
		countryRepo: countryRepo,
		cityRepo: cityRepo,
		attractionRepo: attractionRepo,
	}
}

func (ds *DataSeeder) SeedCountries() error {
	log.Println("Starting country seeding process...")

	countriesToSeed,err := FetchAllCountries()
	if err != nil {
		return fmt.Errorf("country API fetch failed: %w",err)
	}

	log.Printf("Fetched countries from API.Starting db insertion")

	for _,country := range countriesToSeed {
		existing,err := ds.countryRepo.GetByCode(country.Code)
		if err != nil {
			return fmt.Errorf("critical error checking if country %s (%s) exists: %w",country.Name,country.Code,err)
		}

		if existing != nil {
			continue
		}

		insertedID,err := ds.countryRepo.Insert(&country)
		if err != nil {
			log.Printf("ERROR.Failed to insert country %s (%s): %v",country.Name,country.Code,err)
			continue
		}

		log.Printf("SUCCESS: Seeded country %s with ID: %d", country.Name, insertedID)
	}

	log.Printf("Country seeding process finished.")
	return nil 
}


func (s *DataSeeder) SeedCities(filePath string) error {
    log.Println("Starting City Seeding process...")

    topCities, err := dataprocessor.FetchAllCitiesFromFile(filePath,5)
    if err != nil {
        return fmt.Errorf("data processing failed: %w", err)
    }

    countryIDMap, err := s.countryRepo.GetCountryCodeToIDMap()
    if err != nil {
        return fmt.Errorf("failed to get country ID map: %w", err)
    }

    for _, data := range topCities {
        countryID, exists := countryIDMap[data.CountryCode]
        if !exists {
            continue 
        }

        newCity := &models.City{
            CountryID: countryID,
            Name:      data.Name,
            Latitude:  data.Latitude,
            Longitude: data.Longitude,
            Description: fmt.Sprintf("%s is one of the top cities in %s.", data.Name, data.CountryCode),
        }

        cityID, err := s.cityRepo.Insert(newCity)
        if err != nil {
            log.Printf("Failed to insert city(%v) %s (%s): %v",cityID, newCity.Name, data.CountryCode, err)
            continue
        }
    }

    log.Printf("City seeding process finished.")
    return nil
}

const earthRadiusKm = 6371

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var degToRad = func(deg float64) float64 { return deg * (math.Pi / 180) }
	
	rLat1 := degToRad(lat1)
	rLon1 := degToRad(lon1)
	rLat2 := degToRad(lat2)
	rLon2 := degToRad(lon2)

	dLat := rLat2 - rLat1
	dLon := rLon2 - rLon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + 
	     math.Cos(rLat1)*math.Cos(rLat2)*
	     math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func (s *DataSeeder) SeedAttractions(filePath string) error {
    log.Println("Starting Attraction Seeding...")

    attractionData, err := dataprocessor.FetchAttractionFromEuropeanTour(filePath) 
    if err != nil {
        return fmt.Errorf("attraction data processing failed: %w", err)
    }

	cityLocations, err := s.cityRepo.GetAllCityLocations()
    if err != nil {
        return fmt.Errorf("failed to get city locations: %w", err)
    }
    if len(cityLocations) == 0 {
        return fmt.Errorf("no cities found in database to link attractions to")
    }

    for _, data := range attractionData {
        var closestCityID int
        minDistance := float64(math.MaxFloat64)
        
        for _, city := range cityLocations {
            distance := haversine(data.Latitude, data.Longitude, city.Latitude, city.Longitude)
            if distance < minDistance && distance < 200 { 
                minDistance = distance
                closestCityID = city.ID
            }
        }
        
        if closestCityID == 0 {
            log.Printf("Skipping attraction '%s': No close city found in DB (Lat:%.4f, Lon:%.4f).", data.Name, data.Latitude, data.Longitude)
            continue
        }

        newAttraction := &models.Attraction{
            CityID:      closestCityID, 
            Name:        data.Name,
            Category:    data.Category,
            Latitude:    data.Latitude,
            Longitude:   data.Longitude,
            Rating:      data.Rating, 
            EntryFee:    data.EntryFee,
            Currency:    "USD", 
            OpeningHours: "",
            Description: data.Description,
            ImageURL:    "", 
            Website:     "",
        }

        _, err := s.attractionRepo.Insert(newAttraction)
        if err != nil {
            log.Printf("Failed to insert attraction %s: %v", newAttraction.Name, err)
            continue
        }
    }

    log.Printf("Attraction seeding process finished.")
    return nil
}