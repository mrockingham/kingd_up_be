package db

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	PrintfulID   int64  `gorm:"uniqueIndex"`
	ExternalID   string `gorm:"index"`
	Name         string
	ThumbnailURL string
	Synced       bool
	Variants     []Variant `gorm:"foreignKey:ProductID"`
}
