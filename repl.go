package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/KhazAkar/pokedexcli/internal/pokeapi"
	"github.com/KhazAkar/pokedexcli/internal/pokecache"
)

type Config struct {
	NextURL     string
	PreviousURL string
	Cache       *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

var supportedCommands map[string]cliCommand

func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range supportedCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(cfg *Config, args ...string) error {
	locationAreaMap, err := pokeapi.GetLocationAreas(cfg.Cache, cfg.NextURL)
	if err != nil {
		return err
	}

	for _, area := range locationAreaMap.Results {
		fmt.Println(area.Name)
	}
	cfg.NextURL = locationAreaMap.Next
	if locationAreaMap.Previous != nil {
		cfg.PreviousURL = *locationAreaMap.Previous
	} else {
		cfg.PreviousURL = ""
	}
	return nil
}

func commandMapBack(cfg *Config, args ...string) error {
	if cfg.PreviousURL == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	cfg.NextURL = cfg.PreviousURL
	cfg.PreviousURL = ""

	locationAreaMap, err := pokeapi.GetLocationAreas(cfg.Cache, cfg.NextURL)
	if err != nil {
		return err
	}

	for _, area := range locationAreaMap.Results {
		fmt.Println(area.Name)
	}
	cfg.NextURL = locationAreaMap.Next
	if locationAreaMap.Previous != nil {
		cfg.PreviousURL = *locationAreaMap.Previous
	} else {
		cfg.PreviousURL = ""
	}
	return nil
}

func commandExplore(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("please provide a location area name")
	}

	areaName := args[0]
	locationAreaDetails, err := pokeapi.GetLocationAreaDetails(cfg.Cache, cfg.NextURL, areaName)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %v...\n", locationAreaDetails.Name)
	fmt.Println("Found Pokemon: ")
	for _, value := range locationAreaDetails.PokemonEncounters {
		fmt.Printf("- %v\n", value.Pokemon.Name)
	}
	return nil
}

func initValues() map[string]cliCommand {
	supportedCommands = map[string]cliCommand{
		"exit": {
			"exit",
			"Exit the Pokedex",
			commandExit,
		},
		"help": {
			"help",
			"Displays a help message",
			commandHelp,
		},
		"map": {
			"map",
			"Displays a map",
			commandMap,
		},
		"mapb": {
			"mapb",
			"Displays the previous map",
			commandMapBack,
		},
		"explore": {
			"explore",
			"Explore a specific location area",
			commandExplore,
		},
	}
	return supportedCommands
}

func startRepl() {
	reader := bufio.NewScanner(os.Stdin)
	supportedCommandsLocal := initValues()

	// Create a cache with a 5-minute reap interval
	cache := pokecache.NewCache(5 * time.Minute)

	cfg := &Config{
		NextURL:     "https://pokeapi.co/api/v2/location-area?limit=20",
		PreviousURL: "",
		Cache:       cache,
	}
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		handler, ok := supportedCommandsLocal[commandName]
		if ok {
			if len(words) > 1 {
				cfg.NextURL = "https://pokeapi.co/api/v2/location-area"
				handler.callback(cfg, words[1])
			} else {
				handler.callback(cfg)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
