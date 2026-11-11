package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

// authMiddleware verifies the JWT before allowing access.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1. Read the Authorization header.
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		// 2. Extract the JWT.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// 3. Read the signing secret.
		secret := os.Getenv("JWT_SECRET")

		if secret == "" {
			http.Error(w, "Authentication unavailable", http.StatusInternalServerError)
			return
		}

		// 4. Parse and verify the JWT.
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 5. Extract the user ID from the token.
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, err := claims.GetExpirationTime()
		if err != nil || userID == nil {
			http.Error(w, "Invalid token expiration", http.StatusUnauthorized)
			return
		}

		idValue, ok := claims["id"].(float64)
		if !ok || idValue < 1 || idValue != float64(int64(idValue)) {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		// 6. Store the user ID in the request context.
		ctx := context.WithValue(
			r.Context(),
			userIDKey,
			int64(idValue),
		)

		// 7. Continue to the protected handler.
		next(w, r.WithContext(ctx))
	}

}

func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the authenticated user ID
		userID, ok := r.Context().Value(userIDKey).(int64)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Get the user's current role from MySQL
		var role string

		err := db.QueryRow(
			"SELECT role FROM users WHERE id = ?",
			userID,
		).Scan(&role)

		// 3. User no longer exists
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 4. Handle database errors
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		// 5. Check admin permission
		if role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// 6. Allow the admin request
		next(w, r)
	}
}
