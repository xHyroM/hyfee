package handler

import (
	"strings"

	"github.com/disgoorg/disgo/events"
)

type ComponentHandler func(event *events.ComponentInteractionCreate) error

type Component struct {
	Name    string
	Check   Check[*events.ComponentInteractionCreate]
	Handler ComponentHandler
}

func (h *Handler) handleComponent(event *events.ComponentInteractionCreate) {
	customID := event.Data.CustomID()
	if !strings.HasPrefix(customID, "handler:") {
		return
	}

	componentName := strings.Split(customID, ":")[1]
	component, ok := h.Components[componentName]
	if !ok || component.Handler == nil {
		h.Logger.Errorf("No component handler for \"%s\" found", componentName)
	}

	if component.Check != nil && !component.Check(event) {
		return
	}

	if err := component.Handler(event); err != nil {
		h.Logger.Errorf("Failed to handle component interaction for \"%s\" : %s", componentName, err)
	}
}