package postgres

import (
	"context"
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
		insert into product(id, name, description, price, category_id, Artisan_id, quanity, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(c, query,
		&newId,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CategoryId,
		&product.ArtisanId,
		&product.Quantity,
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

	return &pro.EditProductResponse{}, err
}

func (r *productRepo) DeleteProduct(c context.Context, request *pro.DeleteProductRequest) (*pro.DeleteProductResponse, error) {
	query := `
		UPDATE products SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(c, query,
		request.Id)
	return &pro.DeleteProductResponse{}, err
}

func (r *productRepo) GetProducts(c context.Context, request *pro.GetProductsRequest) (*pro.GetProductsResponse, error) {
	return &pro.GetProductsResponse{}, nil
}

func (r *productRepo) GetProduct(c context.Context, request *pro.GetProductRequest) (*pro.GetProductResponse, error) {
	var product *pro.Product
	err := r.db.SelectContext(c, product, "SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at FROM products WHERE id = $1 AND deleted_at IS NULL")
	return &pro.GetProductResponse{Product: product}, err
}

func (r *productRepo) SearchProducts(c context.Context, request *pro.SearchProductsRequest) (*pro.SearchProductsResponse, error) {
	return &pro.SearchProductsResponse{}, nil
}

func (r *productRepo) AddRating(c context.Context, rating *pro.AddRatingRequest) (*pro.AddRatingResponse, error) {
	return &pro.AddRatingResponse{}, nil
}

func (r *productRepo) GetRatings(c context.Context, request *pro.GetRatingsRequest) (*pro.GetRatingsResponse, error) {
	return &pro.GetRatingsResponse{}, nil
}

func (r *productRepo) PlaceOrder(c context.Context, order *pro.PlaceOrderRequest) (*pro.PlaceOrderResponse, error) {
	return &pro.PlaceOrderResponse{}, nil
}

func (r *productRepo) CancelOrder(c context.Context, request *pro.CancelOrderRequest) (*pro.CancelOrderResponse, error) {
	return &pro.CancelOrderResponse{}, nil
}

func (r *productRepo) UpdateOrderStatus(c context.Context, request *pro.UpdateOrderStatusRequest) (*pro.UpdateOrderStatusResponse, error) {
	return &pro.UpdateOrderStatusResponse{}, nil
}

func (r *productRepo) GetOrders(c context.Context, request *pro.GetOrdersRequest) (*pro.GetOrdersResponse, error) {
	return &pro.GetOrdersResponse{}, nil
}

func (r *productRepo) GetOrder(c context.Context, request *pro.GetOrderRequest) (*pro.GetOrderResponse, error) {
	return &pro.GetOrderResponse{}, nil
}

func (r *productRepo) PayOrder(c context.Context, request *pro.PayOrderRequest) (*pro.PayOrderResponse, error) {
	return &pro.PayOrderResponse{}, nil
}

func (r *productRepo) CheckPaymentStatus(c context.Context, request *pro.CheckPaymentStatusRequest) (*pro.CheckPaymentStatusResponse, error) {
	return &pro.CheckPaymentStatusResponse{}, nil
}

func (r *productRepo) UpdateShippingInfo(c context.Context, request *pro.UpdateShippingInfoRequest) (*pro.UpdateShippingInfoResponse, error) {
	return &pro.UpdateShippingInfoResponse{}, nil
}

func (r *productRepo) AddArtisanCategory(c context.Context, request *pro.AddArtisanCategoryRequest) (*pro.AddArtisanCategoryResponse, error) {
	return &pro.AddArtisanCategoryResponse{}, nil
}

func (r *productRepo) AddProductCategory(c context.Context, request *pro.AddProductCategoryRequest) (*pro.AddProductCategoryResponse, error) {
	return &pro.AddProductCategoryResponse{}, nil
}

func (r *productRepo) GetStatistics(c context.Context, request *pro.GetStatisticsRequest) (*pro.GetStatisticsResponse, error) {
	return &pro.GetStatisticsResponse{}, nil
}

func (r *productRepo) GetUserActivity(c context.Context, request *pro.GetUserActivityRequest) (*pro.GetUserActivityResponse, error) {
	return &pro.GetUserActivityResponse{}, nil
}

func (r *productRepo) GetRecommendations(c context.Context, request *pro.GetRecommendationsRequest) (*pro.GetRecommendationsResponse, error) {
	return &pro.GetRecommendationsResponse{}, nil
}

func (r *productRepo) GetArtisanRankings(c context.Context, request *pro.GetArtisanRankingsRequest) (*pro.GetArtisanRankingsResponse, error) {
	return &pro.GetArtisanRankingsResponse{}, nil
}
