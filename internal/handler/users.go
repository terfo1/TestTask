package handler

import (
	"TestTask/internal/database"
	"TestTask/internal/models"
	"TestTask/pkg/enrich"
	"TestTask/pkg/logger"
	"encoding/json"
	"net/http"
	"strconv"
)

// GetUsers godoc
// @Summary      Получение пользователей
// @Description  Получить список пользователей с фильтрами и пагинацией
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page        query   int     false  "Номер страницы"
// @Param        limit       query   int     false  "Количество на странице"
// @Param        age_min     query   int     false  "Мин. возраст"
// @Param        age_max     query   int     false  "Макс. возраст"
// @Param        gender      query   string  false  "Пол"
// @Param        nationality query   string  false  "Национальность"
// @Success      200  {array}  models.User
// @Failure      400  {string}  string "Invalid request"
// @Router       /user [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Logger.Println("Invalid request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 1
	}
	age_min := r.URL.Query().Get("age_min")
	age_max := r.URL.Query().Get("age_max")
	gender := r.URL.Query().Get("gender")
	nationality := r.URL.Query().Get("nationality")

	offset := (page - 1) * limit
	query := database.DB.Model(&models.User{})

	// Filtering
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

	// Pagination
	var users []models.User
	res := query.Limit(limit).Offset(offset).Find(&users)
	if res.Error != nil {
		logger.Logger.Println("Pagination failed")
		http.Error(w, "Pagination failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// CreateUser godoc
// @Summary      Создание пользователя
// @Description  Добавить нового пользователя и обогатить его данными
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  models.User  true  "User Data"
// @Success      201  {object}  models.User
// @Failure      400  {string}  string "Bad request"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /createuser [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Logger.Println("Invalid request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var body struct {
		Name    string
		Surname string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Logger.Println("Could not parse request body!", err)
		http.Error(w, "Could not parse request body!", http.StatusBadRequest)
		return
	}

	enriched, err := enrich.EnrichData(body.Name)
	if err != nil {
		logger.Logger.Println("Enrichment failed:", err)
		http.Error(w, "Enrichment failed", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:        body.Name,
		Surname:     body.Surname,
		Age:         enriched.Age,
		Gender:      enriched.Gender,
		Nationality: enriched.Nationality,
	}

	res := database.DB.Create(&user)
	if res.Error != nil {
		logger.Logger.Println("Could not create user!")
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	logger.Logger.Println("User created successfully!")
	w.WriteHeader(http.StatusCreated)
}

// DeleteUser godoc
// @Summary      Удаление пользователя
// @Description  Удалить пользователя по ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  query  int  true  "ID пользователя"
// @Success      200  {string}  string "User deleted"
// @Failure      400  {string}  string "Bad request"
// @Failure      404  {string}  string "User not found"
// @Router       /deleteuser [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		logger.Logger.Println("Invalid request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 {
		logger.Logger.Println("Id field is empty")
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
	}

	res := database.DB.Delete(&models.User{}, id)
	if res.Error != nil {
		logger.Logger.Printf("Could not delete user with id %d!", id)
		http.Error(w, "Could not delete user", http.StatusNotFound)
		return
	}
	if res.RowsAffected == 0 {
		logger.Logger.Printf("User with id=%d not found", id)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	logger.Logger.Println("User deleted successfully!")
	w.WriteHeader(http.StatusOK)
}

// UpdateUser godoc
// @Summary      Обновление пользователя
// @Description  Обновить пользователя по ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    query  int         true  "user id"
// @Param        user  body   models.User true  "updated data"
// @Success      200  {object}  models.User
// @Failure      400  {string}  string "Bad request"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /updateuser [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		logger.Logger.Println("Invalid request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 {
		logger.Logger.Println("Id field is empty")
		http.Error(w, "Id field is empty", http.StatusBadRequest)
	}

	var body models.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Logger.Println("Could not parse request body!", err)
		http.Error(w, "Could not parse request body!", http.StatusBadRequest)
		return
	}

	var user models.User
	res := database.DB.First(&user, id)
	if res.Error != nil {
		logger.Logger.Printf("Could not find user with id %d", id)
		http.Error(w, "Could not find user with id", http.StatusBadRequest)
		return
	}

	user.Name = body.Name
	user.Surname = body.Surname
	user.Age = body.Age
	user.Gender = body.Gender
	user.Nationality = body.Nationality

	save := database.DB.Save(&user)
	if save.Error != nil {
		logger.Logger.Printf("Could not update user with id %d", id)
		http.Error(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	logger.Logger.Println("User updated successfully!")
	w.WriteHeader(http.StatusOK)
}
