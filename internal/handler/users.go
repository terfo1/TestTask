package handler

import (
	"TestTask/internal/models"
	"TestTask/internal/repository"
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
		logger.Logger.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	filters := repository.UserFilter{
		Gender:      r.URL.Query().Get("gender"),
		Nationality: r.URL.Query().Get("nationality"),
	}
	if ageMin := r.URL.Query().Get("age_min"); ageMin != "" {
		if val, err := strconv.Atoi(ageMin); err == nil {
			filters.AgeMin = &val
		}
	}
	if ageMax := r.URL.Query().Get("age_max"); ageMax != "" {
		if val, err := strconv.Atoi(ageMax); err == nil {
			filters.AgeMax = &val
		}
	}

	users, err := repository.GetByParams(filters, page, limit)
	if err != nil {
		logger.Logger.Printf("Error retrieving users: %v", err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
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
	if body.Name == "" || body.Surname == "" {
		logger.Logger.Println("Name and surname are required!")
		http.Error(w, "Name and surname are required", http.StatusBadRequest)
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

	res := repository.CreateInDb(&user)
	if res.Error != nil {
		logger.Logger.Println("Could not create user!")
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	logger.Logger.Println("User created successfully!")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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

	res := repository.DeleteInDb(&models.User{}, id)
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
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
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
	if body.Name == "" || body.Surname == "" {
		logger.Logger.Println("Name and surname are required!")
		http.Error(w, "Name and surname are required", http.StatusBadRequest)
		return
	}
	
	var user models.User
	res := repository.GetById(&user, id)
	if res.Error != nil {
		logger.Logger.Printf("Could not find user with id %d", id)
		http.Error(w, "Could not find user with id", http.StatusNotFound)
		return
	}

	user.Name = body.Name
	user.Surname = body.Surname
	user.Age = body.Age
	user.Gender = body.Gender
	user.Nationality = body.Nationality

	save := repository.SaveInDb(&user)
	if save.Error != nil {
		logger.Logger.Printf("Could not update user with id %d", id)
		http.Error(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	logger.Logger.Println("User updated successfully!")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
