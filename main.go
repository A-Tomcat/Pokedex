package main

import (
	"Pokedex/internal/pokecache"
	"bufio"
	"fmt"
	"os"
	"time"
)

type config struct {
	nextLocationURL     *string
	previousLocationURL *string
	cache               *pokecache.Cache
	caughtPokemon       map[string]Pokemon
}

func main() {
	cfg := &config{
		cache:         pokecache.NewCache(5 * time.Second),
		caughtPokemon: make(map[string]Pokemon),
	}
	cmds := cmd(cfg)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		prompt := cleanInput(scanner.Text())
		if len(prompt) == 0 {
			continue
		}
		commandName := prompt[0]

		args := []string{}
		if len(prompt) > 1 {
			args = prompt[1:]
		}
		command, ok := cmds[commandName]
		if ok {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command, try 'help'")
			continue
		}
	}

}
