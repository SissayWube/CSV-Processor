// Package main serves as the entry point for the backend CSV processing application.
package main

import "csv_processor/handlers"

func main() {

	// Initialize the Gin router with all defined API routes and middleware.
	router := handlers.SetupRouter()

	// Start the HTTP server.
	router.Run()

}
