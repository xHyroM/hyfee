package handler

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var _ bot.EventListener = (*Handler)(nil)

func New(logger log.Logger) *Handler {
	return &Handler{
		Logger:     logger,
		Commands:   map[string]Command{},
		Components: map[string]Component{},
		Modals:     map[string]Modal{},
	}
}

type Handler struct {
	Logger log.Logger

	Commands   map[string]Command
	Components map[string]Component
	Modals     map[string]Modal
}

func (h *Handler) AddCommands(commands ...Command) {
	for _, command := range commands {
		h.Commands[command.Create.CommandName()] = command
	}
}

func (h *Handler) AddComponents(components ...Component) {
	for _, component := range components {
		h.Components[component.Name] = component
	}
}

func (h *Handler) AddModals(modals ...Modal) {
	for _, modal := range modals {
		h.Modals[modal.Name] = modal
	}
}

func (h *Handler) SyncCommands(client bot.Client, guildIds ...snowflake.ID) {
	commands := make([]discord.ApplicationCommandCreate, len(h.Commands))
	var i int
	for _, command := range h.Commands {
		commands[i] = command.Create
		i++
	}

	if len(guildIds) == 0 {
		if _, err := client.Rest().SetGlobalCommands(client.ApplicationID(), commands); err != nil {
			h.Logger.Error("Failed to sync global commands: ", err)
			return
		}
		h.Logger.Infof("Synced %d global commands", len(commands))
		return
	}

	for _, guildId := range guildIds {
		if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), guildId, commands); err != nil {
			h.Logger.Errorf("Failed to sync commands for guild %d: %s", guildId, err)
			continue
		}
		h.Logger.Infof("Synced %d commands for guild %s", len(commands), guildId)
	}
}

func (h *Handler) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.ApplicationCommandInteractionCreate:
		h.handleCommand(e)
	case *events.AutocompleteInteractionCreate:
		h.handleAutocomplete(e)
	case *events.ComponentInteractionCreate:
		h.handleComponent(e)
	case *events.ModalSubmitInteractionCreate:
		h.handleModal(e)
	}
}