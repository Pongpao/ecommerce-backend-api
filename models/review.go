package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_product;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_product;not null"`
	Rating    int       `gorm:"type:int;check:rating >= 1 AND rating <= 5;not null"`
	Comment   string    `gorm:"type:text;size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
	User      User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}