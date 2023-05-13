package utils

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func GetFocusedOption(e *events.AutocompleteInteractionCreate) discord.AutocompleteOption {
	var focused discord.AutocompleteOption

	for _, opt := range e.Data.Options {
		if opt.Focused {
			focused = opt
			break
		}
	}

	return focused
}

func IfThenElse[T any](condition bool, a func() T, b func() T) T {
	if condition {
			return a()
	}
	return b()
}

func SnakeCaseToPascalCaseWithSpaces(text string) string {
	words := strings.Split(text, "_")

	for i, word := range words {
			words[i] = strings.Title(strings.ToLower(word))
	}

	return strings.Join(words, " ")
}