package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"BikeStoreGolang/api-gateway/internal/client"
	"BikeStoreGolang/api-gateway/internal/handlers"
	"BikeStoreGolang/api-gateway/internal/logger"
	"BikeStoreGolang/api-gateway/internal/service"
)

func main() {
    logFile := "api-gateway.log"
    log, err := logger.NewLogrusLoggerToFile(logFile)
    godotenv.Load(".env")
	if err != nil {
        log.Warn(".env file not found or failed to load")
}
    // gRPC connections
    authConn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR"), grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Не удалось подключиться к auth-service: %v", err)
    }
    defer authConn.Close()

    productConn, err := grpc.Dial(os.Getenv("PRODUCT_SERVICE_ADDR"), grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Не удалось подключиться к product-service: %v", err)
    }
    defer productConn.Close()

    // Clients
    authClient := client.NewAuthClient(authConn)
    productClient := client.NewProductClient(productConn)

    // Services
    authService := service.NewAuthService(authClient.Client)
    productService := service.NewProductService(productClient.Client)

    // Handlers
    authHandler := handlers.NewAuthHandler(authService)
    productHandler := handlers.NewProductHandler(productService)

    // Routes
    http.HandleFunc("/login", authHandler.Login)
    http.HandleFunc("/register", authHandler.Register)
    http.HandleFunc("/activate", authHandler.Activate)
    http.HandleFunc("/forgot-password", authHandler.ForgotPassword)
    http.HandleFunc("/reset-password", authHandler.ResetPassword)
    http.HandleFunc("/refresh-token", authHandler.RefreshToken)
    http.HandleFunc("/me", authHandler.GetMe)
    http.HandleFunc("/logout", authHandler.Logout)

	  // Product routes
    http.HandleFunc("/products", productHandler.ListProducts) // GET
    http.HandleFunc("/create", productHandler.CreateProduct) // POST
    http.HandleFunc("/product/", productHandler.ProductCRUD)  // GET, PUT, DELETE by id


    log.Info("API Gateway запущен на :8080",)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Ошибка запуска API Gateway: %v", err)
    }
}