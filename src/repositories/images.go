package repositories

import (
	"log/slog"

	"github.com/yazu-codes/scanme-files.git/src/db"
	"github.com/yazu-codes/scanme-files.git/src/model"
)

type Images struct {
	logger *slog.Logger
	db     *db.Database
}

func NewImages(d *db.Database, l *slog.Logger) *Images {
	return &Images{logger: l, db: d}
}

func (i *Images) Create(image *model.Image) error {
	if err := i.db.Connection.Create(image).Error; err != nil {
		return err
	}
	return nil
}
