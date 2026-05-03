package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cfg, err := newConfig("", nil, os.Stdout, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Pokedex:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for {
		fmt.Fprint(cfg.Output, "Pokedex > ")

		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		command, exists := commands[words[0]]
		if !exists {
			cfg.println("Unknown command")
			continue
		}

		if err := command.callback(cfg, words); err != nil {
			cfg.println("Error:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
