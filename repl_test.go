package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "   spaced    out   words ",
			expected: []string{"spaced", "out", "words"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, tc := range cases {
		actual := cleanInput(tc.input)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("cleanInput(%q) = %v, expected %v", tc.input, actual, tc.expected)
		}
	}
}

func TestParseCatchInput(t *testing.T) {
	name, ball, err := parseCatchInput([]string{"catch", "pikachu", "greatball"})
	if err != nil {
		t.Fatalf("parseCatchInput returned error: %v", err)
	}

	if name != "pikachu" {
		t.Fatalf("expected pokemon name to be pikachu, got %q", name)
	}

	if ball.Key != "great-ball" {
		t.Fatalf("expected great-ball, got %q", ball.Key)
	}
}

func TestParseCatchInputDefaultsToPokeBall(t *testing.T) {
	name, ball, err := parseCatchInput([]string{"catch", "eevee"})
	if err != nil {
		t.Fatalf("parseCatchInput returned error: %v", err)
	}

	if name != "eevee" {
		t.Fatalf("expected pokemon name to be eevee, got %q", name)
	}

	if ball.Key != "poke-ball" {
		t.Fatalf("expected default ball to be poke-ball, got %q", ball.Key)
	}
}

func TestCalculateCatchChanceClamp(t *testing.T) {
	highChance := calculateCatchChance(4, ballCatalog["ultra-ball"])
	if highChance != 90 {
		t.Fatalf("expected high chance to clamp to 90, got %d", highChance)
	}

	lowChance := calculateCatchChance(400, ballCatalog["poke-ball"])
	if lowChance != 15 {
		t.Fatalf("expected low chance to clamp to 15, got %d", lowChance)
	}
}
