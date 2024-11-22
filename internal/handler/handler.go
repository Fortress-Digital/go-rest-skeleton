package handler

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
)

type Handler struct {
	App *config.App
}

func (h *Handler) AuthClient() *supabase.Auth {
	return h.App.AuthClient()
}
