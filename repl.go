package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/khazakar/pokedexcli/internal/pokeapi"
	"github.com/khazakar/pokedexcli/internal/pokecache"
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

var pokedexStorage map[string]pokeapi.PokemonDetails

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

func commandCatch(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("please provide a Pokemon name")
	}
	// Set the base Pokemon API URL
	pokemonURL := "https://pokeapi.co/api/v2/pokemon"
	fmt.Printf("Throwing a Pokeball at %v...\n", args[0])
	pokemonDetails, err := pokeapi.GetPokemonDetails(cfg.Cache, pokemonURL, args[0])
	if err != nil {
		fmt.Printf("%v was not found\n", args[0])
		return err
	}
	probability := float64(pokemonDetails.BaseExperience) / 200.0
	if rand.Float64() < probability {
		fmt.Printf("%v was caught!\n", pokemonDetails.Name)
		pokedexStorage[pokemonDetails.Name] = *pokemonDetails
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%v escaped!\n", pokemonDetails.Name)
	}

	return nil
}

func commandInspect(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("please provide a Pokemon name")
	}

	pokemonName := args[0]
	pokemon, ok := pokedexStorage[pokemonName]
	if !ok {
		fmt.Printf("you have not caught %s\n", pokemonName)
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *Config, args ...string) error {
	for _, val := range pokedexStorage {
		fmt.Printf("- %s\n", val.Name)
	}
	return nil
}

func initCommands() map[string]cliCommand {
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
		"catch": {
			"catch",
			"Catch a Pokemon",
			commandCatch,
		},
		"inspect": {
			"inspect",
			"Inspect a caught Pokemon",
			commandInspect,
		},
		"pokedex": {
			"pokedex",
			"Displays all caught Pokemon",
			commandPokedex,
		},
	}
	return supportedCommands
}

func startRepl() {
	reader := bufio.NewScanner(os.Stdin)
	supportedCommandsLocal := initCommands()
	pokedexStorage = make(map[string]pokeapi.PokemonDetails)
	// Create a cache with a 5-minute reap interval
	cache := pokecache.NewCache(5 * time.Minute)

	cfg := &Config{
		NextURL:     "https://pokeapi.co/api/v2/location-area",
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
				handler.callback(cfg, words[1:]...)
			} else {
				// Only reset the NextURL for the first call to map command
				// or for commands other than map/mapb
				if commandName != "map" && commandName != "mapb" && cfg.NextURL == "" {
					cfg.NextURL = "https://pokeapi.co/api/v2/location-area?limit=20"
				}
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
