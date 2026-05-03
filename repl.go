package main

import "strings"

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)
	for i := range words {
		words[i] = normalizeName(words[i])
	}

	return words
}

func normalizeName(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}
