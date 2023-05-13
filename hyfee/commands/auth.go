package commands

import (
	"hyros_coffee/handler"
	"hyros_coffee/hyfee"
	"os"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var baseUrl string

func Auth(bot *hyfee.Bot) handler.Command {
	baseUrl = os.Getenv("BASE_URL")
	
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name: "auth",
			Description: "Authenicate yourself with Discord OAuth2",
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": authHandler(bot),
		},
	}
}

func authHandler(bot *hyfee.Bot) handler.CommandHandler {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		url := bot.OAuth2Client.GenerateAuthorizationURL(baseUrl+"/callback", 0, *event.GuildID(), false, discord.OAuth2ScopeGuilds, discord.OAuth2ScopeIdentify)

		return event.CreateMessage(discord.MessageCreate{
			Content: "Click [here](" + url + ") to authenicate yourself with Discord OAuth2",
			Flags: discord.MessageFlagEphemeral,
		})
	}
}