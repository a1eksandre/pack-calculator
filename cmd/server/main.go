package main

import (
	"log"
	"net/http"
	"os"

	"github.com/a1eksandre/pack-calculator/internal/api"
)

func main() {
	defaultPackSizes := []int{250, 500, 1000, 2000, 5000}
	server := api.NewServer(defaultPackSizes)

	mux := http.NewServeMux()

	// Serve UI
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	// Serve API (no StripPrefix here)
	mux.Handle("/api/", server.Routes())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
