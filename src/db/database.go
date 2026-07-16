package db

import (
	"log/slog"

	"github.com/yazu-codes/scanme-files.git/src/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DSN        string
	Connection *gorm.DB
	logger     *slog.Logger
}

func NewDatabase(dsn string, l *slog.Logger) *Database {
	return &Database{DSN: dsn, Connection: nil, logger: l}
}

func (d *Database) Connect() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: d.DSN,
	}), &gorm.Config{})

	if err != nil {
		d.logger.Error(err.Error())
	}

	d.Connection = db
}

func (d *Database) AutoMigrate() {
	d.logger.Info("Auto-migrating database schema...")
	if err := d.Connection.AutoMigrate(
		&model.Image{},
	); err != nil {
		d.logger.Error(err.Error())
	}
}
