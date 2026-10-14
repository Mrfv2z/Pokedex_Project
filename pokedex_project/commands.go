package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type JSONconvert struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type EncounteredPokemon struct {
	Name           string         `json:"name"`
	URL            string         `json:"url"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          []PokemonStats `json:"stats"`
}

type PokemonStats struct {
	BaseStats int `json:"base_stat"`
	Stat      struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

type PokemonEncounter struct {
	Pokemon EncounteredPokemon `json:"pokemon"`
}

type ExploreResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func commandExit(_ *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func helpcommand(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, command := range cfg.commands {
		fmt.Printf("\n%s: %s\n", command.name, command.description)
	}
	return nil
}

func mapcommand(cfg *config, args ...string) error {
	myStruct := JSONconvert{}
	link := "https://pokeapi.co/api/v2/location-area/"
	URLfetch := link
	if cfg.next != nil {
		URLfetch = *cfg.next
	}

	var body []byte
	if check, ok := cfg.cache.Get(URLfetch); ok {
		body = check
	} else {
		res, err := http.Get(URLfetch)
		if err != nil {
			log.Fatal(err)
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		cfg.cache.Add(URLfetch, body)
	}
	if err := json.Unmarshal(body, &myStruct); err != nil {
		log.Fatal(err)
	}
	for i := range myStruct.Results {
		fmt.Println(myStruct.Results[i].Name)
	}
	cfg.next = myStruct.Next
	cfg.previous = myStruct.Previous
	return nil
}

func mapbcommand(cfg *config, args ...string) error {
	myStruct := JSONconvert{}
	link := "https://pokeapi.co/api/v2/location-area/"
	URLfetch := link
	if cfg.previous == nil {
		return fmt.Errorf("you're on the first page")
	}
	URLfetch = *cfg.previous
	var body []byte
	if check, ok := cfg.cache.Get(URLfetch); ok {
		body = check
	} else {
		res, err := http.Get(URLfetch)
		if err != nil {
			log.Fatal(err)
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		cfg.cache.Add(URLfetch, body)
	}
	if err := json.Unmarshal(body, &myStruct); err != nil {
		log.Fatal(err)
	}
	for i := range myStruct.Results {
		fmt.Println(myStruct.Results[i].Name)
	}
	cfg.next = myStruct.Next
	cfg.previous = myStruct.Previous
	return nil
}

func explorecommand(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Please provide a location area name")
	}

	myStruct := ExploreResponse{}
	link := "https://pokeapi.co/api/v2/location-area/"
	URLfetch := fmt.Sprintf("%s%s/", link, args[0])
	var body []byte
	if check, ok := cfg.cache.Get(URLfetch); ok {
		body = check
	} else {
		res, err := http.Get(URLfetch)
		if err != nil {
			return fmt.Errorf("Failed to fetch location: %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Failed to read response body: %w", err)
		}

		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d", res.StatusCode)
		}
		cfg.cache.Add(URLfetch, body)
	}

	if err := json.Unmarshal(body, &myStruct); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for i := range myStruct.PokemonEncounters {
		fmt.Printf(" - %s\n", myStruct.PokemonEncounters[i].Pokemon.Name)
	}
	return nil
}

func catchcommand(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Please provide a Pokemon name")
	}

	myStruct := EncounteredPokemon{}
	URLfetch := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", args[0])

	res, err := http.Get(URLfetch)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("Please provide a valid Pokemon name")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &myStruct); err != nil {
		return err
	}

	var odds int
	if myStruct.BaseExperience <= 100 {
		odds = 65
	} else if myStruct.BaseExperience <= 200 {
		odds = 40
	} else {
		odds = 20
	}

	rng := Intn(100)
	caught := rng < odds

	fmt.Printf("Throwing a Pokeball at %s...\n", myStruct.Name)
	if caught {
		fmt.Printf("%s was caught!\n", myStruct.Name)
		if cfg.pokedex == nil {
			cfg.pokedex = make(map[string]EncounteredPokemon)
		}
		cfg.pokedex[myStruct.Name] = myStruct
		fmt.Print("You may now inspect it with the inspect command.\n")
	} else {
		fmt.Printf("%s escaped!\n", myStruct.Name)
	}

	return nil
}

func inspectcommand(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Please provide a Pokemon name")
	}

	pokemon, ok := cfg.pokedex[args[0]]
	if !ok {
		return fmt.Errorf("you have not caught that Pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")

	for _, stat := range pokemon.Stats {
		fmt.Printf("  %s: %d\n", stat.Stat.Name, stat.BaseStats)
	}

	return nil
}

func pokedexcommand(cfg *config, args ...string) error {
	fmt.Print("Your Pokedex:\n")
	for i := range cfg.pokedex {
		fmt.Printf(" - %s\n", cfg.pokedex[i].Name)
	}
	return nil
}
