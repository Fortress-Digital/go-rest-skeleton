package handler

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/request"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) LoginHandler(c echo.Context) error {
	var r request.LoginRequest

	err := h.decode(c.Request().Body, &r)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	validationErrors := h.validator.Validate(r)

	if len(validationErrors.ValidationErrors) > 0 {
		return response.ValidationErrorResponse(validationErrors)
	}

	uc := supabase.UserCredentials{
		Email:    r.Email,
		Password: r.Password,
	}
	user, serviceErr, err := h.auth.SignIn(uc)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		if serviceErr.Code == http.StatusUnauthorized {
			return response.UnauthorizedResponse(c, serviceErr)
		}

		return response.BadRequestResponse(serviceErr)
	}

	return response.SuccessResponse(c, user)
}
