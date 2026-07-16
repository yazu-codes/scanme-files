package model

import "time"

type Image struct {
	ID         int64  `json:"id" gorm:"primaryKey"`
	UUID       string `json:"uuid" gorm:"uniqueIndex;not null"`
	StorageKey string `json:"storage_key" gorm:"not null"`

	CreatedAt time.Time `json:"created_at"`
}
