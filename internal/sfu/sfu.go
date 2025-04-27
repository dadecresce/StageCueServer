// internal/sfu/sfu.go
package sfu

import (
	"net/http"

	"github.com/StageCue/StageCueServer/internal/config"
	"go.uber.org/zap"

	"github.com/go-logr/zapr"
	"github.com/gorilla/websocket"
	jsonrpcserver "github.com/pion/ion-sfu/cmd/signal/json-rpc/server"
	"github.com/pion/ion-sfu/pkg/middlewares/datachannel"
	"github.com/pion/ion-sfu/pkg/sfu"
	jsonrpc2 "github.com/sourcegraph/jsonrpc2"
	wsjsonrpc "github.com/sourcegraph/jsonrpc2/websocket"
)

// Server incapsula l’istanza di Ion-SFU e il logger.
type Server struct {
	s   *sfu.SFU
	log *zap.Logger
	cfg *config.Config
}

// New crea e configura l’istanza SFU (estendi sfu.Config leggendo cfg se serve).
func New(cfg *config.Config, log *zap.Logger) (*Server, error) {
	// Configurazione base SFU
	sfuConfig := sfu.Config{}
	engine := sfu.NewSFU(sfuConfig)

	// Abilita il Subscriber API via DataChannel
	dc := engine.NewDatachannel(sfu.APIChannelLabel)
	dc.Use(datachannel.SubscriberAPI)

	return &Server{s: engine, log: log, cfg: cfg}, nil
}

// WebSocketHandler ritorna l’handler HTTP per il signaling via JSON-RPC over WebSocket.
func (srv *Server) WebSocketHandler() http.Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Prepara il wrapper zap→logr
	logrLogger := zapr.NewLogger(srv.log)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP -> WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			srv.log.Error("WebSocket upgrade fallito", zap.Error(err))
			return
		}
		defer conn.Close()

		// Crea un nuovo PeerLocal per questo SFU
		peer := sfu.NewPeer(srv.s)

		// Avvolgi il WS in un transport JSON-RPC
		stream := wsjsonrpc.NewObjectStream(conn)
		handler := jsonrpcserver.NewJSONSignal(peer, logrLogger)

		// Avvia la connessione JSON-RPC 2.0
		rpcConn := jsonrpc2.NewConn(r.Context(), stream, handler)
		<-rpcConn.DisconnectNotify()
	})
}
