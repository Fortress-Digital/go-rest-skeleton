package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/validation"
	"github.com/labstack/echo/v4"
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

func (h *Handler) HomeHandler(c echo.Context) error {
	r := map[string]string{
		"message": fmt.Sprintf("Welcome to %s.", h.cfg.Application.Name),
	}
	return response.SuccessResponse(c, r)
}
