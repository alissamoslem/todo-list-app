package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func listUsers(w http.ResponseWriter, r *http.Request) {
	// Query only the fields we want to expose
	rows, err := db.Query(
		"SELECT id, name, email, role FROM users ORDER BY id",
	)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Define the response structure
	type AdminUser struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	users := []AdminUser{}

	for rows.Next() {
		var user AdminUser

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
		)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func listAllTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT
            t.id,
            t.title,
            t.is_completed,
            t.priority,
            t.user_id,
            COALESCE(u.name, '')
        FROM tasks t
        LEFT JOIN users u ON t.user_id = u.id
        ORDER BY t.id
    `)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type AdminTask struct {
		ID        int64  `json:"id"`
		Text      string `json:"text"`
		Completed bool   `json:"completed"`
		Priority  int    `json:"priority"`
		UserID    *int64 `json:"user_id"`
		OwnerName string `json:"owner_name"`
	}

	tasks := []AdminTask{}

	for rows.Next() {
		var task AdminTask
		var completed int

		err := rows.Scan(
			&task.ID,
			&task.Text,
			&completed,
			&task.Priority,
			&task.UserID,
			&task.OwnerName,
		)

		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		task.Completed = completed == 1
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// 1. Get authenticated admin ID
	adminID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Get the user ID to delete
	targetID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || targetID <= 0 {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	// 3. Prevent admin from deleting their own account
	if adminID == targetID {
		http.Error(w, "cannot delete your own account", http.StatusForbidden)
		return
	}

	// 4. Require explicit confirmation
	var body struct {
		Confirmation string `json:"confirmation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "confirmation required", http.StatusBadRequest)
		return
	}

	if body.Confirmation != "DELETE" {
		http.Error(w, "type DELETE to confirm deletion", http.StatusBadRequest)
		return
	}

	// 5. Delete user from database
	result, err := db.Exec(
		"DELETE FROM users WHERE id = ?",
		targetID,
	)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	// 6. Check if user existed
	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// 7. Return success
	w.WriteHeader(http.StatusNoContent)
}
