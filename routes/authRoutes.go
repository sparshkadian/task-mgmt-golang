package routes

import (
	"taskManagement/controllers"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r * mux.Router){
	r.HandleFunc("/api/signup", controllers.Signup).Methods("POST")
}