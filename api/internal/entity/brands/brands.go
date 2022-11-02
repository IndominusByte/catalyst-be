package brands

import "time"

type JsonCreateSchema struct {
	Name string `json:"name" validate:"required,min=3,max=100" db:"name"`
}

type Brand struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
