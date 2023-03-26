package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func main() {
	// Load .env file and later use variables and credentials
	if err := godotenv.Load("development.env"); err != nil {
		log.Fatal("Error while loading .env file!")
	}

	// mysql database config
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id INT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		done BOOLEAN NOT NULL DEFAULT FALSE
	)`); err != nil {
		log.Fatal(err)
	}

	// pingErr := db.Ping()
	// if pingErr != nil {
	// 	log.Fatal(pingErr)
	// }
	// fmt.Println("Connected!")

	// Register the HTTP handlers
	http.HandleFunc("/api/tasks", getAllTasksHandler(db))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getAllTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := getAllTasks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJsonResponse(w, tasks)
	}
}

func getAllTasks(db *sql.DB) (tasks []Task, err error) {
	rows, err := db.Query("SELECT id, title, description, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks = make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func writeJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
