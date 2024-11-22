package cmd

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/model"
)

type Register interface {
	Register(app *config.App) error
}

func Execute(app *config.App) error {
	err := register(
		app,
		config.Configuration{},
		model.DB{},
	)

	if err != nil {
		app.Logger.Info("Registration error", err)
		return err
	}

	err = Server(app)
	if err != nil {
		app.Logger.Info("Server error", err)
		return err
	}

	return nil
}

func register(app *config.App, registers ...Register) error {
	for _, register := range registers {
		err := register.Register(app)
		if err != nil {
			return err
		}
	}

	return nil
}
