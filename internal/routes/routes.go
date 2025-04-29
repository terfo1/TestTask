package routes

import (
	"TestTask/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", httpSwagger.WrapHandler))
	mux.HandleFunc("/createuser", handler.CreateUser)
	mux.HandleFunc("/user", handler.GetUsers)
	mux.HandleFunc("/updateuser", handler.UpdateUser)
	mux.HandleFunc("/deleteuser", handler.DeleteUser)
	return mux
}
