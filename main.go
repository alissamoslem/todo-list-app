package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// getenv reads an environment variable and uses a fallback if empty.
func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// requireEnv ensures required environment variables are configured.
func requireEnv(keys ...string) {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			log.Fatalf("%s environment variable is required", key)
		}
	}
}

type Task struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	Priority  int    `json:"priority"` // 1 red, 2 orange, 3 green
}

var db *sql.DB

func main() {
	// 1. Load .env for local development.
	// Production hosts should provide environment variables directly.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	// 2. Validate required configuration.
	requireEnv(
		"JWT_SECRET",
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	)

	// 3. Build MySQL configuration safely.
	config := mysql.NewConfig()

	config.User = os.Getenv("DB_USER")
	config.Passwd = os.Getenv("DB_PASSWORD")
	config.Net = "tcp"
	config.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	config.DBName = os.Getenv("DB_NAME")

	config.ParseTime = true
	config.ClientFoundRows = true

	dsn := config.FormatDSN()

	// 4. Open the database connection pool.
	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Cannot initialize MySQL: ", err)
	}
	defer db.Close()

	// Configure connection pooling.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Minute)

	// Verify that MySQL is reachable.
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to MySQL: ", err)
	}

	log.Println("Connected to MySQL successfully")

	// 5. Configure HTTP routes.
	mux := http.NewServeMux()

	// Authentication routes.
	mux.HandleFunc("POST /api/auth/register", registerHandler)
	mux.HandleFunc("POST /api/auth/login", loginHandler)
	mux.HandleFunc("GET /api/auth/me", authMiddleware(meHandler))

	// User task routes.
	mux.HandleFunc("GET /api/tasks", authMiddleware(listTasks))
	mux.HandleFunc("POST /api/tasks", authMiddleware(createTask))
	mux.HandleFunc("PATCH /api/tasks/{id}", authMiddleware(updateTask))
	mux.HandleFunc("DELETE /api/tasks/{id}", authMiddleware(deleteTask))

	// Admin routes.
	mux.HandleFunc(
		"GET /api/admin/users",
		authMiddleware(adminMiddleware(listUsers)),
	)
	mux.HandleFunc(
		"GET /api/admin/tasks",
		authMiddleware(adminMiddleware(listAllTasks)),
	)
	mux.HandleFunc(
		"DELETE /api/admin/users/{id}",
		authMiddleware(adminMiddleware(deleteUser)),
	)

	// Serve HTML, CSS, and JavaScript.
	mux.Handle("/", http.FileServer(http.Dir("static")))

	// 6. Start the HTTP server.
	port := getenv("PORT", "8080")
	addr := ":" + port

	server := &http.Server{
		Addr:              addr,
		Handler:           corsMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Server listening on %s", addr)

	log.Fatal(server.ListenAndServe())
}

// writeJSON sends a JSON response with the requested HTTP status.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}

// listTasks returns tasks belonging to the authenticated user.
func listTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.QueryContext(
		r.Context(),
		`SELECT id, title, is_completed, priority
		 FROM tasks
		 WHERE user_id = ?
		 ORDER BY id`,
		userID,
	)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var task Task
		var done int

		if err := rows.Scan(
			&task.ID,
			&task.Text,
			&done,
			&task.Priority,
		); err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		task.Completed = done == 1
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

// createTask creates a task owned by the authenticated user.
func createTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	result, err := db.ExecContext(
		r.Context(),
		"INSERT INTO tasks (title, user_id) VALUES (?, ?)",
		text,
		userID,
	)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, Task{
		ID:       id,
		Text:     text,
		Priority: 3,
	})
}

// updateTask updates completion and/or priority of an owned task.
func updateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

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

	if body.Priority != nil &&
		(*body.Priority < 1 || *body.Priority > 3) {
		http.Error(
			w,
			"priority must be 1, 2 or 3",
			http.StatusBadRequest,
		)
		return
	}

	// Verify task ownership.
	var exists int

	err = db.QueryRowContext(
		r.Context(),
		"SELECT 1 FROM tasks WHERE id = ? AND user_id = ?",
		id,
		userID,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	// Update completion when provided.
	if body.Completed != nil {
		done := 0
		if *body.Completed {
			done = 1
		}

		_, err := db.ExecContext(
			r.Context(),
			`UPDATE tasks
			 SET is_completed = ?
			 WHERE id = ? AND user_id = ?`,
			done,
			id,
			userID,
		)

		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
	}

	// Update priority when provided.
	if body.Priority != nil {
		_, err := db.ExecContext(
			r.Context(),
			`UPDATE tasks
			 SET priority = ?
			 WHERE id = ? AND user_id = ?`,
			*body.Priority,
			id,
			userID,
		)

		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// deleteTask deletes a task only if it belongs to the user.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := db.ExecContext(
		r.Context(),
		"DELETE FROM tasks WHERE id = ? AND user_id = ?",
		id,
		userID,
	)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
