package handler

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/request"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ForgottenPasswordHandler(c echo.Context) error {
	var r request.ForgottenPasswordRequest

	err := h.decode(c.Request().Body, &r)

	if err != nil {
		return response.ServerErrorResponse(err)
	}

	validationErrors := h.validator.Validate(r)

	if len(validationErrors.ValidationErrors) > 0 {
		return response.ValidationErrorResponse(validationErrors)
	}

	serviceErr, err := h.auth.ForgottenPassword(r.Email)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.NoContentResponse(c)
}
