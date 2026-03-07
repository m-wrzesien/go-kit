package httpkit

import (
	"log/slog"
	"net/http"
)

type Server struct {
	log *slog.Logger
	*http.Server
}

// Start launches [http.Server] from this [Server] instance
// It will panic on any errors except for [http.ErrServerClosed]
func (s Server) Start() {
	s.log.Info("App starts", "address", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// New creates instance of [Server]
func New(log *slog.Logger, s *http.Server) *Server {
	return &Server{
		log:    log,
		Server: s,
	}

}
