package handler

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
)

type Handler struct {
	App *config.App
}

func (h *Handler) NewSupabaseClient() *supabase.Client {
	return h.App.NewSupabaseClient()
}
