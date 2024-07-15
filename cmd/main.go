package main

import (
	"net"
	"product-service/configs"
	"product-service/genproto/authentication_service"
	"product-service/genproto/product_service"
	"product-service/logger"
	"product-service/services"
	"product-service/storage/postgres"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logger.InitLogger()

	log := logger.GetLogger()
	log.WithFields(logrus.Fields{
		"TestLogger": "test-logger",
	})

	config, err := configs.InitConfig(".")
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// Connect to PostgreSQL
	db, err := postgres.ConnectDB(config)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	log.Infof("Successfully connected to database")

	authConn, err := grpc.NewClient(":50051")
	if err != nil {
		log.Fatalf("Error connecting to authentication service: %v", err)
	}
	// serving product service
	listener, err := net.Listen("tcp", config.Host+":"+config.GrpcServerPort)
	if err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}

	log.Infof("gRPC server started on %s", config.Host+":"+config.GrpcServerPort)

	grpcServer := grpc.NewServer()

	authClient := authentication_service.NewAuthenticationServiceClient(authConn)

	productRepo := postgres.NewProductRepo(db)

	log.Infof("Initializing product service")

	productService := services.NewProductService(authClient, productRepo)

	product_service.RegisterProductServiceServer(grpcServer, productService.(product_service.ProductServiceServer))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
