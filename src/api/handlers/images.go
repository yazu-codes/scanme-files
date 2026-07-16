package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/yazu-codes/scanme-files.git/src/services"
)

type Images struct {
	logger  *slog.Logger
	service *services.Images
}

func NewImages(service *services.Images, logger *slog.Logger) *Images {
	return &Images{service: service, logger: logger}
}

func (i *Images) AddImage(c *gin.Context) {
	i.logger.Info("ADDING IMAGE")
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "missing image"})
		return
	}
	defer file.Close()

	i.logger.Info("FOUND IMAGE")
	i.logger.Info("ADDING METADATA FOR IMAGE")
	created, err := i.service.Create(
		c.Request.Context(),
		file,
		header.Filename,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, created)
}

func (i *Images) AddImages(c *gin.Context) {
	i.logger.Info("Add Images Handler reached.")
}

func (i *Images) GetImages(c *gin.Context) {
	c.JSON(200, "Hiii")
	i.logger.Info("Get Images Handler reached.")
}

func (i *Images) GetImageById(c *gin.Context) {
	i.logger.Info("Get Image by Id Handler reached.")
}

func (i *Images) DeleteImageById(c *gin.Context) {
	i.logger.Info("Delete Image by Id Handler reached.")
}
