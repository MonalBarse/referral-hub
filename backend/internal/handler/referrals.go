package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MonalBarse/referral-hub/internal/middleware"
	"github.com/MonalBarse/referral-hub/internal/respond"
	"github.com/MonalBarse/referral-hub/internal/store"
	"github.com/MonalBarse/referral-hub/internal/ws"
)

type createReferralRequest struct {
	CandidateName  string `json:"candidateName"`
	CandidateEmail string `json:"candidateEmail"`
	Message        string `json:"message"`
}

func (h *Handler) CreateReferral(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	jobID, ok := pathUUID(w, chi.URLParam(r, "jobID"), "job")
	if !ok {
		return
	}

	var req createReferralRequest
	if !decode(w, r, &req) {
		return
	}

	fe := fieldErrors{}
	in := store.NewReferral{
		JobID:          jobID,
		ReferrerID:     userID,
		CandidateName:  fe.required("candidateName", req.CandidateName, 160),
		CandidateEmail: fe.email("candidateEmail", req.CandidateEmail, 320),
		Message:        fe.optional("message", req.Message, 2000),
	}
	if !fe.ok() {
		respond.ValidationError(w, fe)
		return
	}

	referral, err := h.store.CreateReferral(r.Context(), in)
	if err != nil {
		fail(w, err, "job not found", "this candidate has already been referred for this job")
		return
	}

	// Job page sees the new referral; the dashboard bumps its referral count.
	h.hub.Publish("job:"+jobID.String(), ws.EventReferralCreated, referral)
	h.hub.Publish(ws.TopicJobs, ws.EventReferralCreated, referral)

	respond.JSON(w, http.StatusCreated, referral)
}

func (h *Handler) ListReferrals(w http.ResponseWriter, r *http.Request) {
	jobID, ok := pathUUID(w, chi.URLParam(r, "jobID"), "job")
	if !ok {
		return
	}

	referrals, err := h.store.ListReferralsByJob(r.Context(), jobID)
	if err != nil {
		fail(w, err, "", "")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"referrals": referrals})
}
