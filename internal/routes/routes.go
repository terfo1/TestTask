package routes

import (
	"TestTask/internal/handler"
	"net/http"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/createuser", handler.CreateUser)
	mux.HandleFunc("/user", handler.GetUsers)
	// Not implemented yet
	// mux.HandleFunc("/user", handler.UpdateUsers)
	// mux.HandleFunc("/user", handler.DeleteUsers)
	return mux
}
