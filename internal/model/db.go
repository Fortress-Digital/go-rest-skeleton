package model

import (
	"github.com/Fortress-Digital/go-rest-skeleton/internal/config"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
}

func ormConnection(cfg *config.Config, log log.LoggerInterface) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.Dsn), &gorm.Config{})
	if err != nil {
		log.Info("DB configuration error", err)
		return nil, err
	}

	return db, nil
}

func NewDB(cfg *config.Config, log log.LoggerInterface) (*gorm.DB, error) {
	db, err := ormConnection(cfg, log)
	if err != nil {
		return nil, err
	}

	// db.AutoMigrate(&Test{})

	return db, nil
}
