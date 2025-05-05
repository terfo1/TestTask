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
	mux.Post("/user", handler.CreateUser)
	mux.Get("/user", handler.GetUsers)
	mux.Put("/user", handler.UpdateUser)
	mux.Delete("/user", handler.DeleteUser)
	return mux
}
