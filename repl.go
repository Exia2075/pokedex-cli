package main

import (
	"strings"
)

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)

	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)

	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return words
}
