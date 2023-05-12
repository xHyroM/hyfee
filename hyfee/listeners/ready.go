package listeners

import (
	"hyros_coffe/hyfee"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

func Ready(b *hyfee.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event bot.Event) {
		switch event.(type) {
			case *events.Ready:
				b.Logger.Info("Bot is ready!")
		}
	})
}