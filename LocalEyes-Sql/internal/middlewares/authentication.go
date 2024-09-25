package middlewares

import (
	"encoding/json"
	"localEyes/utils"
	"net/http"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response := utils.NewUnauthorizedError("Missing authentication token")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("ERROR: Error encoding response")
			}
			return
		}
		if !utils.ValidateTokenFunc(authHeader) {
			w.WriteHeader(http.StatusUnauthorized)
			response := utils.NewUnauthorizedError("Invalid token")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("ERROR: Error encoding response")
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response := utils.NewUnauthorizedError("Missing authentication token")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("ERROR: Error encoding response")
			}
			return
		}
		if !utils.ValidateAdminTokenFunc(authHeader) {
			w.WriteHeader(http.StatusUnauthorized)
			response := utils.NewUnauthorizedError("Invalid token")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("ERROR: Error encoding response")
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
