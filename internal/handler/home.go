package handler

import (
	"fmt"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HomeHandler(c echo.Context) error {
	r := map[string]string{
		"message": fmt.Sprintf("Welcome to %s.", h.cfg.Application.Name),
	}
	return response.SuccessResponse(c, r)
}
