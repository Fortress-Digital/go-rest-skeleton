package handler

import (
	"encoding/json"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/validation"
	"io"
)

type Handler struct {
	cfg       *config.Config
	auth      supabase.AuthClientInterface
	validator validation.ValidatorInterface
}

func NewHandler(cfg *config.Config, auth supabase.AuthClientInterface, validator validation.ValidatorInterface) *Handler {
	return &Handler{
		cfg:       cfg,
		auth:      auth,
		validator: validator,
	}
}

func (h *Handler) decode(data io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(data)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}

	return nil
}
