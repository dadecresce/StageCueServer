// internal/sfu/sfu.go
package sfu

import (
	"net/http"

	"github.com/StageCue/StageCueServer/internal/config"
	"go.uber.org/zap"
)

// Server gestisce lo signaling SFU (placeholder).
type Server struct {
	log  *zap.Logger
	conf *config.Config
}

// New crea un’istanza di Server (sostituire con init reale SFU).
func New(cfg *config.Config, log *zap.Logger) (*Server, error) {
	return &Server{log: log, conf: cfg}, nil
}

// WebSocketHandler restituisce l’handler HTTP per lo signaling.
// Attualmente risponde 501 Not Implemented.
func (s *Server) WebSocketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Warn("SFU signaling not implemented", zap.String("path", r.URL.Path))
		http.Error(w, "SFU signaling not implemented", http.StatusNotImplemented)
	})
}
