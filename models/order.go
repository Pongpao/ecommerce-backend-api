package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Total     float64   `gorm:"type:numeric(12,2);not null"`
	Status    string    `gorm:"type:varchar(20);default:'pending';index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Payments   []Payment   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderStatusHistory struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Order      Order     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	FromStatus string    `gorm:"type:varchar(20);not null"`
	ToStatus   string    `gorm:"type:varchar(20);not null"`
	ChangedBy  uuid.UUID `gorm:"type:uuid;not null"` // admin หรือ user ที่เปลี่ยน
	CreatedAt  time.Time
}

const (
	OrderPending    = "pending"
	OrderPaid       = "paid"
	OrderProcessing = "processing"
	OrderShipped    = "shipped"
	OrderCompleted  = "completed"
	OrderCancelled  = "cancelled"
	OrderRefunded   = "refunded"
)