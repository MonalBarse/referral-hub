// Package router wires middleware and handlers onto a single chi mux. Keeping
// the whole URL surface in one file makes the API easy to review at a glance.
package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MonalBarse/referral-hub/internal/config"
	"github.com/MonalBarse/referral-hub/internal/respond"
)

// Deps collects everything the handlers need, passed in from main so nothing
// reaches for a global.
type Deps struct {
	Cfg  config.Config
	Pool *pgxpool.Pool
}

func New(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer) // a panic in one handler must not kill the process
	r.Use(chimw.Timeout(15 * time.Second))

	// The browser calls this API cross-origin from the Next.js dev server, and
	// sends an Authorization header, so CORS has to allow that header exactly.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{d.Cfg.CORSOrigin},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Liveness check: confirms the process is up AND that it can still reach
	// Postgres, which is what you actually want from a health endpoint.
	r.Get("/healthz", func(w http.ResponseWriter, req *http.Request) {
		if err := d.Pool.Ping(req.Context()); err != nil {
			respond.Error(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}
		respond.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return r
}
