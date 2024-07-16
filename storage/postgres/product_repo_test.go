package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	pro "product-service/genproto/product_service"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB *sqlx.DB
	repo   ProductRepo
)

func setupTestDB() (*sqlx.DB, error) {
	connStr := "user=postgres port=5432 password=1702 host=localhost dbname=product_service sslmode=disable" // Adjust as needed
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			category_id UUID,
			artisan_id UUID,
			quantity INTEGER NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestMain(m *testing.M) {
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		panic(err)
	}

	repo = NewProductRepo(testDB)

	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func TestAddProduct(t *testing.T) {
	product := &pro.AddProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		CategoryId:  "af97dcdb-9431-4f85-b98b-5bb7366d4428",
		ArtisanId:   uuid.NewString(),
		Quantity:    10,
	}

	resp, err := repo.AddProduct(context.Background(), product)
	require.NoError(t, err)

	var dbProduct pro.Product
	err = testDB.Get(&dbProduct, "SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at FROM products WHERE id = $1", resp.Id)
	require.NoError(t, err)

	assert.Equal(t, product.Name, dbProduct.Name)
	assert.Equal(t, product.Description, dbProduct.Description)
	assert.Equal(t, product.Price, dbProduct.Price)
	assert.Equal(t, product.CategoryId, dbProduct.CategoryId)
	assert.Equal(t, product.ArtisanId, dbProduct.ArtisanId)
	assert.Equal(t, product.Quantity, dbProduct.Quantity)
}

func TestEditProduct(t *testing.T) {
	productId := uuid.NewString()
	_, err := testDB.Exec(`
		INSERT INTO products (id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		productId, "Old Name", "Old Description", 50.00, uuid.NewString(), uuid.NewString(), 5, time.Now(), time.Now(),
	)
	require.NoError(t, err)

	product := &pro.EditProductRequest{
		Id:          productId,
		Name:        "Updated Name",
		Description: "Updated Description",
		Price:       89.99,
		CategoryId:  uuid.NewString(),
		ArtisanId:   uuid.NewString(),
		Quantity:    15,
	}

	_, err = repo.EditProduct(context.Background(), product)
	require.NoError(t, err)

	var dbProduct pro.Product
	err = testDB.Get(&dbProduct, "SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at FROM products WHERE id = $1", productId)
	require.NoError(t, err)

	assert.Equal(t, product.Name, dbProduct.Name)
	assert.Equal(t, product.Description, dbProduct.Description)
	assert.Equal(t, product.Price, dbProduct.Price)
	assert.Equal(t, product.CategoryId, dbProduct.CategoryId)
	assert.Equal(t, product.ArtisanId, dbProduct.ArtisanId)
	assert.Equal(t, product.Quantity, dbProduct.Quantity)
}

func TestDeleteProduct(t *testing.T) {
	productId := uuid.NewString()
	_, err := testDB.Exec(`
		INSERT INTO products (id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		productId, "Product to Delete", "Description", 20.00, uuid.NewString(), uuid.NewString(), 10, time.Now(), time.Now(),
	)
	require.NoError(t, err)

	_, err = repo.DeleteProduct(context.Background(), &pro.DeleteProductRequest{Id: productId})
	require.NoError(t, err)

	var dbProduct pro.Product
	err = testDB.Get(&dbProduct, "SELECT deleted_at FROM products WHERE id = $1", productId)
	require.NoError(t, err)

	assert.NotNil(t, dbProduct.UpdatedAt)
}

func TestGetProducts(t *testing.T) {
	_, err := testDB.Exec(`
		INSERT INTO products (id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9),
		       ($10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		uuid.NewString(), "Product1", "Description1", 30.00, uuid.NewString(), uuid.NewString(), 10, time.Now(), time.Now(),
		uuid.NewString(), "Product2", "Description2", 40.00, uuid.NewString(), uuid.NewString(), 20, time.Now(), time.Now(),
	)
	require.NoError(t, err)

	resp, err := repo.GetProducts(context.Background(), &pro.GetProductsRequest{Limit: "10", Page: "0"})
	require.NoError(t, err)
	require.Len(t, resp.Products, 2)
}

func TestGetProduct(t *testing.T) {
	productId := uuid.NewString()
	_, err := testDB.Exec(`
		INSERT INTO products (id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		productId, "Single Product", "Single Description", 25.00, uuid.NewString(), uuid.NewString(), 15, time.Now(), time.Now(),
	)
	require.NoError(t, err)

	resp, err := repo.GetProduct(context.Background(), &pro.GetProductRequest{Id: productId})
	require.NoError(t, err)

	assert.Equal(t, productId, resp.Product.Id)
	assert.Equal(t, "Single Product", resp.Product.Name)
	assert.Equal(t, "Single Description", resp.Product.Description)
	assert.Equal(t, 25.00, resp.Product.Price)
}
func TestAddProductCategory(t *testing.T) {
	request := &pro.AddProductCategoryRequest{
		Name:        "Test Category",
		Description: "Test Description",
	}

	resp, err := repo.AddProductCategory(context.Background(), request)
	require.NoError(t, err)

	var dbCategory struct {
		Id          string    `db:"id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}

	err = testDB.Get(&dbCategory, "SELECT id, name, description, created_at, updated_at FROM product_categories WHERE id = $1", resp.Id)
	require.NoError(t, err)

	assert.Equal(t, request.Name, dbCategory.Name)
	assert.Equal(t, request.Description, dbCategory.Description)
	assert.WithinDuration(t, time.Now(), dbCategory.CreatedAt, 1*time.Second)
	assert.WithinDuration(t, time.Now(), dbCategory.UpdatedAt, 1*time.Second)
}
