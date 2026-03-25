package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cmd(cfg *config) map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    func() error { return commandHelp(cfg) },
		},
		"map": {
			name:        "map",
			description: "Displays 20 next Locations",
			callback:    func() error { return commandMap(cfg) },
		},
		"mapb": {
			name:        "mapb",
			description: "Displays 20 previous Locations",
			callback:    func() error { return commandMapb(cfg) },
		},
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
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

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.nextLocationURL != nil {
		url = *cfg.nextLocationURL
	}
	var body []byte
	locations := RespShallowLocations{}
	if val, ok := cfg.cache.Get(url); ok {
		body = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, body)
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

func commandMapb(cfg *config) error {
	if cfg.previousLocationURL == nil {
		return fmt.Errorf("you're on the first page")
	}
	url := *cfg.previousLocationURL

	var body []byte
	locations := RespShallowLocations{}
	if val, ok := cfg.cache.Get(url); ok {
		body = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, body)
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
