package organization

import (
	"time"

	"github.com/camelhr/camelhr-api/internal/base"
)

type Organization struct {
	ID          int64      `db:"organization_id"`
	Name        string     `db:"name"`
	SuspendedAt *time.Time `db:"suspended_at"`
	base.Timestamps
}

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	SuspendedAt *time.Time `json:"suspended_at"`
	base.Timestamps
}
