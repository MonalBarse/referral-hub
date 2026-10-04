package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MonalBarse/referral-hub/internal/middleware"
	"github.com/MonalBarse/referral-hub/internal/models"
	"github.com/MonalBarse/referral-hub/internal/respond"
	"github.com/MonalBarse/referral-hub/internal/store"
	"github.com/MonalBarse/referral-hub/internal/ws"
)

const jobFeedLimit = 100

type createJobRequest struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserID(r.Context())

	var req createJobRequest
	if !decode(w, r, &req) {
		return
	}

	fe := fieldErrors{}
	in := store.NewJob{
		Title:       fe.required("title", req.Title, 160),
		Company:     fe.required("company", req.Company, 160),
		Description: fe.required("description", req.Description, 5000),
		Status:      fe.oneOf("status", req.Status, "open", "open", "paused", "closed"),
		PostedBy:    userID,
	}
	if !fe.ok() {
		respond.ValidationError(w, fe)
		return
	}

	job, err := h.store.CreateJob(r.Context(), in)
	if err != nil {
		fail(w, err, "poster not found", "job already exists")
		return
	}

	// Published after the write commits, and not special-cased for the author.
	h.hub.Publish(ws.TopicJobs, ws.EventJobCreated, job)
	respond.JSON(w, http.StatusCreated, job)
}

func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.store.ListJobs(r.Context(), jobFeedLimit)
	if err != nil {
		fail(w, err, "", "")
		return
	}
	respond.JSON(w, http.StatusOK, map[string]any{"jobs": jobs})
}

type jobDetailResponse struct {
	Job       models.Job        `json:"job"`
	Referrals []models.Referral `json:"referrals"`
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	jobID, ok := pathUUID(w, chi.URLParam(r, "jobID"), "job")
	if !ok {
		return
	}

	job, err := h.store.GetJob(r.Context(), jobID)
	if err != nil {
		fail(w, err, "job not found", "")
		return
	}

	referrals, err := h.store.ListReferralsByJob(r.Context(), jobID)
	if err != nil {
		fail(w, err, "", "")
		return
	}

	respond.JSON(w, http.StatusOK, jobDetailResponse{Job: job, Referrals: referrals})
}
