package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	TaskID      int64  `json:"taskId"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	UserID      int64  `json:"userId"`
}

type Response struct {
	Status 		string 			`json:"status"`
	Result 		interface{} 	`json:"result"`
}

func (task *Task) GetAllTasksDB() *Response {
	resultRows, err := db.Query("SELECT * FROM `taskManagement`.`task`")
	if err != nil {
		fmt.Println("Error Fetching tasks from DB", err)
		return nil 
	}

	var tasks []Task
	for resultRows.Next() {
		var task Task
		if err := resultRows.Scan(&task.TaskID, & task.Description, &task.Priority, &task.UserID); err != nil {
			fmt.Println("Error Scanning Values", err)
			return nil  
		}

		tasks = append(tasks, task)
	}
	
	response := Response{Status: "Success", Result: tasks}
	return &response
}

func (task *Task) CreateNewTaskDB() *Task{
	_= db.QueryRow("INSERT INTO `taskManagement`.`task` (`id`, `description`, `priority`, `user_id`) VALUES (?, ?, ?, ?)",task.TaskID, task.Description, task.Priority, task.UserID).Scan(&task.TaskID, &task.Description, &task.Priority, &task.UserID)
	return task
}

func (task *Task) GetTaskByIdDB(w http.ResponseWriter, r *http.Request) *Response{
	params := mux.Vars(r)
	taskId, err := strconv.ParseInt(params["taskId"], 0, 0)
	if err != nil {
		fmt.Println("Error Parsing taskId", err)
		return nil
	}
	_ = db.QueryRow("SELECT * FROM `taskManagement`.`task` WHERE (id = ?)", taskId).Scan(&task.TaskID, &task.Description, &task.Priority, &task.UserID)
	if task.TaskID == 0 {
		response := Response{Status: "Fail", Result: "No such task Exists"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return nil
	}

	response := Response{Status: "Success", Result: task}
	return &response
}

func (task *Task) UpdateTaskByIdDB(w http.ResponseWriter, r *http.Request) *Task{
	params := mux.Vars(r)
	taskId, _ := strconv.ParseInt(params["taskId"], 0, 0)

	row := db.QueryRow("SELECT `id` FROM `taskManagement`.`task` WHERE (id = ?)", taskId)
	err := row.Scan(&task.TaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			response := Response{Status: "Fail", Result: "No such task Exists"}
        	w.Header().Set("Content-Type", "application/json")
        	w.WriteHeader(http.StatusBadRequest)
        	json.NewEncoder(w).Encode(response)
        	return nil
    	}
    	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    	fmt.Println("Error querying database:", err)
    	return nil 
	}

	stmt, err := db.Prepare("UPDATE `taskManagement`.`task` SET `description` = ?, `priority` = ?, `user_id` = ? WHERE `id` = ?")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Invalid SQL syntax", err)
		return nil
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Description, task.Priority, task.UserID, taskId)
	if err != nil {
		fmt.Println("Error Updating Data in DB", err)
		return nil
	}

	_ = db.QueryRow("SELECT * FROM `taskManagement`.`task` WHERE (id = ?)", taskId).Scan(&task.TaskID, &task.Description, &task.Priority, &task.UserID)
	return task
}


func (task *Task) DeleteTaskByIdDB(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskId, _ := strconv.ParseInt(params["taskId"], 0, 0)

	_ = db.QueryRow("SELECT `id` FROM `taskManagement`.`task` WHERE (id = ?)", taskId).Scan(&task.TaskID)

	if task.TaskID == 0 {
		response := Response{Status: "Fail", Result: "No such task Exists"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return 
	}

	_ = db.QueryRow("DELETE FROM `taskManagement`.`task` WHERE (id = ?)", task.TaskID)
		w.WriteHeader(http.StatusNoContent)
}