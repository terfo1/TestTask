package main

import (
	"TestTask/internal/config"
	"TestTask/internal/database"
	"TestTask/internal/routes"
	"TestTask/pkg/logger"
	"flag"
	"net/http"
)

// @title 			Test Task
// @version         1.0
// @description     This is an API for enriching user info with age, gender, and nationality.
// @host            localhost:8080
// @BasePath        /

func main() {
	logger.InitLog()
	config.LoadEnv()
	database.ConnectToDB()
	database.SyncDB()
	mux := routes.SetupRoutes()
	addr := flag.String("addr", ":8080", "http network addr")
	logger.Logger.Println("Server starting at", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		logger.Logger.Fatal("Could not start a server!", err)
	}
}
