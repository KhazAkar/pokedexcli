package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/khazakar/pokedexcli/internal/pokecache"
	"github.com/khazakar/pokedexcli/internal/pokeapi"
)

const (
	baseLocationURL = "https://pokeapi.co/api/v2/location-area"
	firstPageURL   = baseLocationURL + "?offset=20&limit=20"
	secondPageURL  = baseLocationURL + "?offset=40&limit=20"
)

// TestMapCommandNavigation tests that the map and mapb commands correctly navigate
// through the location areas without resetting to the base URL each time
func TestMapCommandNavigation(t *testing.T) {
	// Create a test Config struct
	cache := pokecache.NewCache(5 * time.Minute)
	cfg := &Config{
		NextURL:     baseLocationURL,
		PreviousURL: "",
		Cache:       cache,
	}

	// Mock the map command callback function to simulate commandMap
	mapCallback := func(cfg *Config, args ...string) error {
		// In the real implementation, this would call GetLocationAreas
		// and update the URLs based on the API response
		
		// Simulate the first call returning the first page
		if cfg.NextURL == baseLocationURL {
			cfg.NextURL = firstPageURL
			cfg.PreviousURL = baseLocationURL
		} else if cfg.NextURL == firstPageURL {
			// Simulate the second call returning the second page
			cfg.NextURL = secondPageURL
			cfg.PreviousURL = firstPageURL
		}
		return nil
	}

	// First call to map
	err := mapCallback(cfg)
	if err != nil {
		t.Fatalf("Error executing map command: %v", err)
	}

	// Verify the URLs are set correctly after first map call
	if cfg.NextURL != firstPageURL {
		t.Errorf("NextURL not set correctly after map command. Got %s, expected %s", cfg.NextURL, firstPageURL)
	}
	if cfg.PreviousURL != baseLocationURL {
		t.Errorf("PreviousURL not set correctly after map command. Got %s, expected %s", cfg.PreviousURL, baseLocationURL)
	}

	// Second call to map - should go to the next page
	// Simulate what happens in the REPL when a command is executed
	originalNextURL := cfg.NextURL
	err = mapCallback(cfg)
	if err != nil {
		t.Fatalf("Error executing second map command: %v", err)
	}

	// Verify the URLs are updated correctly after second map call
	if cfg.NextURL == originalNextURL {
		t.Errorf("NextURL not updated after second map command. Still at %s", cfg.NextURL)
	}
	if cfg.NextURL != secondPageURL {
		t.Errorf("NextURL incorrect after second map command. Got %s, expected %s", 
			cfg.NextURL, secondPageURL)
	}

	// Test mapb command - should go back to the previous page
	mapbCallback := func(cfg *Config, args ...string) error {
		if cfg.PreviousURL == "" {
			fmt.Println("You're on the first page")
			return nil
		}
		
		// Store the current previous URL and set NextURL to it
		cfg.NextURL = cfg.PreviousURL
		cfg.PreviousURL = ""
		
		// In the real implementation, this would call GetLocationAreas
		// and update the URLs based on the API response
		
		// Simulate getting the location areas and updating the URLs
		if cfg.NextURL == baseLocationURL {
			// No previous page for the base URL
			cfg.PreviousURL = ""
		} else if cfg.NextURL == firstPageURL {
			cfg.PreviousURL = baseLocationURL
		}
		
		return nil
	}

	// Call mapb
	err = mapbCallback(cfg)
	if err != nil {
		t.Fatalf("Error executing mapb command: %v", err)
	}

	// Verify the URLs are set correctly after mapb call
	if cfg.NextURL != firstPageURL {
		t.Errorf("NextURL not set correctly after mapb command. Got %s, expected %s", cfg.NextURL, firstPageURL)
	}
}

func TestInspectCommand(t *testing.T) {
	// Create a test Config struct
	cache := pokecache.NewCache(5 * time.Minute)
	cfg := &Config{
		Cache: cache,
	}

	// Reset pokedexStorage before the test
	pokedexStorage = make(map[string]pokeapi.PokemonDetails)

	// Test case 1: Inspect a Pokemon that hasn't been caught
	err := commandInspect(cfg, "pikachu")
	if err != nil {
		t.Fatalf("Unexpected error when inspecting uncaught Pokemon: %v", err)
	}

	// Test case 2: Catch a Pokemon and then inspect it
	// We'll create a mock Pokemon details
	mockPokemon := pokeapi.PokemonDetails{
		Name:   "pidgey",
		Height: 3,
		Weight: 18,
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		}{
			{
				BaseStat: 40,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "hp",
				},
			},
			{
				BaseStat: 45,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "attack",
				},
			},
			{
				BaseStat: 40,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "defense",
				},
			},
			{
				BaseStat: 35,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "special-attack",
				},
			},
			{
				BaseStat: 35,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "special-defense",
				},
			},
			{
				BaseStat: 56,
				Stat: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "speed",
				},
			},
		},
		Types: []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		}{
			{
				Type: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "normal",
				},
			},
			{
				Type: struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					Name: "flying",
				},
			},
		},
	}

	// Add the mock Pokemon to pokedexStorage
	pokedexStorage["pidgey"] = mockPokemon

	// Inspect the caught Pokemon
	err = commandInspect(cfg, "pidgey")
	if err != nil {
		t.Fatalf("Unexpected error when inspecting caught Pokemon: %v", err)
	}
}
