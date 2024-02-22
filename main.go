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

var (
	ascii_art_start_message = `
	
███████╗ ██████╗ ██████╗ ███████╗███████╗████████╗██████╗  ██████╗ ████████╗     ██████╗██████╗     ███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗     
██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔════╝╚══██╔══╝██╔══██╗██╔═══██╗╚══██╔══╝    ██╔════╝╚════██╗    ██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗    
█████╗  ██║   ██║██████╔╝█████╗  ███████╗   ██║   ██████╔╝██║   ██║   ██║       ██║      █████╔╝    ███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝    
██╔══╝  ██║   ██║██╔══██╗██╔══╝  ╚════██║   ██║   ██╔══██╗██║   ██║   ██║       ██║     ██╔═══╝     ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗    
██║     ╚██████╔╝██║  ██║███████╗███████║   ██║   ██████╔╝╚██████╔╝   ██║       ╚██████╗███████╗    ███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║    
╚═╝      ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝   ╚═╝   ╚═════╝  ╚═════╝    ╚═╝        ╚═════╝╚══════╝    ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝                                                                                                                                                	




`
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

	logger.Success("Connected to the database")

	if err != nil {
		logger.Error(err.Error())
	}

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	routes.LoadAndHandleRoutes(r, db, logger)

	r.Use(mux.CORSMethodMiddleware(r))

	logger.Info(fmt.Sprintf("Server is Starting on port: %s", os.Getenv("SERVER_PORT")))

	fmt.Println(ascii_art_start_message)

	server := ServerConfig(r)
	StartServer(server)

	err = db.CloseDb()
	if err != nil {
		logger.Error(fmt.Sprint("Error closing database connection", err.Error()))
	}
	logger.Info("Database connection closed")
	logger.Info("Server has stopped.")
}
