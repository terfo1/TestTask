package handler

import (
	"TestTask/internal/database"
	"TestTask/internal/models"
	"TestTask/pkg/logger"
	"encoding/json"
	"net/http"
	"strconv"
)

type Pagination struct {
	Next          int
	Previous      int
	RecordPerPage int
	CurrentPage   int
	TotalPage     int
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Logger.Println("Invalid request")
		http.Error(w, "error", http.StatusBadRequest)
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 0 {
		limit = 1
	}
	age_min := r.URL.Query().Get("age_min")
	age_max := r.URL.Query().Get("age_max")
	gender := r.URL.Query().Get("gender")
	nationality := r.URL.Query().Get("nationality")

	offset := (page - 1) * limit
	query := database.DB.Model(&models.User{})

	//Filtering
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if nationality != "" {
		query = query.Where("nationality = ?", nationality)
	}
	if age_min != "" {
		ageMin, _ := strconv.Atoi(age_min)
		query = query.Where("age >= ?", ageMin)
	}
	if age_max != "" {
		ageMax, _ := strconv.Atoi(age_max)
		query = query.Where("age <= ?", ageMax)
	}

	//Pagination
	var users []models.User
	res := query.Limit(limit).Offset(offset).Find(&users)
	if res.Error != nil {
		logger.Logger.Println("Pagination failed")
		http.Error(w, "error", http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusCreated)
}
