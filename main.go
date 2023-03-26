package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

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
	http.HandleFunc("/api/tasks/", taskHandler(db))
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

func taskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := getTaskIdFromUrl(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch r.Method {
		// Handle GET method
		case http.MethodGet:
			task, err := getTaskById(db, taskId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if task == nil {
				http.NotFound(w, r)
				return
			}
			writeJsonResponse(w, task)

		// Handle PUT method
		case http.MethodPut:
			task, err := decodeTaskFromBody(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Id = taskId
			err = updateTask(db, task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJsonResponse(w, task)

		// Handle POST method
		case http.MethodPost:
			task, err := decodeTaskFromBody(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			addTask(db, task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJsonResponse(w, task)

		// Handke DELETE method
		case http.MethodDelete:
			err = deleteTask(db, taskId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Task with ID %v deleted successfully\n", taskId)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
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
	return
}

func writeJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addTask(db *sql.DB, task *Task) error {
	_, err := db.Exec("INSERT INTO tasks (title, description, done) VALUES ('?', '?', ?);", task.Title, task.Description, task.Done)
	return err
}

func updateTask(db *sql.DB, task *Task) error {
	_, err := db.Exec("UPDATE tasks SET title=?, description=?, done=? WHERE id=?", task.Title, task.Description, task.Done, task.Id)
	return err
}

func deleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id=?", id)
	return err
}

func getTaskById(db *sql.DB, id int) (*Task, error) {
	row := db.QueryRow("SELECT id, title, description, done FROM tasks WHERE id=?", id)
	var task Task
	if err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Done); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func getTaskIdFromUrl(urlPath string) (int, error) {
	re := regexp.MustCompile(`^/api/tasks/(\d+)$`)
	matches := re.FindStringSubmatch(urlPath)
	if matches == nil {
		return 0, fmt.Errorf("%s", "Invalid task ID")
	}
	taskId, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("%s", "Invalid task ID")
	}
	return taskId, nil
}

func decodeTaskFromBody(body io.ReadCloser) (*Task, error) {
	defer body.Close()
	var task Task
	if err := json.NewDecoder(body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}
