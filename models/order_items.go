package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	OrderID   uuid.UUID `gorm:"type:uuid;index;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;index;not null"`

	Quantity int     `gorm:"not null;check:quantity > 0"`
	Price    float64 `gorm:"type:numeric(12,2);not null"`

	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}