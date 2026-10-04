// Package models holds the API's wire types.
package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserRef struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type Job struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Company       string    `json:"company"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	PostedBy      UserRef   `json:"postedBy"`
	ReferralCount int       `json:"referralCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Referral struct {
	ID             uuid.UUID `json:"id"`
	JobID          uuid.UUID `json:"jobId"`
	CandidateName  string    `json:"candidateName"`
	CandidateEmail string    `json:"candidateEmail"`
	Message        string    `json:"message"`
	Status         string    `json:"status"`
	Referrer       UserRef   `json:"referrer"`
	CommentCount   int       `json:"commentCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Comment struct {
	ID         uuid.UUID `json:"id"`
	ReferralID uuid.UUID `json:"referralId"`
	Body       string    `json:"body"`
	Author     UserRef   `json:"author"`
	CreatedAt  time.Time `json:"createdAt"`
}
