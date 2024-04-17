package organization

type Organization struct {
	ID   int64
	Name string
}

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
