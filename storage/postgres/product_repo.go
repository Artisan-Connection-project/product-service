package postgres

import (
	"context"
	"database/sql"
	pro "product-service/genproto/product_service"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	AddProduct(c context.Context, product *pro.AddProductRequest) (*pro.AddProductResponse, error)
	EditProduct(c context.Context, product *pro.EditProductRequest) (*pro.EditProductResponse, error)
	DeleteProduct(c context.Context, request *pro.DeleteProductRequest) (*pro.DeleteProductResponse, error)
	GetProducts(c context.Context, request *pro.GetProductsRequest) (*pro.GetProductsResponse, error)
	GetProduct(c context.Context, request *pro.GetProductRequest) (*pro.GetProductResponse, error)
	SearchProducts(c context.Context, request *pro.SearchProductsRequest) (*pro.SearchProductsResponse, error)
	AddRating(c context.Context, rating *pro.AddRatingRequest) (*pro.AddRatingResponse, error)
	GetRatings(c context.Context, request *pro.GetRatingsRequest) (*pro.GetRatingsResponse, error)
	PlaceOrder(c context.Context, order *pro.PlaceOrderRequest) (*pro.PlaceOrderResponse, error)
	CancelOrder(c context.Context, request *pro.CancelOrderRequest) (*pro.CancelOrderResponse, error)
	UpdateOrderStatus(c context.Context, request *pro.UpdateOrderStatusRequest) (*pro.UpdateOrderStatusResponse, error)
	GetOrders(c context.Context, request *pro.GetOrdersRequest) (*pro.GetOrdersResponse, error)
	GetOrder(c context.Context, request *pro.GetOrderRequest) (*pro.GetOrderResponse, error)
	PayOrder(c context.Context, request *pro.PayOrderRequest) (*pro.PayOrderResponse, error)
	CheckPaymentStatus(c context.Context, request *pro.CheckPaymentStatusRequest) (*pro.CheckPaymentStatusResponse, error)
	UpdateShippingInfo(c context.Context, request *pro.UpdateShippingInfoRequest) (*pro.UpdateShippingInfoResponse, error)
	AddArtisanCategory(c context.Context, request *pro.AddArtisanCategoryRequest) (*pro.AddArtisanCategoryResponse, error)
	AddProductCategory(c context.Context, request *pro.AddProductCategoryRequest) (*pro.AddProductCategoryResponse, error)
	GetStatistics(c context.Context, request *pro.GetStatisticsRequest) (*pro.GetStatisticsResponse, error)
	GetUserActivity(c context.Context, request *pro.GetUserActivityRequest) (*pro.GetUserActivityResponse, error)
	GetRecommendations(c context.Context, request *pro.GetRecommendationsRequest) (*pro.GetRecommendationsResponse, error)
	GetArtisanRankings(c context.Context, request *pro.GetArtisanRankingsRequest) (*pro.GetArtisanRankingsResponse, error)
}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) AddProduct(c context.Context, product *pro.AddProductRequest) (*pro.AddProductResponse, error) {
	newId := uuid.NewString()
	query := `
		INSERT INTO products (id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(c, query,
		newId,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		product.ArtisanId,
		product.Quantity,
		time.Now(),
		time.Now())

	if err != nil {
		return nil, err
	}

	return &pro.AddProductResponse{Id: newId}, nil
}

func (r *productRepo) EditProduct(c context.Context, product *pro.EditProductRequest) (*pro.EditProductResponse, error) {
	query := `
		UPDATE products SET name = $1, description = $2, price = $3, category_id = $4, artisan_id = $5, quantity = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(c, query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		product.ArtisanId,
		product.Quantity,
		product.Id)

	if err != nil {
		return nil, err
	}

	return &pro.EditProductResponse{}, nil
}

func (r *productRepo) DeleteProduct(c context.Context, request *pro.DeleteProductRequest) (*pro.DeleteProductResponse, error) {
	query := `
		UPDATE products SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(c, query, request.Id)
	if err != nil {
		return nil, err
	}

	return &pro.DeleteProductResponse{}, nil
}

func (r *productRepo) GetProducts(c context.Context, request *pro.GetProductsRequest) (*pro.GetProductsResponse, error) {
	var products []*pro.Product
	query := `
		SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at 
		FROM products 
		WHERE deleted_at IS NULL 
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(c, &products, query, request.Limit, request.Page)
	if err != nil {
		return nil, err
	}

	return &pro.GetProductsResponse{Products: products}, nil
}

func (r *productRepo) GetProduct(c context.Context, request *pro.GetProductRequest) (*pro.GetProductResponse, error) {
	var product pro.Product
	err := r.db.GetContext(c, &product, `
		SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at 
		FROM products 
		WHERE id = $1 AND deleted_at IS NULL
	`, request.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Product not found
		}
		return nil, err
	}

	return &pro.GetProductResponse{Product: &product}, nil
}

func (r *productRepo) SearchProducts(c context.Context, request *pro.SearchProductsRequest) (*pro.SearchProductsResponse, error) {
	var products []*pro.Product
	query := `
		SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at 
		FROM products 
		WHERE deleted_at IS NULL AND name ILIKE $1
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(c, &products, query, "%"+request.Query+"%", request.Limit, request.Page)
	if err != nil {
		return nil, err
	}

	return &pro.SearchProductsResponse{Products: products}, nil
}

func (r *productRepo) AddRating(c context.Context, rating *pro.AddRatingRequest) (*pro.AddRatingResponse, error) {
	newId := uuid.NewString()
	query := `
		INSERT INTO ratings (id, product_id, user_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(c, query,
		newId,
		rating.ProductId,
		rating.UserId,
		rating.Rating,
		rating.Comment,
		time.Now())

	if err != nil {
		return nil, err
	}

	return &pro.AddRatingResponse{Id: newId}, nil
}

func (r *productRepo) GetRatings(c context.Context, request *pro.GetRatingsRequest) (*pro.GetRatingsResponse, error) {
	var ratings []*pro.Rating
	query := `
		SELECT id, product_id, user_id, rating, comment, created_at 
		FROM ratings 
		WHERE product_id = $1
	`
	err := r.db.SelectContext(c, &ratings, query, request.ProductId)
	if err != nil {
		return nil, err
	}

	return &pro.GetRatingsResponse{Ratings: ratings}, nil
}

func (r *productRepo) PlaceOrder(c context.Context, order *pro.PlaceOrderRequest) (*pro.PlaceOrderResponse, error) {
	newId := uuid.NewString()
	query := `
		INSERT INTO orders (id, user_id, status, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(c, query,
		newId,
		order.UserId,
		order.Status,
		order.TotalAmount,
		time.Now(),
		time.Now())

	if err != nil {
		return nil, err
	}

	return &pro.PlaceOrderResponse{Id: newId}, nil
}

func (r *productRepo) CancelOrder(c context.Context, request *pro.CancelOrderRequest) (*pro.CancelOrderResponse, error) {
	query := `
		UPDATE orders SET status = 'cancelled', updated_at = now() WHERE id = $1 AND status = 'pending'
	`
	_, err := r.db.ExecContext(c, query, request.Id)
	if err != nil {
		return nil, err
	}

	return &pro.CancelOrderResponse{}, nil
}

func (r *productRepo) UpdateOrderStatus(c context.Context, request *pro.UpdateOrderStatusRequest) (*pro.UpdateOrderStatusResponse, error) {
	query := `
		UPDATE orders SET status = $1, updated_at = now() WHERE id = $2
	`
	_, err := r.db.ExecContext(c, query, request.Status, request.Id)
	if err != nil {
		return nil, err
	}

	return &pro.UpdateOrderStatusResponse{}, nil
}

// GetOrders retrieves a list of orders with pagination
func (r *productRepo) GetOrders(c context.Context, request *pro.GetOrdersRequest) (*pro.GetOrdersResponse, error) {
	var orders []*pro.Order
	query := `
		SELECT id, user_id, status, total_amount, created_at, updated_at 
		FROM orders 
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(c, &orders, query, request.Limit, request.Page)
	if err != nil {
		return nil, err
	}

	return &pro.GetOrdersResponse{Orders: orders}, nil
}

func (r *productRepo) GetOrder(c context.Context, request *pro.GetOrderRequest) (*pro.GetOrderResponse, error) {
	var order pro.Order
	err := r.db.GetContext(c, &order, `
		SELECT id, user_id, status, total_amount, created_at, updated_at 
		FROM orders 
		WHERE id = $1 AND deleted_at IS NULL
	`, request.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Order not found
		}
		return nil, err
	}

	return &pro.GetOrderResponse{Order: &order}, nil
}

func (r *productRepo) PayOrder(c context.Context, request *pro.PayOrderRequest) (*pro.PayOrderResponse, error) {
	query := `
		UPDATE orders SET status = 'paid', updated_at = now() WHERE id = $1 AND status = 'pending'
	`
	_, err := r.db.ExecContext(c, query, request.OrderId)
	if err != nil {
		return nil, err
	}

	return &pro.PayOrderResponse{}, nil
}

func (r *productRepo) CheckPaymentStatus(c context.Context, request *pro.CheckPaymentStatusRequest) (*pro.CheckPaymentStatusResponse, error) {
	var order pro.Order
	err := r.db.GetContext(c, &order, `
		SELECT status 
		FROM orders 
		WHERE id = $1
	`, request.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Order not found
		}
		return nil, err
	}

	return &pro.CheckPaymentStatusResponse{}, nil
}

func (r *productRepo) UpdateShippingInfo(c context.Context, request *pro.UpdateShippingInfoRequest) (*pro.UpdateShippingInfoResponse, error) {
	query := `
		UPDATE orders SET shipping_address = $1, updated_at = now() WHERE id = $2
	`
	_, err := r.db.ExecContext(c, query, request.ShippingAddress, request.Id)
	if err != nil {
		return nil, err
	}

	return &pro.UpdateShippingInfoResponse{}, nil
}

func (r *productRepo) AddArtisanCategory(c context.Context, request *pro.AddArtisanCategoryRequest) (*pro.AddArtisanCategoryResponse, error) {
	newId := uuid.NewString()
	query := `
		INSERT INTO artisan_categories (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(c, query,
		newId,
		request.Name,
		time.Now(),
		time.Now())

	if err != nil {
		return nil, err
	}

	return &pro.AddArtisanCategoryResponse{Id: newId}, nil
}

func (r *productRepo) AddProductCategory(c context.Context, request *pro.AddProductCategoryRequest) (*pro.AddProductCategoryResponse, error) {
	newId := uuid.NewString()
	query := `
		INSERT INTO product_categories (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(c, query,
		newId,
		request.Name,
		time.Now(),
		time.Now())

	if err != nil {
		return nil, err
	}

	return &pro.AddProductCategoryResponse{Id: newId}, nil
}

func (r *productRepo) GetStatistics(c context.Context, request *pro.GetStatisticsRequest) (*pro.GetStatisticsResponse, error) {
	var totalProducts, totalOrders int
	err := r.db.QueryRowContext(c, `
		SELECT (SELECT COUNT(*) FROM products WHERE deleted_at IS NULL),
		       (SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL)
	`).Scan(&totalProducts, &totalOrders)
	if err != nil {
		return nil, err
	}

	return &pro.GetStatisticsResponse{}, nil
}

func (r *productRepo) GetUserActivity(c context.Context, request *pro.GetUserActivityRequest) (*pro.GetUserActivityResponse, error) {
	var activities []*pro.UserActivity
	query := `
		SELECT id, user_id, action, created_at 
		FROM user_activities 
		WHERE user_id = $1
	`
	err := r.db.SelectContext(c, &activities, query, request)
	if err != nil {
		return nil, err
	}

	return &pro.GetUserActivityResponse{Activities: activities}, nil
}

func (r *productRepo) GetRecommendations(c context.Context, request *pro.GetRecommendationsRequest) (*pro.GetRecommendationsResponse, error) {
	var recommendations []*pro.Product
	query := `
		SELECT p.id, p.name, p.description, p.price, p.category_id, p.artisan_id, p.quantity, p.created_at, p.updated_at
		FROM products p
		JOIN user_preferences up ON p.category_id = up.preferred_category
		WHERE up.user_id = $1
		LIMIT $2
	`
	err := r.db.SelectContext(c, &recommendations, query, request, request)
	if err != nil {
		return nil, err
	}

	return &pro.GetRecommendationsResponse{}, nil
}

func (r *productRepo) GetArtisanRankings(c context.Context, request *pro.GetArtisanRankingsRequest) (*pro.GetArtisanRankingsResponse, error) {
	var rankings []*pro.ArtisanRanking
	query := `
		SELECT a.id, a.name, COUNT(o.id) as order_count
		FROM artisans a
		JOIN orders o ON a.id = o.artisan_id
		GROUP BY a.id, a.name
		ORDER BY order_count DESC
		LIMIT $1
	`
	err := r.db.SelectContext(c, &rankings, query, request)
	if err != nil {
		return nil, err
	}

	return &pro.GetArtisanRankingsResponse{Rankings: rankings}, nil
}
