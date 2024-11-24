package response

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Error struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func ErrorResponse(status int, message any) *echo.HTTPError {
	return echo.NewHTTPError(status, message)
}

func ServerErrorResponse(errors ...error) *echo.HTTPError {
	if len(errors) > 0 {
		return ErrorResponse(http.StatusInternalServerError, Error{
			Message: errors[0].Error(),
		})
	}

	return ErrorResponse(http.StatusInternalServerError, Error{
		Message: "the server encountered a problem and could not process your request",
	})
}

func BadRequestResponse(err any) *echo.HTTPError {
	if err, ok := err.(error); ok {
		return ErrorResponse(http.StatusBadRequest, err.Error())
	}

	return ErrorResponse(http.StatusBadRequest, err)
}

func ValidationErrorResponse(err validation.ValidationErrors) *echo.HTTPError {
	return ErrorResponse(http.StatusUnprocessableEntity, err)
}

func NoContentResponse(c echo.Context) error {
	_ = c.NoContent(http.StatusNoContent)

	return nil
}

func Response(c echo.Context, status int, data any) *echo.HTTPError {
	err := c.JSON(status, data)
	if err != nil {
		return ServerErrorResponse(err)
	}

	return nil
}

func SuccessResponse(c echo.Context, data any) *echo.HTTPError {
	return Response(c, http.StatusOK, data)
}

func CreatedResponse(c echo.Context, data any) *echo.HTTPError {
	return Response(c, http.StatusCreated, data)
}

func UnauthorizedResponse(c echo.Context, data any) *echo.HTTPError {
	return Response(c, http.StatusUnauthorized, data)
}
