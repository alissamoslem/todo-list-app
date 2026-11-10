package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os" // تغيير MySQL: لقراءة اسم المستخدم وكلمة المرور من متغيرات البيئة
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql" // تغيير MySQL: مشغّل MySQL بدل مشغّل SQLite
)

// جديد: يقرأ متغير بيئة، ويستعمل القيمة الافتراضية إن كان فارغًا
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type Task struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	Priority  int    `json:"priority"` // جديد: 1 أحمر، 2 برتقالي، 3 أخضر
}

var db *sql.DB

func main() {
	// جديد: تحميل ملف .env إن وُجد. المتغيرات الموجودة أصلًا في الطرفية لها الأولوية.
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}
	// تغيير MySQL: نص الاتصال بالشكل user:password@tcp(host:port)/database
	// parseTime=true: يحوّل أعمدة التاريخ إلى time.Time
	// clientFoundRows=true: يجعل RowsAffected يعدّ الصفوف المطابقة وليس المتغيّرة فقط،
	// وإلا يعيد MySQL الرقم 0 عند تحديث قيمة إلى نفسها فيظهر "task not found" خطأً
	// جديد: كل قيمة تأتي من .env مع قيمة افتراضية مناسبة. كلمة المرور بلا قيمة افتراضية
	dsn := getenv("DB_USER", "root") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + getenv("DB_HOST", "127.0.0.1") + ":" + getenv("DB_PORT", "3306") + ")/" +
		getenv("DB_NAME", "todo_app") + "?parseTime=true&clientFoundRows=true"

	var err error
	db, err = sql.Open("mysql", dsn) // تغيير MySQL: "mysql" بدل "sqlite"، واسم القاعدة todo_app داخل نص الاتصال
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// تغيير MySQL: Ping يتصل فعليًا بالخادم، فيظهر خطأ كلمة المرور أو توقف MySQL فورًا
	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to MySQL: ", err)
	}

	// تغيير MySQL: حُذف db.SetMaxOpenConns(1) لأنه خاص بـ SQLite
	// تغيير MySQL: حُذف CREATE TABLE و ALTER TABLE لأن الجدول يُنشأ من ملف schema.sql

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/tasks", listTasks)
	mux.HandleFunc("POST /api/tasks", createTask)
	mux.HandleFunc("PATCH /api/tasks/{id}", updateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", deleteTask)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	addr := ":" + getenv("PORT", "8080") // جديد
	log.Println("Running at http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Try make it using ORM
func listTasks(w http.ResponseWriter, r *http.Request) {
	// تغيير MySQL: text → title و completed → is_completed (priority كما هو)
	rows, err := db.Query("SELECT id, title, is_completed, priority FROM tasks ORDER BY id")
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		var done int
		if err := rows.Scan(&t.ID, &t.Text, &done, &t.Priority); err != nil {
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

	// تغيير MySQL: اسم العمود title بدل text. العمود priority يأخذ DEFAULT 3 من قاعدة البيانات
	result, err := db.Exec("INSERT INTO tasks (title) VALUES (?)", text)
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
		// تغيير MySQL: completed → is_completed. و MySQL يحدّث updated_at تلقائيًا
		result, err := db.Exec("UPDATE tasks SET is_completed = ? WHERE id = ?", done, id)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		n, _ := result.RowsAffected()
		found += n
	}
	if body.Priority != nil {
		// بدون تغيير: العمود priority له الاسم نفسه في MySQL
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
	// بدون تغيير: جملة DELETE نفسها تعمل في MySQL
	if _, err := db.Exec("DELETE FROM tasks WHERE id = ?", id); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
