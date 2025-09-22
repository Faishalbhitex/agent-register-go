// main.go
package main

import (
	"agent-register-go/database"
	"agent-register-go/routers"
	"log"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Setup router and start server
	r := routers.SetupRouter()
	log.Println("Agent Register server starting on :8080")
	r.Run(":8080")
}
