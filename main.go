package main

import (
	"log"
	"net/http"

	"github.com/BerylCAtieno/string-analyzer/database"
	"github.com/BerylCAtieno/string-analyzer/routes"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Register API endpoints
	routes.RegisterRoutes()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
