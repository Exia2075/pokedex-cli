package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Exia2075/pokedexcli/internal/pokecache"
)

const locationListURL = "https://pokeapi.co/api/v2/location-area?limit=20"

var cache = pokecache.NewCache(5 * time.Second)

type fetchFunc func(string, any) error

type pokemonAPI interface {
	ListLocations(string) (LocationAreaResponse, error)
	GetLocationArea(string) (LocationAreaDetail, error)
	GetPokemon(string) (PokemonDetail, error)
}

type pokeAPIClient struct {
	fetch fetchFunc
}

type LocationAreaResponse struct {
	Count    int                  `json:"count"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Results  []locationAreaResult `json:"results"`
}

type locationAreaResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaDetail struct {
	PokemonEncounters []locationAreaPokemon `json:"pokemon_encounters"`
}

type locationAreaPokemon struct {
	Pokemon namedAPIResource `json:"pokemon"`
}

type namedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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

func newPokeAPIClient(fetch fetchFunc) *pokeAPIClient {
	if fetch == nil {
		fetch = fetchJSON
	}

	return &pokeAPIClient{fetch: fetch}
}

func (c *pokeAPIClient) ListLocations(url string) (LocationAreaResponse, error) {
	var data LocationAreaResponse
	err := c.fetch(url, &data)
	return data, err
}

func (c *pokeAPIClient) GetLocationArea(name string) (LocationAreaDetail, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", normalizeName(name))

	var detail LocationAreaDetail
	err := c.fetch(url, &detail)
	return detail, err
}

func (c *pokeAPIClient) GetPokemon(name string) (PokemonDetail, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", normalizeName(name))

	var detail PokemonDetail
	err := c.fetch(url, &detail)
	return detail, err
}

func fetchJSON(url string, out any) error {
	if data, found := cache.Get(url); found {
		return json.Unmarshal(data, out)
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("pokeapi returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cache.Add(url, body)

	return json.Unmarshal(body, out)
}
