package services

import (
	"context"
	"fmt"
	"log"
	auth "product-service/genproto/authentication_service"
	pro "product-service/genproto/product_service"
	"product-service/storage/postgres"

	"github.com/sirupsen/logrus"
)

type ProductService interface {
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

type productService struct {
	authService auth.AuthenticationServiceClient
	productRepo postgres.ProductRepo
	pro.UnimplementedProductServiceServer
	log *logrus.Logger
}

func NewProductService(authService auth.AuthenticationServiceClient, productRepo postgres.ProductRepo, log *logrus.Logger) ProductService {
	return &productService{authService: authService, productRepo: productRepo, log: log}
}

func (p *productService) AddProduct(c context.Context, product *pro.AddProductRequest) (*pro.AddProductResponse, error) {
	// Validate user id
	res, err := p.authService.GetUserInfo(c, &auth.GetUserInfoRequest{Id: product.ArtisanId})
	if err != nil {
		p.log.Errorf("user is not found: %v", err)
		return nil, fmt.Errorf("user is not found: %v", err)
	}
	if res.User.UserType != "artisan" {
		p.log.Errorf("user is not an artistan: %v", err)
		return nil, fmt.Errorf("user is not an artistant")
	}
	return p.productRepo.AddProduct(c, product)
}

func (p *productService) EditProduct(c context.Context, product *pro.EditProductRequest) (*pro.EditProductResponse, error) {
	// Validate user id
	res, err := p.authService.GetUserInfo(c, &auth.GetUserInfoRequest{Id: product.ArtisanId})
	if err != nil {
		p.log.Errorf("user is not found: %v", err)
		return nil, fmt.Errorf("user is not found: %v", err)
	}

	log.Println(res)
	if res.User.UserType != "artisan" {
		p.log.Errorf("user is not an artistant")
		return nil, fmt.Errorf("user is not an artistant")
	}
	return p.productRepo.EditProduct(c, product)
}

func (p *productService) DeleteProduct(c context.Context, request *pro.DeleteProductRequest) (*pro.DeleteProductResponse, error) {
	// Validate user id
	res, err := p.productRepo.DeleteProduct(c, request)
	if err != nil {
		p.log.Errorf("failed to delete product: %v", err)
		return nil, err
	}
	res.Message = "product deleted successfully"
	return res, nil
}

func (p *productService) GetProduct(c context.Context, request *pro.GetProductRequest) (*pro.GetProductResponse, error) {
	res, err := p.productRepo.GetProduct(c, request)
	if err != nil {
		p.log.Errorf("failed to get product: %v", err)
		return nil, err
	}
	fmt.Println(res)

	return res, nil
}

func (p *productService) GetProducts(c context.Context, request *pro.GetProductsRequest) (*pro.GetProductsResponse, error) {
	res, err := p.productRepo.GetProducts(c, request)
	if err != nil {
		p.log.Errorf("failed to get product: %v", err)
		return nil, err
	}
	return res, nil
}

func (p *productService) AddProductCategory(c context.Context, request *pro.AddProductCategoryRequest) (*pro.AddProductCategoryResponse, error) {
	return p.productRepo.AddProductCategory(c, request)
}

func (p *productService) AddArtisanCategory(c context.Context, request *pro.AddArtisanCategoryRequest) (*pro.AddArtisanCategoryResponse, error) {
	return p.productRepo.AddArtisanCategory(c, request)
}

func (p *productService) SearchProducts(c context.Context, request *pro.SearchProductsRequest) (*pro.SearchProductsResponse, error) {
	if request.Limit == "" {
		request.Limit = "10"
	}
	if request.Page == "" {
		request.Page = "1"
	}
	res, err := p.productRepo.SearchProducts(c, request)
	if err != nil {
		p.log.Error()
	}
	return res, nil
}

func (p *productService) AddRating(c context.Context, rating *pro.AddRatingRequest) (*pro.AddRatingResponse, error) {
	_, err := p.authService.GetUserInfo(c, &auth.GetUserInfoRequest{Id: rating.UserId})
	if err != nil {
		p.log.Errorf("user is not found: %v", err)
		return nil, fmt.Errorf("user is not found: %v", err)
	}

	return p.productRepo.AddRating(c, rating)
}

func (p *productService) GetRatings(c context.Context, request *pro.GetRatingsRequest) (*pro.GetRatingsResponse, error) {
	return p.productRepo.GetRatings(c, request)
}

func (p *productService) PlaceOrder(c context.Context, order *pro.PlaceOrderRequest) (*pro.PlaceOrderResponse, error) {
	return p.productRepo.PlaceOrder(c, order)
}

func (p *productService) CancelOrder(c context.Context, request *pro.CancelOrderRequest) (*pro.CancelOrderResponse, error) {
	return p.productRepo.CancelOrder(c, request)
}

func (p *productService) UpdateOrderStatus(c context.Context, request *pro.UpdateOrderStatusRequest) (*pro.UpdateOrderStatusResponse, error) {
	return p.productRepo.UpdateOrderStatus(c, request)
}

func (p *productService) GetOrders(c context.Context, request *pro.GetOrdersRequest) (*pro.GetOrdersResponse, error) {
	return p.productRepo.GetOrders(c, request)
}

func (p *productService) GetOrder(c context.Context, request *pro.GetOrderRequest) (*pro.GetOrderResponse, error) {
	return p.productRepo.GetOrder(c, request)
}

func (p *productService) PayOrder(c context.Context, request *pro.PayOrderRequest) (*pro.PayOrderResponse, error) {
	return p.productRepo.PayOrder(c, request)
}

func (p *productService) CheckPaymentStatus(c context.Context, request *pro.CheckPaymentStatusRequest) (*pro.CheckPaymentStatusResponse, error) {
	return p.productRepo.CheckPaymentStatus(c, request)
}

func (p *productService) UpdateShippingInfo(c context.Context, request *pro.UpdateShippingInfoRequest) (*pro.UpdateShippingInfoResponse, error) {
	return p.productRepo.UpdateShippingInfo(c, request)
}

func (p *productService) GetStatistics(c context.Context, request *pro.GetStatisticsRequest) (*pro.GetStatisticsResponse, error) {
	return p.productRepo.GetStatistics(c, request)
}

func (p *productService) GetUserActivity(c context.Context, request *pro.GetUserActivityRequest) (*pro.GetUserActivityResponse, error) {
	return p.productRepo.GetUserActivity(c, request)
}

func (p *productService) GetRecommendations(c context.Context, request *pro.GetRecommendationsRequest) (*pro.GetRecommendationsResponse, error) {
	return p.productRepo.GetRecommendations(c, request)
}

func (p *productService) GetArtisanRankings(c context.Context, request *pro.GetArtisanRankingsRequest) (*pro.GetArtisanRankingsResponse, error) {
	return p.productRepo.GetArtisanRankings(c, request)
}
