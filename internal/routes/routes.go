package routes

import (
	"TestTask/internal/handler"
	"net/http"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/createuser", handler.CreateUser)
	mux.HandleFunc("/user", handler.GetUsers)
	return mux
}
