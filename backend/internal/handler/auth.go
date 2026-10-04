package handler

import (
	"net/http"

	"github.com/MonalBarse/referral-hub/internal/middleware"
	"github.com/MonalBarse/referral-hub/internal/models"
	"github.com/MonalBarse/referral-hub/internal/respond"
)

type loginRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type loginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Login is the mocked auth: an email identifies the account, no password.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decode(w, r, &req) {
		return
	}

	fe := fieldErrors{}
	email := fe.email("email", req.Email, 320)
	name := fe.required("name", req.Name, 120)
	if !fe.ok() {
		respond.ValidationError(w, fe)
		return
	}

	user, err := h.store.UpsertByEmail(r.Context(), email, name)
	if err != nil {
		fail(w, err, "user not found", "email already registered")
		return
	}

	token, err := h.auth.Issue(user.ID)
	if err != nil {
		fail(w, err, "", "")
		return
	}

	respond.JSON(w, http.StatusOK, loginResponse{Token: token, User: user})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.store.GetUser(r.Context(), userID)
	if err != nil {
		fail(w, err, "user not found", "")
		return
	}
	respond.JSON(w, http.StatusOK, user)
}
