package server

import "net/http"

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// origin := r.Header.Get("Origin")

		// Allow ONLY your admin portal
		// if origin == "https://admin.aharsuchi.com" {
		// 	w.Header().Set("Access-Control-Allow-Origin", origin)
		// 	w.Header().Set("Vary", "Origin")
		// }

		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 🔴 REQUIRED FOR PRIVATE NETWORK ACCESS
		w.Header().Set("Access-Control-Allow-Private-Network", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
