package models

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"taskManagement/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   uint8  `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	UserID   uint8  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

var db *sql.DB

func init(){
	db  = utils.ReturnDB()
}

func (user *User) PerformSignup(w http.ResponseWriter) *NewUser{
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		fmt.Println("All fields are required")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 1)
	if err != nil {
		fmt.Println("Error hashing Password", err)
		return nil 
	}

	user.Password = string(hashedPassword)
	rand.NewSource(time.Now().UnixNano())
	user.UserID = uint8(rand.Intn(100))

	resultRow, err := db.Query("INSERT INTO `taskManagement`.`users`(`id`, `name`, `email`, `password`) VALUES (?, ?, ? ,?)", user.UserID, user.Name, user.Email, user.Password)
	if err != nil {
		fmt.Println("Error Creating new User", err)
		return nil
	}
	for resultRow.Next() {
		if err := resultRow.Scan(&user.UserID, &user.Name, &user.Email); err != nil {
			fmt.Println("Error scanning values", err)
			return nil
		}
	}

	newUser := NewUser{
		UserID: user.UserID,
		Name: user.Name,
		Email: user.Email,
	}

	return &newUser
}