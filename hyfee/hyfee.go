package hyfee

import (
	"context"
	"hyros_coffee/db"
	"hyros_coffee/handler"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"
	"github.com/redis/go-redis/v9"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/oauth2"
)


type Bot struct {
	Client bot.Client
	OAuth2Client oauth2.Client
	Logger log.Logger
	Handler *handler.Handler
	Database db.DB
	HTTPHandler *http.ServeMux
	HTTPServer *http.Server
	RedisClient *redis.Client
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
	b.HTTPHandler = http.NewServeMux()
	b.HTTPHandler.HandleFunc("/callback", b.OAuthHandler)
	b.HTTPHandler.HandleFunc("/linked-roles", b.LinkedRolesHandler)

	b.RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	b.Client, err = disgo.New(os.Getenv("DISCORD_TOKEN"),
		bot.WithLogger(slog.Default()),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentGuildPresences)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds), cache.WithCaches(cache.FlagChannels)),
		bot.WithEventListeners(append([]bot.EventListener{b.Handler}, listeners...)...),
	)

	if b.Database, err = db.Setup(config.SyncDatabaseTables, config.Debug); err != nil {
		b.Logger.Error("Failed to setup database", err)
	}

	return err
}

func (b *Bot) SetupOAuth2() {
	b.OAuth2Client = oauth2.New(b.Client.ApplicationID(), os.Getenv("CLIENT_SECRET"))
}

func (b *Bot) Start() {
	if err := b.Client.OpenGateway(context.TODO()); err != nil {
		b.Logger.Error("Failed to open gateway", err)
	}

	b.HTTPServer = &http.Server{Addr: ":"+os.Getenv("MUX_SERVER_HTTP_PORT"), Handler: b.HTTPHandler}
	if err := b.HTTPServer.ListenAndServe(); err != nil {
		b.Logger.Error("Failed to listen and serve", err)
	}

	defer func() {
		b.Client.Close(context.TODO())
		b.Database.Close()
		b.HTTPServer.Shutdown(context.TODO())
	}()

	b.Logger.Info("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}