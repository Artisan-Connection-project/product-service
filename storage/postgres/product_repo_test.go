package postgres

import (
	"context"
	"testing"
	"time"

	pro "product-service/genproto/product_service"
	"product-service/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql db, %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, func() {
		db.Close()
	}
}

func TestAddProduct(t *testing.T) {
	db, mock, teardown := setupMockDB(t)
	defer teardown()

	repo := NewProductRepo(db)
	ctx := context.Background()
	newId := uuid.NewString()
	q := 10
	p := 99.99
	product := &pro.AddProductRequest{
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       p,
		CategoryId:  "cat-123",
		ArtisanId:   "art-123",
		Quantity:    int32(q),
	}

	mock.ExpectQuery("INSERT INTO products").
		WithArgs(sqlmock.AnyArg(), product.Name, product.Description, product.Price, product.CategoryId, product.ArtisanId, product.Quantity).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "price", "category_id", "artisan_id", "quantity", "created_at"}).
			AddRow(newId, product.Name, product.Description, product.Price, product.CategoryId, product.ArtisanId, product.Quantity, time.Now()))

	response, err := repo.AddProduct(ctx, product)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, newId, response.Id)
	assert.Equal(t, product.Name, response.Name)
	assert.Equal(t, product.Description, response.Description)
	assert.Equal(t, product.Price, response.Price)
	assert.Equal(t, product.CategoryId, response.CategoryId)
	assert.Equal(t, product.ArtisanId, response.ArtisanId)
	assert.Equal(t, product.Quantity, response.Quantity)
}

func TestEditProduct(t *testing.T) {
	db, mock, teardown := setupMockDB(t)
	defer teardown()

	repo := NewProductRepo(db)
	ctx := context.Background()
	product := &pro.EditProductRequest{
		Id:          "prod-123",
		Name:        "Updated Product",
		Description: "Updated description",
		Price:       49.99,
		CategoryId:  "cat-456",
		ArtisanId:   "art-456",
		Quantity:    5,
	}

	mock.ExpectExec("UPDATE products").
		WithArgs(product.Name, product.Description, product.Price, product.CategoryId, product.ArtisanId, product.Quantity, product.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	response, err := repo.EditProduct(ctx, product)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, product.Id, response.Id)
	assert.Equal(t, product.Name, response.Name)
	assert.Equal(t, product.Description, response.Description)
	assert.Equal(t, product.Price, response.Price)
	assert.Equal(t, product.CategoryId, response.CategoryId)
	assert.Equal(t, product.ArtisanId, response.ArtisanId)
}

func TestDeleteProduct(t *testing.T) {
	db, mock, teardown := setupMockDB(t)
	defer teardown()

	repo := NewProductRepo(db)
	ctx := context.Background()
	request := &pro.DeleteProductRequest{Id: "prod-123"}

	mock.ExpectExec("UPDATE products SET deleted_at").
		WithArgs(request.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	response, err := repo.DeleteProduct(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestGetProducts(t *testing.T) {
	db, mock, teardown := setupMockDB(t)
	defer teardown()

	repo := NewProductRepo(db)
	ctx := context.Background()
	request := &pro.GetProductsRequest{Limit: "10", Page: "0"}

	mockProducts := []*pro.Product{
		{
			Id:          "prod-1",
			Name:        "Product 1",
			Description: "Description 1",
			Price:       100.0,
			CategoryId:  "cat-1",
			ArtisanId:   "art-1",
			Quantity:    10,
			CreatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			Id:          "prod-2",
			Name:        "Product 2",
			Description: "Description 2",
			Price:       200.0,
			CategoryId:  "cat-2",
			ArtisanId:   "art-2",
			Quantity:    20,
			CreatedAt:   time.Now().Format(time.RFC3339),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "category_id", "artisan_id", "quantity", "created_at"}).
		AddRow(mockProducts[0].Id, mockProducts[0].Name, mockProducts[0].Description, mockProducts[0].Price, mockProducts[0].CategoryId, mockProducts[0].ArtisanId, mockProducts[0].Quantity, mockProducts[0].CreatedAt).
		AddRow(mockProducts[1].Id, mockProducts[1].Name, mockProducts[1].Description, mockProducts[1].Price, mockProducts[1].CategoryId, mockProducts[1].ArtisanId, mockProducts[1].Quantity, mockProducts[1].CreatedAt)

	mock.ExpectQuery("SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at FROM products").
		WithArgs(request.Limit, request.Page).
		WillReturnRows(rows)

	response, err := repo.GetProducts(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Products))
	assert.Equal(t, mockProducts[0].Id, response.Products[0].Id)
	assert.Equal(t, mockProducts[1].Id, response.Products[1].Id)
}

func TestGetProduct(t *testing.T) {
	db, mock, teardown := setupMockDB(t)
	defer teardown()

	repo := NewProductRepo(db)
	ctx := context.Background()
	request := &pro.GetProductRequest{Id: "prod-123"}

	mockProduct := models.Product{
		Id:          "prod-123",
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       99.99,
		CategoryId:  "cat-123",
		ArtisanId:   "art-123",
		Quantity:    10,
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	}

	row := sqlmock.NewRows([]string{"id", "name", "description", "price", "category_id", "artisan_id", "quantity", "created_at", "updated_at"}).
		AddRow(mockProduct.Id, mockProduct.Name, mockProduct.Description, mockProduct.Price, mockProduct.CategoryId, mockProduct.ArtisanId, mockProduct.Quantity, mockProduct.CreatedAt, mockProduct.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at FROM products").
		WithArgs(request.Id).
		WillReturnRows(row)

	response, err := repo.GetProduct(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, mockProduct.Id, response.Product.Id)
	assert.Equal(t, mockProduct.Name, response.Product.Name)
	assert.Equal(t, mockProduct.Description, response.Product.Description)
	assert.Equal(t, mockProduct.Price, response.Product.Price)
	assert.Equal(t, mockProduct.CategoryId, response.Product.CategoryId)
	assert.Equal(t, mockProduct.ArtisanId, response.Product.ArtisanId)
	assert.Equal(t, mockProduct.Quantity, response.Product.Quantity)
	assert.Equal(t, mockProduct.CreatedAt, response.Product.CreatedAt)
	assert.Equal(t, mockProduct.UpdatedAt, response.Product.UpdatedAt)
}
