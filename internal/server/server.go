package server

import (
	"net/http"
	"visitor/internal/hash"
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

	return &Server{
		addr: 	addr,
		db: 	db,
		hasher: hasher,
		mux: 	mux,
	}
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr: 		s.addr,
		Handler: 	s.mux,
	}

	return srv.ListenAndServe()
}