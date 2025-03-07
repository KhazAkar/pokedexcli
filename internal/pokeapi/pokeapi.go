package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/KhazAkar/pokedexcli/internal/pokecache"
)

type LocationAreaMap struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaExplore struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func fetchAPIData[T any](cache *pokecache.Cache, url string) (*T, error) {
	// First, check if the data is in the cache
	if cachedData, exists := cache.Get(url); exists {
		var data T
		err := json.Unmarshal(cachedData, &data)
		if err == nil {
			return &data, nil
		}
	}

	// If not in cache, make the network request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the response
	var data T
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// Add the response to the cache
	cache.Add(url, body)

	return &data, nil
}

func GetLocationAreas(cache *pokecache.Cache, url string) (*LocationAreaMap, error) {
	return fetchAPIData[LocationAreaMap](cache, url)
}

func GetLocationAreaDetails(cache *pokecache.Cache, url string, areaName string) (*LocationAreaExplore, error) {
	return fetchAPIData[LocationAreaExplore](cache, fmt.Sprintf("%s/%s", url, areaName))
}
