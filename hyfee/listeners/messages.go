package listeners

import (
	"hyros_coffee/hyfee"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func MessageCreate(b *hyfee.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event bot.Event) {
		switch event.(type) {
			case *events.MessageCreate:
				event := event.(*events.MessageCreate)
				channel, ok := event.Channel()
				if ok != true {
					b.Logger.Error("Failed to get channel")
					return
				}

				channel, ok = channel.(discord.GuildNewsChannel)
				if ok != true {
					b.Logger.Error("Failed to get channel as news")
					return
				}

				if *channel.ParentID() != snowflake.MustParse("1111909235756896286") {
					return
				}

				_, err := b.Client.Rest().CrosspostMessage(event.ChannelID, event.MessageID)
				if err != nil {
					b.Logger.Error("Failed to crosspost message", err)
				}

				b.Logger.Info("Crossposted message")
		}
	})
}