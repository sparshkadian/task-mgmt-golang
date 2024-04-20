package main

import (
	"fmt"
	"log"
	"net/http"
	"taskManagement/routes"

	"github.com/gorilla/mux"
)


func main() {
	
	r := mux.NewRouter()

	routes.RegisterTaskRoutes(r)
	routes.RegisterAuthRoutes(r)

	fmt.Println("Server started at Port:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Error starting server", err)
	}
	
}