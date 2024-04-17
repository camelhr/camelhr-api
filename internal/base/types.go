package base

import "time"

// Timestamps is a struct that contains the timestamps for a record.
// It can be embedded to db model and response structs to add the timestamps fields.
type Timestamps struct {
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
