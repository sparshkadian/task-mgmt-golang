package routes

import (
	"taskManagement/controllers"

	"github.com/gorilla/mux"
)

func RegisterTaskRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", controllers.GetAllTasks).Methods("GET")
	r.HandleFunc("/tasks", controllers.CreateNewTask).Methods("POST")
	r.HandleFunc("/tasks/{taskId}", controllers.GetTaskById).Methods("GET")
	r.HandleFunc("/tasks/{taskId}", controllers.DeleteTaskById).Methods("DELETE")
	r.HandleFunc("/tasks/{taskId}", controllers.UpdateTaskById).Methods("PUT")
}