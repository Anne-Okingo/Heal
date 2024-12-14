package main

import (
	"fmt"
	"net/http"

	"Heal/internals/routes"
	"Heal/utils"
)

func main() {
	// Initialize the database connection
	utils.Getdb("./Heal.db")

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	wrappedMux := routes.RouteChecker(mux)

	// Create an HTTP server and handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	fmt.Println("server running @http://localhost:8080\nTo stop server type 'exit'\n=====================================")

	// Start the HTTP server and listen for incoming requests
	err := server.ListenAndServe()
	if err != nil {
		utils.ErrorHandler("web")
	}
}
