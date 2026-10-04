// Package router wires middleware and handlers onto a single chi mux.
package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MonalBarse/referral-hub/internal/auth"
	"github.com/MonalBarse/referral-hub/internal/config"
	"github.com/MonalBarse/referral-hub/internal/handler"
	"github.com/MonalBarse/referral-hub/internal/middleware"
	"github.com/MonalBarse/referral-hub/internal/respond"
	"github.com/MonalBarse/referral-hub/internal/store"
	"github.com/MonalBarse/referral-hub/internal/ws"
)

type Deps struct {
	Cfg  config.Config
	Pool *pgxpool.Pool
	Hub  *ws.Hub
}

func New(d Deps) http.Handler {
	authManager := auth.New(d.Cfg.JWTSecret)
	h := handler.New(d.Cfg, store.New(d.Pool), authManager, d.Hub)
	requireAuth := middleware.RequireAuth(authManager)

	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer) // a panic in one handler must not kill the process

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{d.Cfg.CORSOrigin},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, req *http.Request) {
		if err := d.Pool.Ping(req.Context()); err != nil {
			respond.Error(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}
		respond.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Outside /api/v1: a long-lived socket must not inherit a request timeout.
	r.Get("/ws", h.WebSocket)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(chimw.Timeout(15 * time.Second))

		r.Post("/auth/login", h.Login)

		r.Group(func(r chi.Router) {
			r.Use(requireAuth)
			r.Get("/auth/me", h.Me)
		})

		r.Get("/jobs", h.ListJobs)
		r.Get("/jobs/{jobID}", h.GetJob)
		r.Get("/jobs/{jobID}/referrals", h.ListReferrals)
		r.Get("/referrals/{referralID}/comments", h.ListComments)

		r.Group(func(r chi.Router) {
			r.Use(requireAuth)
			r.Post("/jobs", h.CreateJob)
			r.Post("/jobs/{jobID}/referrals", h.CreateReferral)
			r.Post("/referrals/{referralID}/comments", h.CreateComment)
		})
	})

	return r
}
