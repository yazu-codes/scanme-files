package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/yazu-codes/scanme-files.git/src/model"
	"github.com/yazu-codes/scanme-files.git/src/repositories"
	"github.com/yazu-codes/scanme-files.git/src/storage"
)

type Images struct {
	logger     *slog.Logger
	repository *repositories.Images
	storage    storage.Storage
}

func NewImages(r *repositories.Images, l *slog.Logger, s storage.Storage) *Images {
	return &Images{
		logger:     l,
		repository: r,
		storage:    s,
	}
}

func (s *Images) Create(
	ctx context.Context,
	file io.Reader,
	filename string,
) (*model.Image, error) {

	id := uuid.New().String()

	key := fmt.Sprintf(
		"/images/%s%s",
		id,
		filepath.Ext(filename),
	)

	err := s.storage.Save(ctx, key, file)
	if err != nil {
		return nil, err
	}

	image := model.Image{
		UUID:       id,
		StorageKey: key,
	}

	if err := s.repository.Create(&image); err != nil {
		return nil, err
	}

	return &image, nil
}

func (s *Images) GetByMenuId(menuId string) ([]model.ImageDTO, error) {
	images, err := s.repository.GetByMenuId(menuId)
	if err != nil {
		return nil, err
	}

	return images.ToDTOs(), nil
}
