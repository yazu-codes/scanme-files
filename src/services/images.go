package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/yazu-codes/scanme-files.git/src/model"
	"github.com/yazu-codes/scanme-files.git/src/repositories"
	"github.com/yazu-codes/scanme-files.git/src/storage"
)

const (
	maxWidth    = 300
	maxHeight   = 400
	jpegQuality = 82 // 80-85 is the sweet spot: near-invisible loss, big size drop
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

// fitWithinBounds resizes src to fit within maxW x maxH, preserving aspect
// ratio, and never upscales an already-smaller image.
func fitWithinBounds(src image.Image, maxW, maxH int) image.Image {
	b := src.Bounds()
	if b.Dx() <= maxW && b.Dy() <= maxH {
		return src
	}
	return imaging.Fit(src, maxW, maxH, imaging.Lanczos)
}

func (s *Images) Create(
	ctx context.Context,
	file io.Reader,
	filename string,
	menuId string,
) (*model.Image, error) {
	// Decode whatever format was uploaded (jpeg/png/webp/etc — imaging
	// auto-orients based on EXIF too, which raw image.Decode won't do).
	src, err := imaging.Decode(file, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	// Fit within 300x400, preserving aspect ratio, only downscaling
	// (never upscaling a smaller source image).
	resized := fitWithinBounds(src, maxWidth, maxHeight)

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, resized, imaging.JPEG, imaging.JPEGQuality(jpegQuality)); err != nil {
		return nil, fmt.Errorf("encode image: %w", err)
	}

	id := uuid.New().String()

	key := fmt.Sprintf(
		"/images/%s%s",
		id,
		filepath.Ext(filename),
	)

	err = s.storage.Save(ctx, key, file)
	if err != nil {
		return nil, err
	}

	image := model.Image{
		UUID:       id,
		StorageKey: key,
		MenuID:     menuId,
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
