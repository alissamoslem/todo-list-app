package main

import (
	"net/http"
	"os"
	"strings"
)

// corsMiddleware controls which browser origins
// are allowed to access our API.
func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigin := os.Getenv("FRONTEND_ORIGIN")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS applies to API requests.
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")

		// Tell caches that responses can vary by Origin.
		w.Header().Add("Vary", "Origin")

		// Handle CORS preflight requests.
		if r.Method == http.MethodOptions {
			w.Header().Add(
				"Vary",
				"Access-Control-Request-Method",
			)
			w.Header().Add(
				"Vary",
				"Access-Control-Request-Headers",
			)

			// Reject requests from unknown origins.
			if origin == "" || origin != allowedOrigin {
				http.Error(
					w,
					"Origin not allowed",
					http.StatusForbidden,
				)
				return
			}

			requestedMethod := r.Header.Get(
				"Access-Control-Request-Method",
			)

			allowedMethods := map[string]bool{
				http.MethodGet:    true,
				http.MethodPost:   true,
				http.MethodPatch:  true,
				http.MethodDelete: true,
			}

			if !allowedMethods[requestedMethod] {
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			w.Header().Set(
				"Access-Control-Allow-Origin",
				allowedOrigin,
			)
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PATCH, DELETE",
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Authorization, Content-Type",
			)
			w.Header().Set(
				"Access-Control-Max-Age",
				"600",
			)

			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Allow only the configured frontend origin.
		if allowedOrigin != "" && origin == allowedOrigin {
			w.Header().Set(
				"Access-Control-Allow-Origin",
				allowedOrigin,
			)
		}

		next.ServeHTTP(w, r)
	})
}
