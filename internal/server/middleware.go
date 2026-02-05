package server

import "net/http"

func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			for domain := range s.allowedDomains {
				if origin == "http://"+domain || origin == "https://"+domain {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return 
		}

		next.ServeHTTP(w, r)
	}))
}

func (s *Server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.password != "" {
			_, pass, ok := r.BasicAuth()
			if !ok || pass != s.password {
				w.Header().Set("WWW-Authenticate", `Basic realm="visitor"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
