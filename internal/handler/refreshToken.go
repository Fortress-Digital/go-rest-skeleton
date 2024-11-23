package handler

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/request"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RefreshTokenHandler(c echo.Context) error {
	var r request.RefreshTokenRequest

	err := h.decode(c.Request().Body, &r)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	validationErrors := h.validator.Validate(r)

	if len(validationErrors.ValidationErrors) > 0 {
		return response.ValidationErrorResponse(validationErrors)
	}

	user, serviceErr, err := h.auth.RefreshToken(r.RefreshToken)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.SuccessResponse(c, user)
}
