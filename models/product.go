package models

type Product struct {
	Id          string  `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	CategoryId  string  `json:"category_id" db:"category_id"`
	ArtisanId   string  `json:"artisan_id" db:"artisan_id"`
	Quantity    int32   `json:"quantity" db:"quantity"`
	CreatedAt   string  `json:"created_at" db:"created_at"`
	UpdatedAt   string  `json:"updated_at" db:"updated_at" `
}

type Order struct {
	Id              string  `db:"id" json:"id"`
	UserId          string  `db:"user_id" json:"user_id"`
	TotalAmount     float64 `db:"total_amount" json:"total_amount"`
	Status          string  `db:"status" json:"status"`
	ShippingAddress string  `db:"shipping_address" json:"shipping_address"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`
}
