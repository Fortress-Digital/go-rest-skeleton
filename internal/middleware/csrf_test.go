package middleware

import (
	"fmt"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/go-playground/assert/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		hasHeader       bool
		hasInvalidToken bool
		message         string
	}{
		{
			"With header",
			"token",
			true,
			false,
			"",
		},
		{
			"Without header",
			"token",
			false,
			false,
			"code=400, message=missing csrf token in request header",
		},
		{
			"Invalid token",
			"token",
			true,
			true,
			"code=403, message=invalid csrf token",
		},
	}

	cfg := config.Config{
		Application: config.Application{
			Name: "test-application",
			Env:  "dev",
		},
	}

	for _, tt := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderCookie, fmt.Sprintf("_csrf=%s", tt.token))

		if tt.hasHeader {
			if tt.hasInvalidToken {
				req.Header.Add("X-CSRF-Token", "invalid")
			} else {
				req.Header.Add("X-CSRF-Token", tt.token)
			}
		}

		c := echo.New().NewContext(req, rec)

		csrf := CSRFMiddleware(&cfg)

		h := csrf(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		result := h(c)

		if tt.hasHeader && !tt.hasInvalidToken {
			assert.Equal(t, nil, result)
		} else {
			assert.Equal(t, tt.message, result.Error())
		}
	}
}
