package routes

import (
	"TestTask/internal/handler"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func SetupRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	mux.HandleFunc("/createuser", handler.CreateUser)
	mux.HandleFunc("/user", handler.GetUsers)
	mux.HandleFunc("/updateuser", handler.UpdateUser)
	mux.HandleFunc("/deleteuser", handler.DeleteUser)
	return mux
}
