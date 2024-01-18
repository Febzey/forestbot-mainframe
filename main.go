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

// players := []types.PlayerList{
// 	{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
// }

// dc := utils.RenderTab(players)

// // Determine the directory of the executabl
// if err := dc.SavePNG("tablist.png"); err != nil {
// 	fmt.Println("Error saving PNG:", err)
// 	return
// }

// fmt.Printf("Tablist image saved successfully")
