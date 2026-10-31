package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	Priority  int    `json:"priority"` // جديد: 1 أحمر، 2 برتقالي، 3 أخضر
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "todos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		text       TEXT NOT NULL,
		completed  INTEGER NOT NULL DEFAULT 0,
		priority   INTEGER NOT NULL DEFAULT 3,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`) // جديد: عمود priority
	if err != nil {
		log.Fatal(err)
	}

	// جديد: لقاعدة بيانات todos.db القديمة التي لا تحتوي العمود.
	// إن كان العمود موجودًا فسيظهر خطأ "duplicate column" ونتجاهله.
	_, err = db.Exec("ALTER TABLE tasks ADD COLUMN priority INTEGER NOT NULL DEFAULT 3")
	if err != nil && !strings.Contains(err.Error(), "duplicate column") {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/tasks", listTasks)
	mux.HandleFunc("POST /api/tasks", createTask)
	mux.HandleFunc("PATCH /api/tasks/{id}", updateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", deleteTask)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, text, completed, priority FROM tasks ORDER BY id") // جديد
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		var done int
		if err := rows.Scan(&t.ID, &t.Text, &done, &t.Priority); err != nil { // جديد
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		t.Completed = done == 1
		tasks = append(tasks, t)
	}
	writeJSON(w, http.StatusOK, tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	text := strings.TrimSpace(body.Text)
	if text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO tasks (text) VALUES (?)", text)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	writeJSON(w, http.StatusCreated, Task{ID: id, Text: text, Priority: 3}) // جديد
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// جديد: المؤشر (*) يسمح بمعرفة هل أرسل المتصفح الحقل أم لا
	var body struct {
		Completed *bool `json:"completed"`
		Priority  *int  `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if body.Completed == nil && body.Priority == nil {
		http.Error(w, "nothing to update", http.StatusBadRequest)
		return
	}
	if body.Priority != nil && (*body.Priority < 1 || *body.Priority > 3) {
		http.Error(w, "priority must be 1, 2 or 3", http.StatusBadRequest)
		return
	}

	var found int64
	if body.Completed != nil {
		done := 0
		if *body.Completed {
			done = 1
		}
		result, err := db.Exec("UPDATE tasks SET completed = ? WHERE id = ?", done, id)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		n, _ := result.RowsAffected()
		found += n
	}
	if body.Priority != nil {
		result, err := db.Exec("UPDATE tasks SET priority = ? WHERE id = ?", *body.Priority, id)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		n, _ := result.RowsAffected()
		found += n
	}
	if found == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if _, err := db.Exec("DELETE FROM tasks WHERE id = ?", id); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
