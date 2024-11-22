package handler

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) LogoutHandler(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	token = token[7:]

	sb := h.NewSupabaseClient()
	serviceErr, err := sb.Auth.SignOut(token)

	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.NoContentResponse(c)
}
