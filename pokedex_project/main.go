package main

import (
	"time"

	pokecache "github.com/Mrfv2z/pokedex_project/internal/pokecache"
)

func main() {
	cfg := &config{
		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    helpcommand,
			},
			"map": {
				name:        "map",
				description: "It displays the names of 20 location areas in the Pokemon world.",
				callback:    mapcommand,
			},
			"mapb": {
				name:        "mapb",
				description: "It displays the names of the previous 20 location areas seen in Pokemon world.",
				callback:    mapbcommand,
			},
			"explore": {
				name:        "explore",
				description: "It displays all Pokemon in selected location",
				callback:    explorecommand,
			},
			"catch": {
				name:        "catch",
				description: "Throws pokeball at selected Pokemon",
				callback:    catchcommand,
			},
			"inspect": {
				name:        "inspect",
				description: "Checks data on caught pokemon",
				callback:    inspectcommand,
			},
			"pokedex": {
				name:        "pokedex",
				description: "Checks Pokemon caught so far",
				callback:    pokedexcommand,
			},
		},
		cache: pokecache.NewCache(8000 * time.Millisecond),
	}
	startRepl(cfg)
}
