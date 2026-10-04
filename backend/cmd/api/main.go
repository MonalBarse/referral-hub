// Command api is the referral-hub HTTP + WebSocket server.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MonalBarse/referral-hub/internal/config"
	"github.com/MonalBarse/referral-hub/internal/db"
	"github.com/MonalBarse/referral-hub/internal/router"
	"github.com/MonalBarse/referral-hub/internal/ws"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool, "migrations"); err != nil {
		log.Fatalf("startup: migrate: %v", err)
	}

	hub := ws.NewHub()

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.New(router.Deps{Cfg: *cfg, Pool: pool, Hub: hub}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("api: listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("api: listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("api: shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("api: shutdown: %v", err)
	}
}
