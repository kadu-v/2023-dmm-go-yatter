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
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create account repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find account from db: %w", err)
	}

	return entity, nil
}

// AddAccount: 新規アカウントを登録
func (r *account) AddAccount(ctx context.Context, a *object.Account) error {
	query := "INSERT INTO account (username, password_hash, display_name, avatar, header, note) VALUES (?, ?, ?, NULL, NULL, NULL)"
	_, err := r.db.Exec(query, a.Username, a.PasswordHash, a.Username)
	if err != nil {
		return fmt.Errorf("failed to register new account to db: %w", err)
	}
	return nil
}
