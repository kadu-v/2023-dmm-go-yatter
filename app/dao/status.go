package dao

import (
	"context"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Status
	status struct {
		db *sqlx.DB
	}
)

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// AddStatus: statusをpost
func (r *status) AddStatus(ctx context.Context, a *object.Account, s *object.Status) (int64, error) {
	query := "INSERT INTO status (account_id, content) VALUES (?, ?);"
	result, err := r.db.Exec(query, a.ID, s.Content)
	if err != nil {
		return -1, fmt.Errorf("failed to post a status: %w", err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to get the last id: %w", err)
	}
	return lastID, nil
}
