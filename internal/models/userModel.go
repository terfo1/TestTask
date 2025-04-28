package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string
	Surname     string
	Age         int
	Gender      string
	Nationality string
}
