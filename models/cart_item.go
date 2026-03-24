package models

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CartID    uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null;check:quantity > 0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Cart      Cart    `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
