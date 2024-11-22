package route

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/handler"
	middlewares "github.com/Fortress-Digital/go-rest-skeleton/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func NewRouter(app *config.App) http.Handler {
	router := echo.New()
	router.Use(middleware.Recover())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	router.Use(middlewares.CSRFMiddleware(app.Config))

	defineRoutes(router, app)

	return router
}

func defineRoutes(router *echo.Echo, app *config.App) {
	h := handler.Handler{App: app}

	router.GET("/", h.HomeHandler)
	router.POST("/register", h.RegisterHandler)
	router.POST("/login", h.LoginHandler)
	router.POST("/forgotten-password", h.ForgottenPasswordHandler)
	router.POST("/reset-password", h.ResetPasswordHandler)
	router.POST("/refresh-token", h.RefreshTokenHandler)
	router.POST("/logout", h.LogoutHandler)
}
