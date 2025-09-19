package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex"`
	Token         string
	PersonID      int
	InitializedAt time.Time `gorm:"not null"`
}
