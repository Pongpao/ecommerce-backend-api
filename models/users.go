package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"unique"`
	Email     string    `gorm:"unique"`
	Password  string
	Role      string
	CreatedAt time.Time
}

