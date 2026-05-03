package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type savedState struct {
	Pokedex              map[string]Pokemon `json:"pokedex"`
	LastExploredArea     string             `json:"last_explored_area,omitempty"`
	LastEncounterOptions []string           `json:"last_encounter_options,omitempty"`
}

func defaultStoragePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".pokedexcli.json"), nil
}

func loadState(path string) (savedState, error) {
	state := savedState{
		Pokedex: make(map[string]Pokemon),
	}

	if path == "" {
		return state, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return state, nil
		}
		return state, err
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return state, nil
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return state, err
	}

	if state.Pokedex == nil {
		state.Pokedex = make(map[string]Pokemon)
	} else {
		state.Pokedex = clonePokedex(state.Pokedex)
	}

	state.LastEncounterOptions = append([]string(nil), state.LastEncounterOptions...)

	return state, nil
}

func saveState(path string, state savedState) error {
	if path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	payload := savedState{
		Pokedex:              clonePokedex(state.Pokedex),
		LastExploredArea:     state.LastExploredArea,
		LastEncounterOptions: append([]string(nil), state.LastEncounterOptions...),
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(path), "pokedex-state-*.json")
	if err != nil {
		return err
	}

	tempName := tempFile.Name()
	defer os.Remove(tempName)

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	return os.Rename(tempName, path)
}

func clonePokedex(src map[string]Pokemon) map[string]Pokemon {
	dst := make(map[string]Pokemon, len(src))

	for name, pokemon := range src {
		copied := pokemon
		copied.Stats = append([]Stat(nil), pokemon.Stats...)
		copied.Types = append([]string(nil), pokemon.Types...)
		dst[name] = copied
	}

	return dst
}
