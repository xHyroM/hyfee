package listeners

import (
	"hyros_coffee/hyfee"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func MessageCreate(b *hyfee.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(event bot.Event) {
		switch event := event.(type) {
			case *events.MessageCreate:
				// Crosspost messages in datamining category
				channel, ok := event.Channel()
				if !ok {
					b.Logger.Error("Failed to get channel")
					return
				}

				if *channel.ParentID() != snowflake.MustParse("1111909235756896286") {
					return
				}

				_, err := b.Client.Rest().CrosspostMessage(event.ChannelID, event.MessageID)
				if err != nil {
					b.Logger.Error("Failed to crosspost message", err)
				}

				b.Logger.Info("Message " + event.Message.ID.String() + " crossposted")
		}
	})
}