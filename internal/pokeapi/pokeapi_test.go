package pokeapi

import (
	"testing"
	"time"

	"github.com/khazakar/pokedexcli/internal/pokecache"
)

func TestGetLocationAreas(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Minute)
	baseURL := "https://pokeapi.co/api/v2/location-area"

	testCases := []struct {
		name           string
		url            string
		expectedResult bool
		expectError    bool
	}{
		{
			name:           "First page of location areas",
			url:            baseURL,
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Invalid URL",
			url:            "https://invalid-url",
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			locationAreas, err := GetLocationAreas(cache, tc.url)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if locationAreas == nil {
				t.Errorf("Expected non-nil location areas")
				return
			}

			if len(locationAreas.Results) == 0 {
				t.Errorf("Expected non-empty location areas")
			}
		})
	}
}

func TestGetLocationAreaDetails(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Minute)
	baseURL := "https://pokeapi.co/api/v2/location-area"

	testCases := []struct {
		name           string
		url            string
		areaName       string
		expectedResult bool
		expectError    bool
	}{
		{
			name:           "Explore specific location area",
			url:            baseURL,
			areaName:       "canalave-city-area",
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Invalid location area",
			url:            baseURL,
			areaName:       "non-existent-area",
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			locationAreaDetails, err := GetLocationAreaDetails(cache, tc.url, tc.areaName)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if locationAreaDetails == nil {
				t.Errorf("Expected non-nil location area details")
				return
			}

			if locationAreaDetails.Name == "" {
				t.Errorf("Expected non-empty location area name")
			}
		})
	}
}
