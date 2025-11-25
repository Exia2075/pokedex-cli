package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/Exia2075/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(5 * time.Second)

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonDetail struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func fetch(url string, out interface{}) error {
	if data, found := cache.Get(url); found {
		// fmt.Println("(cache hit)")
		return json.Unmarshal(data, out)
	}

	// fmt.Println("(cache miss)")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cache.Add(url, b)

	return json.Unmarshal(b, out)
}

func commandMap(cfg *Config, words []string) error {
	if cfg.Next == nil {
		fmt.Println("no more locations")
		return nil
	}

	var data LocationAreaResponse
	if err := fetch(*cfg.Next, &data); err != nil {
		return err
	}

	for _, r := range data.Results {
		fmt.Println(r.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous
	return nil
}

func commandMapBack(cfg *Config, words []string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	var data LocationAreaResponse
	if err := fetch(*cfg.Previous, &data); err != nil {
		return err
	}

	for _, r := range data.Results {
		fmt.Println(r.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous
	return nil
}

func commandExplore(cfg *Config, words []string) error {
	if len(words) < 2 {
		fmt.Println("Usage: explore <location-area-name>")
		return nil
	}

	locName := words[1]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locName)

	var detail LocationAreaDetail
	if err := fetch(url, &detail); err != nil {
		return err
	}

	if len(detail.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon found in this location area.")
		return nil
	}

	for _, p := range detail.PokemonEncounters {
		fmt.Println(p.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *Config, words []string) error {
	if len(words) < 2 {
		fmt.Println("Usage: catch <pokemon-name>")
		return nil
	}

	pokeName := words[1]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokeName)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokeName)
	var detail PokemonDetail
	if err := fetch(url, &detail); err != nil {
		return fmt.Errorf("could not find pokemon %s: %w", pokeName, err)
	}

	chance := rand.Intn(100)
	if chance > detail.BaseExperience {
		fmt.Printf("%s was caught!\n", detail.Name)

		stats := make([]Stat, len(detail.Stats))
		for i, s := range detail.Stats {
			stats[i] = Stat{
				Name:  s.Stat.Name,
				Value: s.BaseStat,
			}
		}

		types := make([]string, len(detail.Types))
		for i, t := range detail.Types {
			types[i] = t.Type.Name
		}

		cfg.Pokedex[detail.Name] = Pokemon{
			Name:           detail.Name,
			BaseExperience: detail.BaseExperience,
			Height:         detail.Height,
			Weight:         detail.Weight,
			Stats:          stats,
			Types:          types,
		}
	} else {
		fmt.Printf("%s escaped!\n", detail.Name)
	}

	return nil
}

func commandInspect(cfg *Config, words []string) error {
	if len(words) < 2 {
		fmt.Println("Usage: inspect <pokemon-name>")
		return nil
	}

	name := words[1]
	pokemon, ok := cfg.Pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Name, s.Value)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t)
	}

	return nil
}

func commandPokedex(cfg *Config, words []string) error {
	if len(cfg.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range cfg.Pokedex {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
