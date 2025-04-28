package handler

import (
	"TestTask/internal/database"
	"TestTask/internal/models"
	"TestTask/pkg/logger"
	"encoding/json"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Logger.Println("Invalid request")
		http.Error(w, "error", http.StatusBadRequest)
	}
	
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Logger.Println("Invalid request")
		http.Error(w, "error", http.StatusBadRequest)
	}
	var body struct {
		Name    string
		Surname string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Logger.Println("Could not parse request body!", err)
		http.Error(w, "error", http.StatusBadRequest)
	}

	user := models.User{Name: body.Name, Surname: body.Surname}
	res := database.DB.Create(&user)
	if res.Error != nil {
		logger.Logger.Println("Could not create user!")
		http.Error(w, "error", http.StatusInternalServerError)
	}

	logger.Logger.Println("User created successfully!")
}
