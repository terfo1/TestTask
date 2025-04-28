package database

import "TestTask/internal/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
