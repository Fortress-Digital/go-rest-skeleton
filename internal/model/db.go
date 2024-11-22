package model

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
}

func (d DB) Register(app *config.App) error {
	db, err := newModels(app)
	if err != nil {
		return err
	}

	app.DB = db

	return nil
}

func ormConnection(app *config.App) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(app.Config.Database.Dsn), &gorm.Config{})
	if err != nil {
		app.Logger.Info("DB configuration error", err)
		return nil, err
	}

	return db, nil
}

func newModels(app *config.App) (*gorm.DB, error) {
	db, err := ormConnection(app)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Test{})

	return db, nil
}
