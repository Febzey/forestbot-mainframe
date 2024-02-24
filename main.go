package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

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
	fmt.Println(ascii_art_start_message)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// Initialize logger
	logger := &logger.Logger{
		Logger: log.New(os.Stdout, "", 0),
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		logger.Error(err.Error())
		log.Fatal("Failed to connect to the database")
	}
	defer func() {
		// Defer the closing of the database connection
		if err := db.CloseDb(); err != nil {
			logger.Error(fmt.Sprintln("Error closing database connection:", err))
		} else {
			logger.Info("Database connection closed")
		}
	}()

	logger.Success("Connected to the database")

	// Create a new router
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	// Create a controller
	controller := routes.NewController(db, logger)

	// Load and handle routes
	routes.LoadAndHandleRoutes(r, controller)

	// Set up CORS middleware
	r.Use(mux.CORSMethodMiddleware(r))

	// Configure server
	server := ServerConfig(r)

	// Start the server in a goroutine
	go func() {
		// Log server start information
		logger.Info(fmt.Sprintf("Server is starting on port: %s", os.Getenv("SERVER_PORT")))
		if err := server.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintln("Server error:", err))
		}
	}()

	// Wait for Ctrl+C signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shutdown the server gracefully
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error(fmt.Sprintln("Error shutting down the server:", err))
	} else {
		logger.Info("Server shut down gracefully")
	}

	// Log server shutdown
	logger.Info("Server has stopped.")
}
