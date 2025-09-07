package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex"`
	Token         string
	PersonId      int
	InitializedAt time.Time `gorm:"not null"`
}
