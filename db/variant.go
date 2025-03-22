package db

import "gorm.io/gorm"

type Variant struct {
	gorm.Model
	PrintfulID   int64 `gorm:"uniqueIndex"`
	ProductID    uint  // foreign key to Product
	Name         string
	SKU          string
	RetailPrice  string
	ThumbnailURL string
}
