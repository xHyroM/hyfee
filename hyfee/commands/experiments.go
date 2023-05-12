package commands

import (
	"hyros_coffee/handler"
	"hyros_coffee/hyfee"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func Experiments(bot *hyfee.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name: "experiments",
			Description: "Discord experiments",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name: "get",
					Description: "Get information about an experiment",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name: "experiment",
							Description: "The experiment id",
							Required: true,
							Autocomplete: true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name: "eligible",
					Description: "Check if a guild is eligible for an experiment",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name: "experiment",
							Description: "The experiment id",
							Required: true,
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionString{
							Name: "guild",
							Description: "The guild id",
							Required: true,
						},
					},
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"get": getExperimentsHandler(bot),
			"eligible": eligibleExperimentsHandler(bot),
		},
		AutocompleteHandlers: map[string]handler.AutocompleteHandler{
			"get": autocompleteHandler(bot),
			"eligible": autocompleteHandler(bot),
		},
	}
}

func getExperimentsHandler(bot *hyfee.Bot) handler.CommandHandler {
	return func(e *events.ApplicationCommandInteractionCreate) error {
		id := e.SlashCommandInteractionData().String("id")

		return e.CreateMessage(discord.MessageCreate{
			Content: "You said: " + id,
			Flags: discord.MessageFlagEphemeral,
		})
	}
}

func eligibleExperimentsHandler(bot *hyfee.Bot) handler.CommandHandler {
	return func(e *events.ApplicationCommandInteractionCreate) error {
		id := e.SlashCommandInteractionData().String("id")
		guild := e.SlashCommandInteractionData().String("guild")

		return e.CreateMessage(discord.MessageCreate{
			Content: "You said: " + id + " " + guild,
			Flags: discord.MessageFlagEphemeral,
		})
	}
}

func autocompleteHandler(bot *hyfee.Bot) handler.AutocompleteHandler {
	return func(e *events.AutocompleteInteractionCreate) error {
		return e.Result([]discord.AutocompleteChoice{
			discord.AutocompleteChoiceString{
				Name: "test experiments",
				Value: "test_experiments",
			},
		})
	}
}