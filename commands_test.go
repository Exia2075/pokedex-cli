package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

type stubAPI struct {
	listLocationsResp LocationAreaResponse
	locationAreaResp  LocationAreaDetail
	pokemonResp       PokemonDetail
	listLocationsErr  error
	locationAreaErr   error
	pokemonErr        error
}

func (s stubAPI) ListLocations(url string) (LocationAreaResponse, error) {
	return s.listLocationsResp, s.listLocationsErr
}

func (s stubAPI) GetLocationArea(name string) (LocationAreaDetail, error) {
	return s.locationAreaResp, s.locationAreaErr
}

func (s stubAPI) GetPokemon(name string) (PokemonDetail, error) {
	return s.pokemonResp, s.pokemonErr
}

type fixedRand struct {
	value int
}

func (r fixedRand) Intn(n int) int {
	if n <= 0 {
		return 0
	}

	if r.value >= n {
		return n - 1
	}

	return r.value
}

func TestCommandExploreUpdatesEncounterPoolAndPersistsState(t *testing.T) {
	var output bytes.Buffer
	path := filepath.Join(t.TempDir(), "state.json")

	cfg, err := newConfig(
		path,
		stubAPI{
			locationAreaResp: LocationAreaDetail{
				PokemonEncounters: []locationAreaPokemon{
					{Pokemon: namedAPIResource{Name: "pikachu"}},
					{Pokemon: namedAPIResource{Name: "pidgey"}},
					{Pokemon: namedAPIResource{Name: "pikachu"}},
				},
			},
		},
		&output,
		fixedRand{value: 0},
	)
	if err != nil {
		t.Fatalf("newConfig returned error: %v", err)
	}

	if err := commandExplore(cfg, []string{"explore", "kanto-route-2"}); err != nil {
		t.Fatalf("commandExplore returned error: %v", err)
	}

	expectedPool := []string{"pidgey", "pikachu"}
	if strings.TrimSpace(cfg.LastExploredArea) != "kanto-route-2" {
		t.Fatalf("expected LastExploredArea to be updated, got %q", cfg.LastExploredArea)
	}

	if len(cfg.LastEncounterOptions) != len(expectedPool) {
		t.Fatalf("expected %d encounter options, got %d", len(expectedPool), len(cfg.LastEncounterOptions))
	}

	for i, name := range expectedPool {
		if cfg.LastEncounterOptions[i] != name {
			t.Fatalf("expected encounter option %d to be %q, got %q", i, name, cfg.LastEncounterOptions[i])
		}
	}

	state, err := loadState(path)
	if err != nil {
		t.Fatalf("loadState returned error: %v", err)
	}

	if state.LastExploredArea != "kanto-route-2" {
		t.Fatalf("expected persisted LastExploredArea to be kanto-route-2, got %q", state.LastExploredArea)
	}
}

func TestCommandEncounterUsesSavedPool(t *testing.T) {
	var output bytes.Buffer

	cfg := &Config{
		Output:               &output,
		Random:               fixedRand{value: 1},
		LastExploredArea:     "viridian-forest",
		LastEncounterOptions: []string{"caterpie", "weedle", "pikachu"},
	}

	if err := commandEncounter(cfg, []string{"encounter"}); err != nil {
		t.Fatalf("commandEncounter returned error: %v", err)
	}

	if !strings.Contains(output.String(), "weedle") {
		t.Fatalf("expected encounter output to mention weedle, got %q", output.String())
	}
}

func TestCommandCatchStoresPokemonAndPersistsState(t *testing.T) {
	var output bytes.Buffer
	path := filepath.Join(t.TempDir(), "state.json")

	cfg, err := newConfig(
		path,
		stubAPI{
			pokemonResp: PokemonDetail{
				Name:           "pikachu",
				BaseExperience: 112,
				Height:         4,
				Weight:         60,
				Stats: []struct {
					BaseStat int `json:"base_stat"`
					Stat     struct {
						Name string `json:"name"`
					} `json:"stat"`
				}{
					{
						BaseStat: 90,
						Stat: struct {
							Name string `json:"name"`
						}{Name: "speed"},
					},
				},
				Types: []struct {
					Type struct {
						Name string `json:"name"`
					} `json:"type"`
				}{
					{
						Type: struct {
							Name string `json:"name"`
						}{Name: "electric"},
					},
				},
			},
		},
		&output,
		fixedRand{value: 0},
	)
	if err != nil {
		t.Fatalf("newConfig returned error: %v", err)
	}

	if err := commandCatch(cfg, []string{"catch", "pikachu", "ultra-ball"}); err != nil {
		t.Fatalf("commandCatch returned error: %v", err)
	}

	if _, exists := cfg.Pokedex["pikachu"]; !exists {
		t.Fatalf("expected pikachu to be stored in pokedex")
	}

	state, err := loadState(path)
	if err != nil {
		t.Fatalf("loadState returned error: %v", err)
	}

	if _, exists := state.Pokedex["pikachu"]; !exists {
		t.Fatalf("expected pikachu to be persisted to disk")
	}
}
