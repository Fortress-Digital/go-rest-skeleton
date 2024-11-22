package main

import (
	"github.com/Fortress-Digital/go-rest-skeleton/cmd"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := config.App{
		Logger: logger,
	}

	err := cmd.Execute(&app)
	logger.Error(err.Error())
	os.Exit(1)
}
