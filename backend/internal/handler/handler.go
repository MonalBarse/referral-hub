// Package handler holds the HTTP handlers: decode, validate, store, publish.
package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/MonalBarse/referral-hub/internal/auth"
	"github.com/MonalBarse/referral-hub/internal/config"
	"github.com/MonalBarse/referral-hub/internal/respond"
	"github.com/MonalBarse/referral-hub/internal/store"
	"github.com/MonalBarse/referral-hub/internal/ws"
)

type Handler struct {
	cfg   config.Config
	store *store.Store
	auth  *auth.Manager
	hub   *ws.Hub
}

func New(cfg config.Config, s *store.Store, a *auth.Manager, hub *ws.Hub) *Handler {
	return &Handler{cfg: cfg, store: s, auth: a, hub: hub}
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	defer r.Body.Close()
	limited := io.LimitReader(r.Body, 64<<10)

	dec := json.NewDecoder(limited)
	dec.DisallowUnknownFields() // typos in field names fail loudly, not silently

	if err := dec.Decode(dst); err != nil {
		respond.Error(w, http.StatusBadRequest, "request body must be valid JSON")
		return false
	}
	return true
}

func pathUUID(w http.ResponseWriter, raw, label string) (uuid.UUID, bool) {
	id, err := uuid.Parse(raw)
	if err != nil {
		respond.Error(w, http.StatusNotFound, label+" not found")
		return uuid.Nil, false
	}
	return id, true
}

func fail(w http.ResponseWriter, err error, notFoundMsg, duplicateMsg string) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		respond.Error(w, http.StatusNotFound, notFoundMsg)
	case errors.Is(err, store.ErrDuplicate):
		respond.Error(w, http.StatusConflict, duplicateMsg)
	default:
		log.Printf("handler: %v", err)
		respond.Error(w, http.StatusInternalServerError, "something went wrong")
	}
}
