package handler

import (
	"encoding/json"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/request"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) LoginHandler(c echo.Context) error {
	var r request.LoginRequest

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&r)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	v := validation.NewValidator()
	validationErrors := v.Validate(r)

	if len(validationErrors.ValidationErrors) > 0 {
		return response.ValidationErrorResponse(validationErrors)
	}

	sb := h.AuthClient()
	uc := supabase.UserCredentials{
		Email:    r.Email,
		Password: r.Password,
	}
	user, serviceErr, err := sb.SignIn(uc)
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
