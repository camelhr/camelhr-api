package organization

import "github.com/camelhr/camelhr-api/internal/base"

type Organization struct {
	ID   int64  `db:"organization_id"`
	Name string `db:"name"`
	base.Timestamps
}

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	base.Timestamps
}
