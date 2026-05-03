package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Ball struct {
	Key         string
	DisplayName string
	Bonus       int
}

var ballCatalog = map[string]Ball{
	"poke-ball": {
		Key:         "poke-ball",
		DisplayName: "Pokeball",
		Bonus:       0,
	},
	"great-ball": {
		Key:         "great-ball",
		DisplayName: "Great Ball",
		Bonus:       12,
	},
	"ultra-ball": {
		Key:         "ultra-ball",
		DisplayName: "Ultra Ball",
		Bonus:       25,
	},
}

var ballOrder = []string{"poke-ball", "great-ball", "ultra-ball"}

func defaultBall() Ball {
	return ballCatalog["poke-ball"]
}

func resolveBall(name string) (Ball, bool) {
	normalized := normalizeName(name)
	switch normalized {
	case "pokeball":
		normalized = "poke-ball"
	case "greatball":
		normalized = "great-ball"
	case "ultraball":
		normalized = "ultra-ball"
	}

	ball, ok := ballCatalog[normalized]
	return ball, ok
}

func acceptedBallKeys() []string {
	keys := make([]string, 0, len(ballOrder))
	for _, key := range ballOrder {
		keys = append(keys, key)
	}
	return keys
}

func parseCatchInput(words []string) (string, Ball, error) {
	if len(words) < 2 || len(words) > 3 {
		return "", Ball{}, errors.New("Usage: catch <pokemon-name> [ball-name]")
	}

	ball := defaultBall()
	if len(words) == 3 {
		resolvedBall, ok := resolveBall(words[2])
		if !ok {
			return "", Ball{}, fmt.Errorf("unknown ball %q. Try one of: %s", words[2], strings.Join(acceptedBallKeys(), ", "))
		}
		ball = resolvedBall
	}

	return normalizeName(words[1]), ball, nil
}

func calculateCatchChance(baseExperience int, ball Ball) int {
	chance := 70 - (baseExperience / 4) + ball.Bonus

	if chance < 15 {
		return 15
	}

	if chance > 90 {
		return 90
	}

	return chance
}

func attemptCatch(detail PokemonDetail, ball Ball, rng randomSource) (Pokemon, int, bool) {
	if rng == nil {
		rng = defaultRandom()
	}

	catchChance := calculateCatchChance(detail.BaseExperience, ball)
	roll := rng.Intn(100) + 1

	return pokemonFromDetail(detail), catchChance, roll <= catchChance
}

func pokemonFromDetail(detail PokemonDetail) Pokemon {
	stats := make([]Stat, len(detail.Stats))
	for i, stat := range detail.Stats {
		stats[i] = Stat{
			Name:  stat.Stat.Name,
			Value: stat.BaseStat,
		}
	}

	types := make([]string, len(detail.Types))
	for i, pokemonType := range detail.Types {
		types[i] = pokemonType.Type.Name
	}

	return Pokemon{
		Name:           normalizeName(detail.Name),
		BaseExperience: detail.BaseExperience,
		Height:         detail.Height,
		Weight:         detail.Weight,
		Stats:          stats,
		Types:          types,
		CaughtAt:       time.Now().UTC(),
	}
}

func chooseRandomEncounter(options []string, rng randomSource) (string, error) {
	if len(options) == 0 {
		return "", errors.New("no encounter options available")
	}

	if rng == nil {
		rng = defaultRandom()
	}

	return options[rng.Intn(len(options))], nil
}

func uniquePokemonNames(encounters []locationAreaPokemon) []string {
	seen := make(map[string]struct{}, len(encounters))
	names := make([]string, 0, len(encounters))

	for _, encounter := range encounters {
		name := normalizeName(encounter.Pokemon.Name)
		if name == "" {
			continue
		}

		if _, exists := seen[name]; exists {
			continue
		}

		seen[name] = struct{}{}
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

func sortedPokemonNames(pokedex map[string]Pokemon) []string {
	names := make([]string, 0, len(pokedex))
	for name := range pokedex {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}
