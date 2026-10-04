package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MonalBarse/referral-hub/internal/middleware"
	"github.com/MonalBarse/referral-hub/internal/respond"
	"github.com/MonalBarse/referral-hub/internal/ws"
)

type createCommentRequest struct {
	Body string `json:"body"`
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	referralID, ok := pathUUID(w, chi.URLParam(r, "referralID"), "referral")
	if !ok {
		return
	}

	var req createCommentRequest
	if !decode(w, r, &req) {
		return
	}

	fe := fieldErrors{}
	body := fe.required("body", req.Body, 2000)
	if !fe.ok() {
		respond.ValidationError(w, fe)
		return
	}

	comment, err := h.store.CreateComment(r.Context(), referralID, userID, body)
	if err != nil {
		fail(w, err, "referral not found", "")
		return
	}

	h.hub.Publish("referral:"+referralID.String(), ws.EventCommentCreated, comment)
	respond.JSON(w, http.StatusCreated, comment)
}

func (h *Handler) ListComments(w http.ResponseWriter, r *http.Request) {
	referralID, ok := pathUUID(w, chi.URLParam(r, "referralID"), "referral")
	if !ok {
		return
	}

	// Unknown referral must 404, not look like a thread with no comments.
	if _, err := h.store.GetReferral(r.Context(), referralID); err != nil {
		fail(w, err, "referral not found", "")
		return
	}

	comments, err := h.store.ListCommentsByReferral(r.Context(), referralID)
	if err != nil {
		fail(w, err, "", "")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"comments": comments})
}
