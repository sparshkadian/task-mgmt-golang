package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"taskManagement/models"
)


func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error Reading request body", err)
		return 
	}
	if err = json.Unmarshal(bytes, &user); err != nil {
		fmt.Println("Error unmarshaling Data", err)
		return 
	}
	newUser := user.PerformSignup(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}