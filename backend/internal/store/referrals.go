package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MonalBarse/referral-hub/internal/models"
)

type NewReferral struct {
	JobID          uuid.UUID
	ReferrerID     uuid.UUID
	CandidateName  string
	CandidateEmail string
	Message        string
}

const referralColumns = `
	r.id, r.job_id, r.candidate_name, r.candidate_email, r.message, r.status, r.created_at,
	u.id, u.name, u.email,
	(SELECT count(*) FROM comments c WHERE c.referral_id = r.id)`

func scanReferral(row pgx.Row) (models.Referral, error) {
	var r models.Referral
	err := row.Scan(
		&r.ID, &r.JobID, &r.CandidateName, &r.CandidateEmail, &r.Message, &r.Status, &r.CreatedAt,
		&r.Referrer.ID, &r.Referrer.Name, &r.Referrer.Email,
		&r.CommentCount,
	)
	return r, err
}

func (s *Store) CreateReferral(ctx context.Context, in NewReferral) (models.Referral, error) {
	const q = `
		WITH inserted AS (
			INSERT INTO referrals (job_id, referrer_id, candidate_name, candidate_email, message)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, job_id, referrer_id, candidate_name, candidate_email, message, status, created_at
		)
		SELECT ` + referralColumns + `
		FROM inserted r
		JOIN users u ON u.id = r.referrer_id`

	ref, err := scanReferral(s.pool.QueryRow(ctx, q,
		in.JobID, in.ReferrerID, in.CandidateName, in.CandidateEmail, in.Message))
	if err != nil {
		return models.Referral{}, fmt.Errorf("create referral: %w", classify(err))
	}
	return ref, nil
}

func (s *Store) ListReferralsByJob(ctx context.Context, jobID uuid.UUID) ([]models.Referral, error) {
	const q = `
		SELECT ` + referralColumns + `
		FROM referrals r
		JOIN users u ON u.id = r.referrer_id
		WHERE r.job_id = $1
		ORDER BY r.created_at DESC`

	rows, err := s.pool.Query(ctx, q, jobID)
	if err != nil {
		return nil, fmt.Errorf("list referrals: %w", err)
	}
	defer rows.Close()

	referrals := []models.Referral{}
	for rows.Next() {
		ref, err := scanReferral(rows)
		if err != nil {
			return nil, fmt.Errorf("scan referral: %w", err)
		}
		referrals = append(referrals, ref)
	}
	return referrals, rows.Err()
}

func (s *Store) GetReferral(ctx context.Context, id uuid.UUID) (models.Referral, error) {
	const q = `
		SELECT ` + referralColumns + `
		FROM referrals r
		JOIN users u ON u.id = r.referrer_id
		WHERE r.id = $1`

	ref, err := scanReferral(s.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Referral{}, ErrNotFound
	}
	if err != nil {
		return models.Referral{}, fmt.Errorf("get referral: %w", err)
	}
	return ref, nil
}
