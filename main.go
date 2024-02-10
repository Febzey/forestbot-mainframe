package main

import (
	"fmt"
	"log"
	"os"

	"github.com/febzey/ForestBot-Mainframe/database"
	"github.com/febzey/ForestBot-Mainframe/logger"
	"github.com/febzey/ForestBot-Mainframe/middleware"
	"github.com/febzey/ForestBot-Mainframe/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	logger := &logger.Logger{
		Logger: log.New(os.Stdout, "", 0),
	}

	db, err := database.Connect()

	defer func() {
		err := db.CloseDb()
		if err != nil {
			logger.Error(err.Error())
		}
		logger.Success("Database connection closed")
	}()

	logger.Success("Connected to the database")

	if err != nil {
		logger.Error(err.Error())
	}

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	routes.LoadAndHandleRoutes(r, db, logger)

	r.Use(mux.CORSMethodMiddleware(r))

	server := ServerConfig(r)
	StartServer(server)

}
