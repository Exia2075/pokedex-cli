package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display help for every command",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Show the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore <location-area-name>",
			description: "List Pokemon in an area and refresh its encounter pool",
			callback:    commandExplore,
		},
		"encounter": {
			name:        "encounter",
			description: "Roll a random wild Pokemon from the last explored area",
			callback:    commandEncounter,
		},
		"balls": {
			name:        "balls",
			description: "Show supported ball types and their catch bonuses",
			callback:    commandBalls,
		},
		"catch": {
			name:        "catch <pokemon-name> [ball-name]",
			description: "Attempt to catch a Pokemon with the selected ball",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon-name>",
			description: "View details of a Pokemon you've caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show all Pokemon you've caught",
			callback:    commandPokedex,
		},
	}
}

func commandExit(cfg *Config, words []string) error {
	cfg.println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, words []string) error {
	cfg.println("Welcome to the Pokedex!")
	cfg.println("Usage:")
	cfg.println()

	commands := getCommands()
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		command := commands[name]
		cfg.printf("%s: %s\n", command.name, command.description)
	}

	cfg.println()
	return nil
}

func commandMap(cfg *Config, words []string) error {
	if cfg.Next == nil {
		cfg.println("no more locations")
		return nil
	}

	return printLocationPage(cfg, *cfg.Next)
}

func commandMapBack(cfg *Config, words []string) error {
	if cfg.Previous == nil {
		cfg.println("you're on the first page")
		return nil
	}

	return printLocationPage(cfg, *cfg.Previous)
}

func printLocationPage(cfg *Config, url string) error {
	data, err := cfg.API.ListLocations(url)
	if err != nil {
		return err
	}

	for _, area := range data.Results {
		cfg.println(area.Name)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous
	return nil
}

func commandExplore(cfg *Config, words []string) error {
	if len(words) != 2 {
		cfg.println("Usage: explore <location-area-name>")
		return nil
	}

	locationName := normalizeName(words[1])
	detail, err := cfg.API.GetLocationArea(locationName)
	if err != nil {
		return err
	}

	encounterPool := uniquePokemonNames(detail.PokemonEncounters)
	cfg.LastExploredArea = locationName
	cfg.LastEncounterOptions = encounterPool
	if err := cfg.save(); err != nil {
		return err
	}

	cfg.printf("Exploring %s...\n", locationName)

	if len(encounterPool) == 0 {
		cfg.println("No Pokemon found in this location area.")
		return nil
	}

	cfg.println("Pokemon found:")
	for _, name := range encounterPool {
		cfg.printf(" - %s\n", name)
	}

	cfg.printf("Encounter pool updated. Use `encounter` to roll a wild Pokemon from %s.\n", locationName)
	return nil
}

func commandEncounter(cfg *Config, words []string) error {
	if len(cfg.LastEncounterOptions) == 0 {
		cfg.println("Explore a location first so there is an encounter pool to roll from.")
		return nil
	}

	pokemonName, err := chooseRandomEncounter(cfg.LastEncounterOptions, cfg.Random)
	if err != nil {
		return err
	}

	cfg.printf("A wild %s appeared in %s!\n", pokemonName, cfg.LastExploredArea)
	cfg.printf("Try: catch %s %s\n", pokemonName, defaultBall().Key)
	return nil
}

func commandBalls(cfg *Config, words []string) error {
	cfg.println("Available balls:")
	for _, key := range ballOrder {
		ball := ballCatalog[key]
		cfg.printf(" - %s (%+d catch bonus)\n", ball.Key, ball.Bonus)
	}
	return nil
}

func commandCatch(cfg *Config, words []string) error {
	pokemonName, ball, err := parseCatchInput(words)
	if err != nil {
		cfg.println(err.Error())
		return nil
	}

	detail, err := cfg.API.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("could not find pokemon %s: %w", pokemonName, err)
	}

	cfg.printf("Throwing a %s at %s...\n", ball.DisplayName, detail.Name)

	pokemon, catchChance, caught := attemptCatch(detail, ball, cfg.Random)
	if !caught {
		cfg.printf("%s escaped! Catch chance was %d%%.\n", detail.Name, catchChance)
		return nil
	}

	cfg.Pokedex[pokemon.Name] = pokemon
	if err := cfg.save(); err != nil {
		return err
	}

	cfg.printf("%s was caught! Catch chance was %d%%.\n", pokemon.Name, catchChance)
	return nil
}

func commandInspect(cfg *Config, words []string) error {
	if len(words) != 2 {
		cfg.println("Usage: inspect <pokemon-name>")
		return nil
	}

	pokemonName := normalizeName(words[1])
	pokemon, exists := cfg.Pokedex[pokemonName]
	if !exists {
		cfg.println("you have not caught that pokemon")
		return nil
	}

	cfg.printf("Name: %s\n", pokemon.Name)
	cfg.printf("Height: %d\n", pokemon.Height)
	cfg.printf("Weight: %d\n", pokemon.Weight)
	cfg.printf("Caught: %s\n", pokemon.CaughtAt.In(time.Local).Format(time.RFC1123))

	cfg.println("Stats:")
	for _, stat := range pokemon.Stats {
		cfg.printf("  - %s: %d\n", stat.Name, stat.Value)
	}

	cfg.println("Types:")
	for _, pokemonType := range pokemon.Types {
		cfg.printf("  - %s\n", pokemonType)
	}

	return nil
}

func commandPokedex(cfg *Config, words []string) error {
	if len(cfg.Pokedex) == 0 {
		cfg.println("Your Pokedex is empty.")
		return nil
	}

	cfg.println("Your Pokedex:")
	for _, name := range sortedPokemonNames(cfg.Pokedex) {
		cfg.printf(" - %s\n", name)
	}

	if cfg.LastExploredArea != "" && len(cfg.LastEncounterOptions) > 0 {
		cfg.printf("Last encounter pool: %s (%s)\n", cfg.LastExploredArea, strings.Join(cfg.LastEncounterOptions, ", "))
	}

	return nil
}
