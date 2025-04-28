package main

import (
	"TestTask/internal/config"
	"TestTask/internal/database"
	"TestTask/internal/routes"
	"TestTask/pkg/logger"
	"flag"
	"net/http"
)

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
