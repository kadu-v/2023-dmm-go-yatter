package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Add new status
	AddStatus(ctx context.Context, a *object.Account, s *object.Status) (int64, error)
	// Find the status corresponding with ID
	FindStatusByID(ctx context.Context, id int64) (*object.Status, error)
}
