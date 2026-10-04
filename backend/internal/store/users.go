package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MonalBarse/referral-hub/internal/models"
)

func (s *Store) UpsertByEmail(ctx context.Context, email, name string) (models.User, error) {
	const q = `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, email, name, created_at`

	var u models.User
	err := s.pool.QueryRow(ctx, q, email, name).
		Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("upsert user: %w", classify(err))
	}
	return u, nil
}

func (s *Store) GetUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	const q = `SELECT id, email, name, created_at FROM users WHERE id = $1`

	var u models.User
	err := s.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}
