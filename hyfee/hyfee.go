package hyfee

import (
	"context"
	"hyros_coffe/handler"
	"os"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)


type Bot struct {
	Client bot.Client
	Logger log.Logger
	Handler *handler.Handler
}

func New() *Bot {
	var logger log.Logger = log.New(log.Ldate | log.Ltime | log.Lshortfile)

	return &Bot{
		Logger: logger,
		Handler: handler.New(logger),
	}
}

func (b *Bot) Setup(listeners ...bot.EventListener) (err error) {
	b.Client, err = disgo.New(os.Getenv("DISCORD_TOKEN"),
		bot.WithLogger(log.New(log.Ldate | log.Ltime | log.Lshortfile)),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds)),
		bot.WithEventListeners(append([]bot.EventListener{b.Handler}, listeners...)...),
	)

	return err
}

func (b *Bot) Start() error {
	return b.Client.OpenGateway(context.TODO())
}