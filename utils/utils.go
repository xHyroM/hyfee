package utils

import (
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