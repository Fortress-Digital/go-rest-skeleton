package cmd

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/handler"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/log"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/route"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/supabase"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/validation"
)

func Execute(log log.LoggerInterface) error {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("Config error", err)
		return err
	}

	auth := supabase.NewAuthClient(cfg.Supabase.Url, cfg.Supabase.Key)
	validator := validation.NewValidator()
	handler := handler.NewHandler(cfg, auth, validator)

	router := route.NewRouter(cfg, handler)

	err = NewServer(cfg, router, log)
	if err != nil {
		log.Error("NewServer error", err)
		return err
	}

	return nil
}
