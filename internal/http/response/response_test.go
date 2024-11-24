package response

import (
	"errors"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/validation"
	"github.com/go-playground/assert/v2"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	result := ErrorResponse(http.StatusInternalServerError, "error")
	expected := echo.HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "error",
	}

	assert.Equal(t, result, expected)
}

func TestServerErrorResponse(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{"No message", ""},
		{"With message", "Test error!"},
	}

	for _, tt := range tests {
		if tt.message == "" {
			result := ServerErrorResponse()
			expected := echo.HTTPError{
				Code: http.StatusInternalServerError,
				Message: Error{
					Message: "the server encountered a problem and could not process your request",
				},
			}

			assert.Equal(t, result, expected)
		} else {
			err := errors.New(tt.message)
			result := ServerErrorResponse(err)
			expected := echo.HTTPError{
				Code: http.StatusInternalServerError,
				Message: Error{
					Message: "Test error!",
				},
			}

			assert.Equal(t, result, expected)
		}
	}
}

func TestBadRequestResponse(t *testing.T) {

	tests := []struct {
		name    string
		message string
		asError bool
	}{
		{"No error", "Normal message!", false},
		{"With error", "Test error!", true},
	}

	for _, tt := range tests {
		var err any

		if tt.asError {
			err = errors.New(tt.message)
		} else {
			err = tt.message
		}

		result := BadRequestResponse(err)
		expected := echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: tt.message,
		}

		assert.Equal(t, result, expected)
	}
}

func TestValidationErrorResponse(t *testing.T) {
	errs := validation.ValidationErrors{
		Message: "Validation error",
	}
	err := validation.ValidationError{
		Message: "Email is a required field",
		Field:   "email",
	}
	errs.ValidationErrors = append(errs.ValidationErrors, err)

	result := ValidationErrorResponse(errs)
	expected := echo.HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: errs,
	}

	assert.Equal(t, result, expected)
}

func TestNoContentResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	NoContentResponse(c)

	assert.Equal(t, rec.Code, http.StatusNoContent)
}

func TestSuccessResponse(t *testing.T) {
	body := map[string]string{"test": "foo"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	_ = SuccessResponse(c, body)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), "{\"test\":\"foo\"}\n")
}

func TestCreatedResponse(t *testing.T) {
	body := map[string]string{"test": "foo"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	_ = CreatedResponse(c, body)

	assert.Equal(t, rec.Code, http.StatusCreated)
	assert.Equal(t, rec.Body.String(), "{\"test\":\"foo\"}\n")
}

func TestUnauthorizedResponse(t *testing.T) {
	body := map[string]string{"test": "foo"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	_ = UnauthorizedResponse(c, body)

	assert.Equal(t, rec.Code, http.StatusUnauthorized)
	assert.Equal(t, rec.Body.String(), "{\"test\":\"foo\"}\n")
}

func TestResponseWithBadDate(t *testing.T) {
	body := math.NaN()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	result := Response(c, http.StatusOK, body)
	expected := echo.HTTPError{
		Code: http.StatusInternalServerError,
		Message: Error{
			Message: "json: unsupported value: NaN",
		},
	}

	assert.Equal(t, result, expected)
}
