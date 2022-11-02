package orders

import "time"

type JsonCreateSchema struct {
	BuyerId    int `json:"buyer_id" validate:"required,min=1" db:"buyer_id"`
	ProductId  int `json:"product_id" validate:"required,min=1" db:"product_id"`
	Qty        int `json:"qty" validate:"required,min=1" db:"qty"`
	Price      int `json:"-" db:"price"`
	TotalPrice int `json:"-" db:"total_price"`
}

type Order struct {
	Id          int       `json:"id" db:"id"`
	BuyerId     int       `json:"buyer_id" db:"buyer_id"`
	BuyerName   string    `json:"buyer_name" db:"buyer_name"`
	ProductId   int       `json:"product_id" db:"product_id"`
	ProductName string    `json:"product_name" db:"product_name"`
	Qty         int       `json:"qty" db:"qty"`
	Price       int       `json:"price" db:"price"`
	TotalPrice  int       `json:"total_price" db:"total_price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
