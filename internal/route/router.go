package route

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/handler"
	middlewares "github.com/Fortress-Digital/go-rest-skeleton/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func NewRouter(cfg *config.Config, handler *handler.Handler) http.Handler {
	router := echo.New()
	router.Use(middleware.Recover())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	router.Use(middlewares.CSRFMiddleware(cfg))

	defineRoutes(router, handler)

	return router
}

func defineRoutes(router *echo.Echo, h *handler.Handler) {
	router.GET("/", h.HomeHandler)
	router.POST("/register", h.RegisterHandler)
	router.POST("/login", h.LoginHandler)
	router.POST("/forgotten-password", h.ForgottenPasswordHandler)
	router.POST("/reset-password", h.ResetPasswordHandler)
	router.POST("/refresh-token", h.RefreshTokenHandler)
	router.POST("/logout", h.LogoutHandler)
}
