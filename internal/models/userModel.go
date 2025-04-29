package models

import (
	"time"
)

type User struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Name        string
	Surname     string
	Age         int
	Gender      string
	Nationality string
}
