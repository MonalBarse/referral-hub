package handler

import (
	"net/http"

	"github.com/MonalBarse/referral-hub/internal/respond"
)

// WebSocket authenticates from a query parameter, because the browser
// WebSocket API cannot set headers on the handshake.
func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respond.Error(w, http.StatusUnauthorized, "missing token")
		return
	}

	userID, err := h.auth.Parse(token)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	h.hub.Serve(w, r, userID, func(req *http.Request) bool {
		origin := req.Header.Get("Origin")
		return origin == "" || origin == h.cfg.CORSOrigin
	})
}
