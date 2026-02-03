package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"visitor/internal/hash"
	"visitor/internal/model"
	"visitor/internal/storage"
)

type Server struct {
	addr	string
	db 		*storage.DB
	hasher 	*hash.Manager
	mux 	*http.ServeMux
}

func New(addr string, db *storage.DB, hasher *hash.Manager) *Server {
	mux := http.NewServeMux()

	s := &Server{
		addr: 	addr,
		db: 	db,
		hasher: hasher,
		mux: 	mux,
	}

	s.mux.HandleFunc("POST /api/event", s.handleEvent)

	return s
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr: 		s.addr,
		Handler: 	cors(s.mux),
	}

	return srv.ListenAndServe()
}

func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request) {
	var event model.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if event.Domain == "" {
		http.Error(w, "Domain is required", http.StatusBadRequest)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	userAgent := r.Header.Get("User-Agent")

	visitorHash, err := s.hasher.GetHash(r.Context(), event.Domain, ip, userAgent)
	if err != nil {
		log.Printf("failed to get visitor hash: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	pv := &model.PageView{
		Domain: 		event.Domain,
		Path: 			event.Path,
		Referrer: 		event.Referrer,
		VisitorHash: 	visitorHash,
	}

	if err := s.db.InsertPageView(r.Context(), pv); err != nil {
		log.Printf("Failed to insert page view: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}