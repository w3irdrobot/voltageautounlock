package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		log.Fatalf("error getting the config: %s", err)
	}

	router := http.NewServeMux()
	server := &http.Server{Addr: cfg.address, Handler: router}
	appCtx, appCancel := context.WithCancel(context.Background())

	router.Handle("/unlock", newUnlockHandler(appCtx, cfg.webhookSecret, cfg.nodeAPI, cfg.walletPassword))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	shutdownComplete := make(chan struct{}, 1)

	go func() {
		log.Println("shutting down the server")
		<-signals

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("error shutting down the server: %s", err)
		}

		close(shutdownComplete)
	}()

	log.Printf("starting the server on %s", cfg.address)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error starting the server: %s", err)
	}

	appCancel()
	<-shutdownComplete
	log.Println("server shutdown successfully")
}
