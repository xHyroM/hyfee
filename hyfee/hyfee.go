package hyfee

import (
	"context"
	"hyros_coffee/db"
	"hyros_coffee/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/oauth2"
)


type Bot struct {
	Client bot.Client
	OAuth2Client oauth2.Client
	Logger log.Logger
	Handler *handler.Handler
	Database db.DB
	HTTPServer *http.ServeMux
}

type Config struct {
	SyncDatabaseTables *bool
	Debug *bool
}

func New() *Bot {
	var logger log.Logger = log.New(log.Ldate | log.Ltime | log.Lshortfile)

	return &Bot{
		Logger: logger,
		Handler: handler.New(logger),
	}
}

func (b *Bot) Setup(config Config, listeners ...bot.EventListener) (err error) {
	b.HTTPServer = http.NewServeMux()
	b.HTTPServer.HandleFunc("/callback", b.OAuthHandler)

	b.Client, err = disgo.New(os.Getenv("DISCORD_TOKEN"),
		bot.WithLogger(log.New(log.Ldate | log.Ltime | log.Lshortfile)),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds)),
		bot.WithEventListeners(append([]bot.EventListener{b.Handler}, listeners...)...),
	)

	if b.Database, err = db.Setup(config.SyncDatabaseTables, config.Debug); err != nil {
		b.Logger.Error("Failed to setup database", err)
	}

	return err
}

func (b *Bot) Start() {
	if err := b.Client.OpenGateway(context.TODO()); err != nil {
		b.Logger.Error("Failed to open gateway", err)
	}

	if err := http.ListenAndServe(":8080", b.HTTPServer); err != nil {
		b.Logger.Error("Failed to listen and serve", err)
	}

	defer func() {
		b.Client.Close(context.TODO())
		b.Database.Close()
	}()

	b.Logger.Info("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}