package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

type Stat struct {
	Name  string
	Value int
}

type Pokemon struct {
	Name           string
	BaseExperience int
	Height         int
	Weight         int
	Stats          []Stat
	Types          []string
}

type Config struct {
	Next     *string
	Previous *string
	Pokedex  map[string]Pokemon
}

func ptr(s string) *string {
	return &s
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	config := &Config{
		Next:     ptr("https://pokeapi.co/api/v2/location-area?limit=20"),
		Previous: nil,
		Pokedex:  make(map[string]Pokemon),
	}

	commands := getCommands()

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		words := cleanInput(input)

		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(config, words)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			break
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
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
			description: "List all Pokémon in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon-name>",
			description: "Attempt to catch a Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon-name>",
			description: "View details of a Pokémon you've caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show all Pokémon you've caught",
			callback:    commandPokedex,
		},
	}
}

func commandExit(cfg *Config, words []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, words []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}
