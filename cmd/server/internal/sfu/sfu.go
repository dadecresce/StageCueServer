// internal/sfu/sfu.go
package sfu

import (
	"net/http"

	"github.com/StageCue/StageCueServer/internal/config"
	"github.com/pion/sfu/v2/pkg/sfu"
	"go.uber.org/zap"
)

type Server struct {
	sfu  *sfu.SFU
	log  *zap.Logger
	conf *config.Config
}

func New(cfg *config.Config, log *zap.Logger) (*Server, error) {
	// Default SFU config; customise later via cfg
	sfuConfig := sfu.Config{}
	instance := sfu.NewSFU(sfuConfig)

	return &Server{
		sfu:  instance,
		log:  log,
		conf: cfg,
	}, nil
}

// WebSocketHandler returns the HTTP handler for signaling.
func (s *Server) WebSocketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.sfu.ServeHTTP(w, r) // Pion SFU embeds its own WS signaling
	})
}
