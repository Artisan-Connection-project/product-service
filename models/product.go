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
