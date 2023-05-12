package main

import (
	"context"
	"hyros_coffe/hyfee"
	"hyros_coffe/hyfee/commands"
	"hyros_coffe/hyfee/listeners"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/snowflake/v2"
	_ "github.com/joho/godotenv/autoload"
)

var (
	botToken = os.Getenv("DISCORD_TOKEN")
	guildId, err =	snowflake.Parse("1046534628577640528")
)

func main() {
	bot := hyfee.New()

	bot.Handler.AddCommands(
		commands.Say(bot),
	)

	if err := bot.Setup(
		listeners.Ready(bot),
	); err != nil {
		bot.Logger.Fatal(err)
	}

	defer bot.Client.Close(context.TODO())

	bot.Handler.SyncCommands(bot.Client, guildId)

	if err = bot.Start() ; err != nil {
		bot.Logger.Fatal(err)
	}

	bot.Logger.Info("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}