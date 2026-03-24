package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string    `gorm:"not null"`
	Description   string
	Price         float64 `gorm:"not null"`
	Stock         int     `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
	AverageRating float64
	Reviews       []Review `gorm:"foreignKey:ProductID"`
}
