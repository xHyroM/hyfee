package main

import (
	"flag"
	"hyros_coffee/hyfee"
	"hyros_coffee/hyfee/commands"
	"hyros_coffee/hyfee/listeners"
	"os"

	"github.com/disgoorg/snowflake/v2"
	_ "github.com/joho/godotenv/autoload"
)

var (
	botToken = os.Getenv("DISCORD_TOKEN")
	guildId, err =	snowflake.Parse("1046534628577640528")

	syncDatabaseTables *bool
	syncCommands *bool
	debug *bool
)

func init() {
	syncDatabaseTables = flag.Bool("sync-db", false, "Whether to sync the database tables")
	syncCommands = flag.Bool("sync-commands", false, "Whether to sync the commands")
	debug = flag.Bool("debug", false, "Whether to enable debug mode")
	flag.Parse()
}

func main() {
	bot := hyfee.New()

	bot.Handler.AddCommands(
		commands.Auth(bot),
		commands.Experiments(bot),
		commands.Say(bot),
	)

	if err := bot.Setup(
		hyfee.Config{
			SyncDatabaseTables: syncDatabaseTables,
			Debug: debug,
		},
		listeners.Ready(bot),
		listeners.MessageCreate(bot),
	); err != nil {
		bot.Logger.Fatal(err)
	}

	bot.SetupOAuth2()

	if *syncCommands {
		bot.Handler.SyncCommands(bot.Client, guildId)
	}

	bot.Start()
}