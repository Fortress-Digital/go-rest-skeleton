package handler

import (
	"encoding/json"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/request"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterHandler(c echo.Context) error {
	var r request.RegisterRequest

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&r)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	sb := h.App.AuthClient()
	uc := supabase.UserCredentials{
		Email:    r.Email,
		Password: r.Password,
	}

	user, serviceErr, err := sb.SignUp(uc)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.CreatedResponse(c, user)
}
