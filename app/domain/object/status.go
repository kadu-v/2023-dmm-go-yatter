package object

import "time"

type Status struct {
	// The internal ID of the account
	ID int64 `json:"id,omitempty" db:"id"`

	// Account of the status
	Account *Account `json:"account"`

	// Content of the status
	Content string `json:"content" db:"content"`

	// // URL
	// Url string `json:"url" db:"url"`

	// The time the status was created
	CreateAt time.Time `json:"create_at,omitempty" db:"create_at"`
}
