package handler

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/request"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/http/response"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) RegisterHandler(c echo.Context) error {
	var r request.RegisterRequest

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

	user, serviceErr, err := h.auth.SignUp(uc)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.CreatedResponse(c, user)
}

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

func (h *Handler) LogoutHandler(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	token = token[7:]

	serviceErr, err := h.auth.SignOut(token)

	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.NoContentResponse(c)
}

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

func (h *Handler) ResetPasswordHandler(c echo.Context) error {
	var r request.ResetPasswordRequest

	err := h.decode(c.Request().Body, &r)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	validationErrors := h.validator.Validate(r)

	if len(validationErrors.ValidationErrors) > 0 {
		return response.ValidationErrorResponse(validationErrors)
	}

	token := c.Request().Header.Get("Authorization")
	token = token[7:]

	serviceErr, err := h.auth.ResetPassword(token, r.Password)
	if err != nil {
		return response.ServerErrorResponse(err)
	}

	if serviceErr != nil {
		return response.BadRequestResponse(serviceErr)
	}

	return response.NoContentResponse(c)
}
