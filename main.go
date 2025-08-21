package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Dsouza10082/documentOCRProcessor/handler"
	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	server := &http.Server{
		Addr:         ":8082",
		Handler:      r,
		ReadTimeout:  240 * time.Second, // 4 minutes testing purposes
		WriteTimeout: 240 * time.Second, // 4 minutes testing purposes
		IdleTimeout:  240 * time.Second, // 4 minutes testing purposes
	}

	// Test route
	handler.Routes(r)
 
	port := os.Getenv("PORT")

	if port == "" {
		port = "8082"
	}

	log.Printf("Starting server on http://localhost:%s", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

