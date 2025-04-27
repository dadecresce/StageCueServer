// cmd/server/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StageCue/StageCueServer/internal/config"
	"github.com/StageCue/StageCueServer/internal/logger"
	"github.com/StageCue/StageCueServer/internal/metrics"
	"github.com/StageCue/StageCueServer/internal/sfu"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	var (
		configPath string
	)
	flag.StringVar(&configPath, "config", "config.toml", "config file path")
	flag.Parse()

	cfg, err := config.Parse(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	log.Info("StageCueServer starting",
		zap.String("addr", cfg.Address),
		zap.String("config", configPath),
		zap.String("version", version))

	// ---- SFU ---------------------------------------------------------------
	sfuInstance, err := sfu.New(cfg, log)
	if err != nil {
		log.Fatal("failed to init SFU", zap.Error(err))
	}
	// -----------------------------------------------------------------------

	metrics.MustRegisterDefault() // Prom metrics

	mux := http.NewServeMux()
	mux.Handle("/ws", sfuInstance.WebSocketHandler()) // signaling
	mux.Handle("/metrics", promhttp.Handler())        // Prometheus
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "StageCueServer %s", version)
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error", zap.Error(err))
		}
	}()
	log.Info("listening...", zap.String("addr", cfg.Address))

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", zap.Error(err))
	}
	log.Sync()
}
