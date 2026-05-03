package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type commandCallback func(*Config, []string) error

type cliCommand struct {
	name        string
	description string
	callback    commandCallback
}

type Stat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Pokemon struct {
	Name           string    `json:"name"`
	BaseExperience int       `json:"base_experience"`
	Height         int       `json:"height"`
	Weight         int       `json:"weight"`
	Stats          []Stat    `json:"stats"`
	Types          []string  `json:"types"`
	CaughtAt       time.Time `json:"caught_at"`
}

type randomSource interface {
	Intn(int) int
}

type Config struct {
	Next                 *string
	Previous             *string
	Pokedex              map[string]Pokemon
	LastExploredArea     string
	LastEncounterOptions []string

	API         pokemonAPI
	Output      io.Writer
	Random      randomSource
	StoragePath string
}

func newConfig(storagePath string, api pokemonAPI, output io.Writer, rng randomSource) (*Config, error) {
	if storagePath == "" {
		var err error
		storagePath, err = defaultStoragePath()
		if err != nil {
			return nil, err
		}
	}

	state, err := loadState(storagePath)
	if err != nil {
		return nil, err
	}

	if api == nil {
		api = newPokeAPIClient(nil)
	}

	if output == nil {
		output = os.Stdout
	}

	if rng == nil {
		rng = defaultRandom()
	}

	return &Config{
		Next:                 ptr(locationListURL),
		Previous:             nil,
		Pokedex:              state.Pokedex,
		LastExploredArea:     state.LastExploredArea,
		LastEncounterOptions: append([]string(nil), state.LastEncounterOptions...),
		API:                  api,
		Output:               output,
		Random:               rng,
		StoragePath:          storagePath,
	}, nil
}

func defaultRandom() randomSource {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func ptr(s string) *string {
	return &s
}

func (cfg *Config) println(args ...any) {
	fmt.Fprintln(cfg.Output, args...)
}

func (cfg *Config) printf(format string, args ...any) {
	fmt.Fprintf(cfg.Output, format, args...)
}

func (cfg *Config) save() error {
	return saveState(cfg.StoragePath, cfg.snapshot())
}

func (cfg *Config) snapshot() savedState {
	return savedState{
		Pokedex:              clonePokedex(cfg.Pokedex),
		LastExploredArea:     cfg.LastExploredArea,
		LastEncounterOptions: append([]string(nil), cfg.LastEncounterOptions...),
	}
}
