package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MonalBarse/referral-hub/internal/models"
)

const commentColumns = `
	c.id, c.referral_id, c.body, c.created_at,
	u.id, u.name, u.email`

func scanComment(row pgx.Row) (models.Comment, error) {
	var c models.Comment
	err := row.Scan(
		&c.ID, &c.ReferralID, &c.Body, &c.CreatedAt,
		&c.Author.ID, &c.Author.Name, &c.Author.Email,
	)
	return c, err
}

func (s *Store) CreateComment(ctx context.Context, referralID, authorID uuid.UUID, body string) (models.Comment, error) {
	const q = `
		WITH inserted AS (
			INSERT INTO comments (referral_id, author_id, body)
			VALUES ($1, $2, $3)
			RETURNING id, referral_id, author_id, body, created_at
		)
		SELECT ` + commentColumns + `
		FROM inserted c
		JOIN users u ON u.id = c.author_id`

	comment, err := scanComment(s.pool.QueryRow(ctx, q, referralID, authorID, body))
	if err != nil {
		return models.Comment{}, fmt.Errorf("create comment: %w", classify(err))
	}
	return comment, nil
}

func (s *Store) ListCommentsByReferral(ctx context.Context, referralID uuid.UUID) ([]models.Comment, error) {
	const q = `
		SELECT ` + commentColumns + `
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.referral_id = $1
		ORDER BY c.created_at ASC`

	rows, err := s.pool.Query(ctx, q, referralID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}
