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
	entity := new(object.Status)
	query := `
		SELECT
			s.id,
			s.content,
			s.create_at,
			a.id AS "account.id",
			a.username AS "account.username",
			a.display_name AS "account.display_name"
		FROM status s
		JOIN account a 
		ON s.account_id = a.id
		WHERE s.id = ?;`
	if err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find status from db: %w", err)
	}
	return entity, nil
}

// Question: DB操作でトランザクションはどのように使えば良いのか？
// Question: アカウントのテーブルを検索するけど，本来はaccount.goにこの処理を書くべきでは？
// 今回のインターン課題だとどこ？
// Question:
// FindStatusesByRange: sinceIDからmaxIDまでのstatusを取得する
func (r *status) FindStatusesByRange(ctx context.Context, sinceID int64, maxID int64, limit int64) (*object.Timeline, error) {
	query := `
		SELECT
			s.id,
			s.content,
			s.create_at,
			a.id AS "account.id",
			a.username AS "account.username",
			a.display_name AS "account.display_name"
		FROM status s
		JOIN account a 
		ON s.account_id = a.id
		WHERE ? <= s.id AND s.id <= ? 
		ORDER BY s.id ASC LIMIT ?;`

	var entities []*object.Status
	err := r.db.SelectContext(ctx, &entities, query, sinceID, maxID, limit)
	if err != nil {
		err = fmt.Errorf("failed to find statuses from since_id \"%d\" to max_id \"%d\" with limit \"%d\": %w", sinceID, maxID, limit, err)
		return nil, err
	}

	timeline := new(object.Timeline)
	timeline.Statuses = entities
	return timeline, nil
}
