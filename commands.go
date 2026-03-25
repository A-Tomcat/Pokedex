package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func cmd(cfg *config, args ...string) map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 next Locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays 20 previous Locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <location>",
			description: "Displays Pokemon at Location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Adds <pokemon> to users pokedex",
			callback:    commandCatch,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	keys := []string{}
	available := cmd(cfg)
	for name := range available {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		command := available[name]
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.nextLocationURL != nil {
		url = *cfg.nextLocationURL
	}
	var body []byte
	locations := RespShallowLocations{}

	body, err := cfg.cache.Check(url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &locations); err != nil {
		return err
	}
	//Update the config!
	cfg.nextLocationURL = locations.Next
	cfg.previousLocationURL = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.previousLocationURL == nil {
		return fmt.Errorf("you're on the first page")
	}
	url := *cfg.previousLocationURL

	var body []byte
	locations := RespShallowLocations{}
	body, err := cfg.cache.Check(url)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &locations); err != nil {
		return err
	}

	//Update the config!
	cfg.nextLocationURL = locations.Next
	cfg.previousLocationURL = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) < 1 {
		return errors.New("Provide a location: explore <location>")
	}
	loc := args[0]
	url := "https://pokeapi.co/api/v2/location-area/"
	area_url := url + loc + "/"

	var location LocationArea
	var body []byte

	body, err := cfg.cache.Check(area_url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &location); err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", loc)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) < 1 {
		return errors.New("Provide a Pokemon: catch <pokemon>")
	}
	pk := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + pk + "/"
	var pokemon Pokemon
	var body []byte

	body, err := cfg.cache.Check(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	chance := rand.Intn(pokemon.BaseExperience)

	if float64(chance) > float64(pokemon.BaseExperience)*0.8 {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.caughtPokemon[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}
