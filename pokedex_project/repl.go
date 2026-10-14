package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	pokecache "github.com/Mrfv2z/pokedex_project/internal/pokecache"
)

type config struct {
	commands map[string]cliCommand
	previous *string
	next     *string
	cache    *pokecache.Cache
	pokedex  map[string]EncounteredPokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string) error
}

func cleanInput(text string) []string {
	var ourSlice []string

	for _, word := range strings.Fields(text) {
		word = strings.ToLower(strings.TrimSpace(word))
		ourSlice = append(ourSlice, word)
	}

	return ourSlice
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		command := strings.ToLower(strings.TrimSpace(scanner.Text()))
		words := strings.Fields(command)

		if len(words) == 0 {
			continue
		}

		if cmd, ok := cfg.commands[words[0]]; ok {
			if err := cmd.callback(cfg, words[1:]...); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func Intn(n int) int {
	if n <= 0 {
		return -1
	}
	return rand.Intn(n)
}
