package postgres

type ProductRepo interface {
	GetAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(product *Product) error
}
