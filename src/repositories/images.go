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

func (i *Images) GetByMenuId(menuId string) (model.Images, error) {
	var images []model.Image
	if err := i.db.Connection.Where("menu_id = ?", menuId).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}
