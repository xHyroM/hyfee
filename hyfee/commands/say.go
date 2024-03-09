package commands

import (
	"hyros_coffee/handler"
	"hyros_coffee/hyfee"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func Say(bot *hyfee.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name: "say",
			Description: "Make the bot say something",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name: "message",
					Description: "The message to say",
					Required: true,
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": sayHandler(bot),
		},
	}
}

func sayHandler(_ *hyfee.Bot) handler.CommandHandler {
	return func(e *events.ApplicationCommandInteractionCreate) error {
		message := e.SlashCommandInteractionData().String("message")

		return e.CreateMessage(discord.MessageCreate{
			Content: "You said: " + message,
			Flags: discord.MessageFlagEphemeral,
		})
	}
}