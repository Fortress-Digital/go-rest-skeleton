package config

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"gorm.io/gorm"
	"log/slog"
)

type App struct {
	Logger *slog.Logger
	Config *Config
	DB     *gorm.DB
}

func (a *App) AuthClient() *supabase.Auth {
	return supabase.CreateAuth(a.Config.Supabase.Url, a.Config.Supabase.Key)
}
