package dao

import (
	"context"
	"database/sql"
	"errors"
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

// FindStatusByID: 対応するIDのstatusを取得
func (r *status) FindStatusByID(ctx context.Context, id int64) (*object.Status, error) {
	status := new(object.Status)
	var accountID int64
	query := "SELECT id, account_id, content, create_at FROM status WHERE id = ?"
	if err := r.db.QueryRowxContext(ctx, query, id).Scan(&status.ID, &accountID, &status.Content, &status.CreateAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find status from db: %w", err)
	}

	account := new(object.Account)
	if err := r.db.QueryRowxContext(ctx, "SELECT * FROM account WHERE id = ?", accountID).StructScan(account); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find account from db: %w", err)
	}
	status.Account = account

	return status, nil
}
