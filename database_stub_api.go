package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nikvkov/database-stub-api/handlers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.IndexHandler).Methods(http.MethodGet)
	router.HandleFunc("/supplementary-data", handlers.LoggingMiddleware(handlers.SupplementaryDataHandler)).Methods(http.MethodPost)
	router.HandleFunc("/disclosure-identification", handlers.LoggingMiddleware(handlers.DisclosureIdentificationHandler)).Methods(http.MethodPost)
	//router.Use(handlers.LoggingMiddleware)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8899"
		log.Printf("Defaulting to port %s", port)
	}

	http.Handle("/", router)
	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
