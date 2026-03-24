
package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	PaymentPending = "pending"
	PaymentSuccess = "success"
	PaymentFailed  = "failed"
)

type Payment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Amount    float64   `gorm:"type:numeric(12,2);not null"`
	Status    string    `gorm:"type:varchar(20)"` // pending, success, failed
	Method    string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time

	Order     Order     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}