package repository

import (
	"TestTask/internal/database"
	"TestTask/internal/models"
	"gorm.io/gorm"
)

type UserFilter struct {
	Gender      string
	Nationality string
	AgeMin      *int
	AgeMax      *int
}

func GetById(user *models.User, id int) *gorm.DB {
	res := database.DB.First(user, id)
	return res
}

func SaveInDb(user *models.User) *gorm.DB {
	res := database.DB.Save(user)
	return res
}

func DeleteInDb(user *models.User, id int) *gorm.DB {
	res := database.DB.Delete(user, id)
	return res
}

func CreateInDb(user *models.User) *gorm.DB {
	res := database.DB.Create(&user)
	return res
}

func GetByParams(filter UserFilter, page, limit int) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * limit

	query := database.DB.Model(&models.User{})

	if filter.Gender != "" {
		query = query.Where("gender = ?", filter.Gender)
	}
	if filter.Nationality != "" {
		query = query.Where("nationality = ?", filter.Nationality)
	}
	if filter.AgeMin != nil {
		query = query.Where("age >= ?", *filter.AgeMin)
	}
	if filter.AgeMax != nil {
		query = query.Where("age <= ?", *filter.AgeMax)
	}

	res := query.Limit(limit).Offset(offset).Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}
