package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	pro "product-service/genproto/product_service"
	"product-service/models"
	"strconv"
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
		INSERT INTO products (id, name, description, price, category_id, artisan_id, quantity)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING 
		id,
		name,
        description,
        price,
        category_id,
        artisan_id,
		quantity,
        created_at
	`

	row := r.db.QueryRowContext(c, query,
		newId,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		product.ArtisanId,
		product.Quantity,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}
	newProduct := pro.AddProductResponse{}
	err := row.Scan(&newProduct.Id, &newProduct.Name, &newProduct.Description, &newProduct.Price, &newProduct.CategoryId, &newProduct.ArtisanId, &newProduct.Quantity, &newProduct.CreatedAt)
	return &newProduct, err
}

func (r *productRepo) EditProduct(c context.Context, product *pro.EditProductRequest) (*pro.EditProductResponse, error) {
	query := `
		UPDATE products SET name = $1, description = $2, price = $3, category_id = $4, artisan_id = $5, quantity = $6
		WHERE id = $7 AND deleted_at IS NULL
	`

	res, err := r.db.ExecContext(c, query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryId,
		product.ArtisanId,
		product.Quantity,
		product.Id)

	log.Println(product, res)

	if err != nil {
		return nil, err
	}

	return &pro.EditProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryId:  product.CategoryId,
		ArtisanId:   product.ArtisanId,
	}, nil
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
	var pr float64
	query := `
		SELECT id, name, description, price, category_id, artisan_id, quantity, created_at, updated_at 
		FROM products 
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(c, query, request.Id).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&pr,
		&product.CategoryId,
		&product.ArtisanId,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Product not found
		}
		return nil, err
	}

	// productRes := &pro.Product{
	// 	Id:          product.Id,
	// 	Name:        product.Name,
	// 	Description: product.Description,
	// 	Price:       pr,
	// 	CategoryId:  product.CategoryId,
	// 	ArtisanId:   product.ArtisanId,
	// 	Quantity:    product.Quantity,
	// 	CreatedAt:   product.CreatedAt,
	// 	UpdatedAt:   product.UpdatedAt,
	// }

	return &pro.GetProductResponse{Product: &product}, nil
}

func (r *productRepo) SearchProducts(ctx context.Context, req *pro.SearchProductsRequest) (*pro.SearchProductsResponse, error) {
	var products []*pro.Product
	query := `
        SELECT p.id, p.name, p.price, p.category_id
        FROM products p
        JOIN product_categories c ON p.category_id = c.id
        WHERE c.name LIKE $1
        AND p.name LIKE $2
        AND p.price BETWEEN $3 AND $4
        AND p.deleted_at IS NULL
        AND c.deleted_at IS NULL
		LIMIT $5
		OFFSET $6
    `
	limit, _ := strconv.Atoi(req.Limit)
	page, _ := strconv.Atoi(req.Page)
	offset := (page - 1) * limit

	rows, err := r.db.QueryContext(ctx, query, "%"+req.Category+"%", "%"+req.Query+"%", req.MinPrice, req.MaxPrice, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product pro.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.CategoryId); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
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
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	newOrderId := uuid.NewString()

	query := `INSERT INTO orders(id, user_id, total_amount, status, shipping_address) VALUES($1, $2, $3, $4, $5)`
	type shipping struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		Country string `json:"country"`
		ZipCode string `json:"zip_code"`
	}

	shippingAddress := shipping{
		Street: order.ShippingAddress,
	}

	shippingAddressJSON, err := json.Marshal(shippingAddress)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error marshaling shipping address: %v", err)
	}
	res, err := tx.ExecContext(c, query, newOrderId, order.UserId, 0, order.Status, shippingAddressJSON)
	if err != nil {
		return nil, fmt.Errorf("error inserting order: %v", err)
	}
	if n, err := res.RowsAffected(); n == 0 || err != nil {
		return nil, fmt.Errorf("error inserting order: %v", err)
	}
	items := order.GetItems()
	if len(items) == 0 {
		return nil, fmt.Errorf("empty order")
	}
	query = `
		SELECT price, quantity FROM products WHERE id = $1 AND deleted_at IS NULL
	`

	var total_amount float64
	for i, item := range items {
		var quantity int
		var price float64
		row := tx.QueryRowContext(c, query, item.ProductId)
		if row.Err() == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}

		err := row.Scan(&price, &quantity)

		if err != nil {
			return nil, fmt.Errorf("error scanning product: %v", err)
		}
		if quantity < int(item.Quantity) {
			return nil, fmt.Errorf("not enough stock for item %s", item.ProductId)
		}
		items[i].Price = float32(price)
		newOrderItemsId := uuid.NewString()
		res, err := tx.ExecContext(c, "INSERT INTO order_items(id, order_id, product_id, quantity, price) VALUES($1, $2, $3, $4, $5)", newOrderItemsId, newOrderId, item.ProductId, item.Quantity, price)
		if err != nil {
			return nil, fmt.Errorf("error inserting order item: %v", err)
		}
		if n, err := res.RowsAffected(); n == 0 || err != nil {
			return nil, fmt.Errorf("error inserting order item: %v", err)
		}

		total_amount += price * float64(item.Quantity)
	}

	query = `UPDATE orders SET total_amount = $1 WHERE id = $2`
	res, err = tx.ExecContext(c, query, total_amount, newOrderId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting order: %v", err)
	}
	if n, err := res.RowsAffected(); n == 0 || err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting order: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}
	return &pro.PlaceOrderResponse{
		Id:              newOrderId,
		UserId:          order.UserId,
		TotalPrice:      float32(total_amount),
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		Items:           items,
		CreatedAt:       time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (r *productRepo) CancelOrder(c context.Context, request *pro.CancelOrderRequest) (*pro.CancelOrderResponse, error) {
	query := `
		UPDATE orders SET status = 'canceled', updated_at = now() WHERE id = $1 AND status = 'pending'
	`
	_, err := r.db.ExecContext(c, query, request.Id)
	if err != nil {
		return nil, err
	}

	return &pro.CancelOrderResponse{
		Id:        request.Id,
		Status:    "canceled",
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *productRepo) UpdateOrderStatus(c context.Context, request *pro.UpdateOrderStatusRequest) (*pro.UpdateOrderStatusResponse, error) {
	query := `
		UPDATE orders SET status = $1, updated_at = now() WHERE id = $2
	`
	_, err := r.db.ExecContext(c, query, request.Status, request.Id)
	if err != nil {
		return nil, err
	}

	return &pro.UpdateOrderStatusResponse{
		Id:        request.Id,
		Status:    request.Status,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// GetOrders retrieves a list of orders with pagination
func (r *productRepo) GetOrders(c context.Context, request *pro.GetOrdersRequest) (*pro.GetOrdersResponse, error) {
	var orders []*models.Order
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

	var protoOrders []*pro.Order
	for _, order := range orders {
		protoOrder := &pro.Order{
			Id:              order.Id,
			UserId:          order.UserId,
			TotalAmount:     order.TotalAmount,
			Status:          order.Status,
			ShippingAddress: order.ShippingAddress,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
		}
		protoOrders = append(protoOrders, protoOrder)
	}

	return &pro.GetOrdersResponse{Orders: protoOrders}, nil
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
	var amount float32
	row := r.db.QueryRowContext(c, "SELECT total_amount FROM orders WHERE id = $1 AND deleted_at IS NULL AND status = 'pending'", request.OrderId)
	if row == nil {
		return nil, fmt.Errorf("order not found")
	}
	err := row.Scan(&amount)
	if err != nil {
		return nil, err
	}
	transacId := uuid.NewString()
	newPaymentId := uuid.NewString()
	query := `
			INSERT INTO payments (
			id, 
			order_id, 
			amount, 
			status,
			transaction_id,
			payment_method)
			VALUES (
			$1, 
			$2, 
			$3,
			'paid', 
			$4,
            $5)
	`
	_, err = r.db.ExecContext(c, query,
		newPaymentId,
		request.OrderId,
		amount,
		transacId,
		request.PaymentMethod,
	)
	if err != nil {
		return nil, err
	}

	return &pro.PayOrderResponse{
		Id:            newPaymentId,
		OrderId:       request.OrderId,
		Amount:        float64(amount),
		PaymentMethod: request.PaymentMethod,
		Status:        "paid",
		TransactionId: transacId,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *productRepo) CheckPaymentStatus(c context.Context, request *pro.CheckPaymentStatusRequest) (*pro.CheckPaymentStatusResponse, error) {

	qury := `
		SELECT id, order_id, amount, status, transaction_id, payment_method, created_at
        FROM payments 
        WHERE id = $1 AND order_id = $2 AND deleted_at IS NULL
	`

	var payment pro.Payment
	row := r.db.QueryRowContext(c, qury, request.PaymentId, request.OrderId)
	if row == nil {
		return nil, fmt.Errorf("payment not found %v", row.Err())
	}
	err := row.Scan(&payment.Id, &payment.OrderId, &payment.Amount, &payment.Status, &payment.TransactionId, &payment.PaymentMethod, &payment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &pro.CheckPaymentStatusResponse{
		Payment: &payment,
	}, nil
}

func (r *productRepo) UpdateShippingInfo(c context.Context, request *pro.UpdateShippingInfoRequest) (*pro.UpdateShippingInfoResponse, error) {
	updatedAt := time.Now()
	query := `
		UPDATE orders SET 
		tracking_number = $1, 
		carrier = $2, 
		estimated_delivery_date = $3, 
		updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(c, query,
		request.TrackingNumber,
		request.Carrier,
		request.EstimatedDeliveryDate,
		updatedAt,
		request.OrderId,
	)
	if err != nil {
		return nil, err
	}

	return &pro.UpdateShippingInfoResponse{
		OrderId:               request.OrderId,
		TrackingNumber:        request.TrackingNumber,
		Carrier:               request.Carrier,
		EstimatedDeliveryDate: request.EstimatedDeliveryDate,
		UpdatedAt:             updatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *productRepo) AddArtisanCategory(c context.Context, request *pro.AddArtisanCategoryRequest) (*pro.AddArtisanCategoryResponse, error) {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	newId := uuid.NewString()
	query := `
		INSERT INTO artisan_categories (id, name, Description, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(c, query,
		newId,
		request.Name,
		request.Description,
		createdAt,
	)

	if err != nil {
		return nil, err
	}

	return &pro.AddArtisanCategoryResponse{
		Id:          newId,
		Name:        request.Name,
		Description: request.Description,
		CreatedAt:   createdAt,
	}, nil
}

func (r *productRepo) AddProductCategory(c context.Context, request *pro.AddProductCategoryRequest) (*pro.AddProductCategoryResponse, error) {
	newId := uuid.NewString()
	query := `
		INSERT INTO product_categories (id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, created_at
	`
	row := r.db.QueryRowContext(c, query, newId, request.Name, request.Description)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var productCat pro.AddProductCategoryResponse
	err := row.Scan(&productCat.Id, &productCat.Name, &productCat.Description, &productCat.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &productCat, nil
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
	// I ned to change this in protobuf
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
	//I need to change this response in protobuf
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
