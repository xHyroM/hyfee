package handler

import (
	"strings"

	"github.com/disgoorg/disgo/events"
)

type ModalHandler func(event *events.ModalSubmitInteractionCreate) error

type Modal struct {
	Name    string
	Check   Check[*events.ModalSubmitInteractionCreate]
	Handler ModalHandler
}

func (h *Handler) handleModal(event *events.ModalSubmitInteractionCreate) {
	customID := event.Data.CustomID
	if !strings.HasPrefix(customID, "handler:") {
		return
	}

	modalName := strings.Split(customID, ":")[1]
	modal, ok := h.Modals[modalName]
	if !ok || modal.Handler == nil {
		h.Logger.Errorf("No modal handler for \"%s\" found", modalName)
	}

	if modal.Check != nil && !modal.Check(event) {
		return
	}

	if err := modal.Handler(event); err != nil {
		h.Logger.Errorf("Failed to handle modal interaction for \"%s\" : %s", modalName, err)
	}
}