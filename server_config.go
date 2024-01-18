package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ConnUrlBuilder(n string) (string, error) {
	var url string

	switch n {
	case "server":
		url = fmt.Sprintf(
			"%s:%s",
			os.Getenv("SERVER_HOST"),
			os.Getenv("SERVER_PORT"),
		)
	default:
		return "", fmt.Errorf("invalid connection type")
	}

	return url, nil
}

func ServerConfig(router *mux.Router) *http.Server {
	serverConnUrl, _ := ConnUrlBuilder("server")
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT_SECONDS_COUNT"))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("*")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	return &http.Server{
		Handler:     handlers.CORS(headersOk, originsOk, methodsOk)(router),
		Addr:        serverConnUrl,
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}

}

func StartServer(server *http.Server) {
	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	log.Println("Stopping server...")

}
