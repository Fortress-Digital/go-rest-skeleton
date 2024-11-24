package middleware

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CSRFMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieSecure:   cfg.Application.Env == "production",
		CookieHTTPOnly: cfg.Application.Env == "production",
		TokenLookup:    "header:X-CSRF-Token",
	})
}
