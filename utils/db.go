package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func DBConnect() {
	var err error
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	db ,err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error verifying conn arguments", err)
		return 
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error Pinging DB", err)
		return 
	}else {
		fmt.Println("DB Connection Successfull!")
	}
}

func ReturnDB() *sql.DB {
	DBConnect()
	return db
}