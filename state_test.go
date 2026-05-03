package main

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestLoadStateMissingFileReturnsEmptyState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	state, err := loadState(path)
	if err != nil {
		t.Fatalf("loadState returned error: %v", err)
	}

	if len(state.Pokedex) != 0 {
		t.Fatalf("expected empty pokedex, got %d entries", len(state.Pokedex))
	}

	if state.LastExploredArea != "" {
		t.Fatalf("expected empty LastExploredArea, got %q", state.LastExploredArea)
	}
}

func TestSaveStateRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")

	expected := savedState{
		Pokedex: map[string]Pokemon{
			"pikachu": {
				Name:           "pikachu",
				BaseExperience: 112,
				Height:         4,
				Weight:         60,
				Stats: []Stat{
					{Name: "speed", Value: 90},
				},
				Types:    []string{"electric"},
				CaughtAt: time.Unix(1700000000, 0).UTC(),
			},
		},
		LastExploredArea:     "kanto-route-2",
		LastEncounterOptions: []string{"pikachu", "pidgey"},
	}

	if err := saveState(path, expected); err != nil {
		t.Fatalf("saveState returned error: %v", err)
	}

	actual, err := loadState(path)
	if err != nil {
		t.Fatalf("loadState returned error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("round-trip mismatch:\nexpected: %#v\nactual: %#v", expected, actual)
	}
}
