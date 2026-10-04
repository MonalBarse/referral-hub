package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MonalBarse/referral-hub/internal/models"
)

type NewJob struct {
	Title       string
	Company     string
	Description string
	Status      string
	PostedBy    uuid.UUID
}

const jobColumns = `
	j.id, j.title, j.company, j.description, j.status, j.created_at,
	u.id, u.name, u.email,
	(SELECT count(*) FROM referrals r WHERE r.job_id = j.id)`

func scanJob(row pgx.Row) (models.Job, error) {
	var j models.Job
	err := row.Scan(
		&j.ID, &j.Title, &j.Company, &j.Description, &j.Status, &j.CreatedAt,
		&j.PostedBy.ID, &j.PostedBy.Name, &j.PostedBy.Email,
		&j.ReferralCount,
	)
	return j, err
}

func (s *Store) CreateJob(ctx context.Context, in NewJob) (models.Job, error) {
	const q = `
		WITH inserted AS (
			INSERT INTO jobs (title, company, description, status, posted_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, title, company, description, status, created_at, posted_by
		)
		SELECT ` + jobColumns + `
		FROM inserted j
		JOIN users u ON u.id = j.posted_by`

	job, err := scanJob(s.pool.QueryRow(ctx, q,
		in.Title, in.Company, in.Description, in.Status, in.PostedBy))
	if err != nil {
		return models.Job{}, fmt.Errorf("create job: %w", classify(err))
	}
	return job, nil
}

func (s *Store) ListJobs(ctx context.Context, limit int) ([]models.Job, error) {
	const q = `
		SELECT ` + jobColumns + `
		FROM jobs j
		JOIN users u ON u.id = j.posted_by
		ORDER BY j.created_at DESC
		LIMIT $1`

	rows, err := s.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	jobs := []models.Job{}
	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, fmt.Errorf("scan job: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (s *Store) GetJob(ctx context.Context, id uuid.UUID) (models.Job, error) {
	const q = `
		SELECT ` + jobColumns + `
		FROM jobs j
		JOIN users u ON u.id = j.posted_by
		WHERE j.id = $1`

	job, err := scanJob(s.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Job{}, ErrNotFound
	}
	if err != nil {
		return models.Job{}, fmt.Errorf("get job: %w", err)
	}
	return job, nil
}
