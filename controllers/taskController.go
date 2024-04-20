package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"taskManagement/utils"

	"github.com/gorilla/mux"
)

type Task struct {
	TaskID   		int64   	`json:"taskId"`
	Description 	string 		`json:"description"`
	Priority 		string 		`json:"priority"`
	UserID			int64		`json:"userId"`
}

type Response struct {
	Status 		string 			`json:"status"`
	Result 		interface{} 	`json:"result"`
}

var db *sql.DB

func init() {
	db = utils.ReturnDB()
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	resultRows, err := db.Query("SELECT * FROM `taskManagement`.`task`")
	if err != nil {
		fmt.Println("Error Fetching tasks from DB", err)
		return 
	}
	
	var tasks []Task
	for resultRows.Next() {
		var task Task
		if err := resultRows.Scan(&task.TaskID, & task.Description, &task.Priority, &task.UserID); err != nil {
			fmt.Println("Error Scanning Values", err)
			return 
		}

		tasks = append(tasks, task)
	}

	bytes, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Error marshaling data", err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error Reading Request body", err)
		return 
	}
	if err  := json.Unmarshal(bytes, &task); err != nil {
		fmt.Println("Error Unmarshaling data", err)
		return 
	}
	
	_= db.QueryRow("INSERT INTO `taskManagement`.`task` (`id`, `description`, `priority`, `user_id`) VALUES (?, ?, ?, ?)",task.TaskID, task.Description, task.Priority, task.UserID).Scan(&task.TaskID, &task.Description, &task.Priority, &task.UserID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetTaskById(w http.ResponseWriter, r * http.Request){
	params := mux.Vars(r)
	taskId, err := strconv.ParseInt(params["taskId"], 0, 0)
	if err != nil {
		fmt.Println("Error Parsing taskId", err)
		return 
	}

	var task Task
	_ = db.QueryRow("SELECT * FROM `taskManagement`.`task` WHERE (id = ?)", taskId).Scan(&task.TaskID, &task.Description, &task.Priority, &task.UserID)

	if task.TaskID == 0 {
		response := Response{Status: "Fail", Result: "No such task Exists"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func UpdateTaskById(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	taskId, _ := strconv.ParseInt(params["taskId"], 0, 0)
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error Reading request body", err)
		return
	}

	var task Task
	unmarshalErr := json.Unmarshal(bytes, &task)
	if unmarshalErr != nil {
		fmt.Println("Error unmarshaling data")
	}

	row := db.QueryRow("SELECT `id` FROM `taskManagement`.`task` WHERE (id = ?)", taskId)
	err = row.Scan(&task.TaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			response := Response{Status: "Fail", Result: "No such task Exists"}
        	w.Header().Set("Content-Type", "application/json")
        	w.WriteHeader(http.StatusBadRequest)
        	json.NewEncoder(w).Encode(response)
        	return
    	}
    	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    	fmt.Println("Error querying database:", err)
    	return
	}

	stmt, err := db.Prepare("UPDATE `taskManagement`.`task` SET `description` = ?, `priority` = ?, `user_id` = ? WHERE `id` = ?")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Invalid SQL syntax", err)
		return 
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Description, task.Priority, task.UserID, taskId)
	if err != nil {
		fmt.Println("Error Updating Data in DB", err)
		return
	}
	
	_ = db.QueryRow("SELECT * FROM `taskManagement`.`task` WHERE (id = ?)", taskId).Scan(&task.TaskID, &task.Description, &task.Priority, &task.UserID)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func DeleteTaskById(w http.ResponseWriter, r * http.Request) {
	params := mux.Vars(r)
	taskId, _ := strconv.ParseInt(params["taskId"], 0, 0)

	var task Task
	_ = db.QueryRow("SELECT `id` FROM `taskManagement`.`task` WHERE (id = ?)", taskId).Scan(&task.TaskID)
	
	if task.TaskID == 0 {
		response := Response{Status: "Fail", Result: "No such task Exists"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return 
	}

	_ = db.QueryRow("DELETE FROM `taskManagement`.`task` WHERE (id = ?)", taskId)
	w.WriteHeader(http.StatusNoContent)
}

