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
	// Find the statuses from since_id to max_id constrained with limit
	FindStatusesByRange(ctx context.Context, sinceID int64, maxID int64, limit int64) (*object.Timeline, error)
}
