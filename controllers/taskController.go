package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"taskManagement/models"
)


func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	tasks := task.GetAllTasksDB()
	bytes, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Error marshaling data", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error Reading Request body", err)
		return 
	}
	if err  := json.Unmarshal(bytes, &task); err != nil {
		fmt.Println("Error Unmarshaling data", err)
		return 
	}
	
	returnedTask := task.CreateNewTaskDB()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnedTask)
}

func GetTaskById(w http.ResponseWriter, r * http.Request){
	var task models.Task
	returnedTask := task.GetTaskByIdDB(w, r)
	if returnedTask == nil {
		return 
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(returnedTask)
	}
}

func UpdateTaskById(w http.ResponseWriter, r *http.Request){
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error Reading request body", err)
		return
	}

	var task models.Task
	unmarshalErr := json.Unmarshal(bytes, &task)
	if unmarshalErr != nil {
		fmt.Println("Error unmarshaling data")
	}

	returnedTask := task.UpdateTaskByIdDB(w, r)
	if returnedTask == nil {
		return 
	} else {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}
	
}

func DeleteTaskById(w http.ResponseWriter, r * http.Request) {
	var task models.Task
	task.DeleteTaskByIdDB(w, r)
}

