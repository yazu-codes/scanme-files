package model

import "time"

type Image struct {
	ID         int64  `json:"id" gorm:"primaryKey"`
	UUID       string `json:"uuid" gorm:"uniqueIndex;not null"`
	StorageKey string `json:"storage_key" gorm:"not null"`
	MenuID     string `json:"menu_id" gorm:"column: menu_id"`

	CreatedAt time.Time `json:"created_at"`
}

type Images []Image

func (i *Images) ToDTOs() []ImageDTO {
	imageDTOs := []ImageDTO{}
	for _, im := range *i {
		imageDTOs = append(imageDTOs, ImageDTO{
			MenuID: im.MenuID,
			Url:    im.StorageKey,
		})
	}
	return imageDTOs
}

type ImageDTO struct {
	MenuID string `json:"menu_id"`
	Url    string `json: "url"`
}
