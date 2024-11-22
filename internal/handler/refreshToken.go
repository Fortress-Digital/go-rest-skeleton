package handler

import (
	"encoding/json"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/request"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/validation"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RefreshTokenHandler(c echo.Context) error {
	var r request.RefreshTokenRequest

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

	sb := h.NewSupabaseClient()
	user, serviceErr, err := sb.Auth.RefreshToken(r.RefreshToken)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.SuccessResponse(c, user)
}