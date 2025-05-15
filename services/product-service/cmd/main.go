package main

import (
    "context"
    "net"
    "os"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    deliverygrpc "BikeStoreGolang/services/product-service/internal/delivery/grpc"
    "BikeStoreGolang/services/product-service/internal/logger"
    "BikeStoreGolang/services/product-service/internal/usecase"
    pb "BikeStoreGolang/services/product-service/proto/gen"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    // Логгер
    logFile := "product-service.log"
    log, err := logger.NewLogrusLoggerToFile(logFile)
    if err != nil {
        panic("Failed to initialize logger: " + err.Error())
    }

    // .env
    if err := godotenv.Load(".env"); err != nil {
        log.Warn("Warning: .env file not found, using system environment variables")
    }

    mongoURI := os.Getenv("MONGO_URI")
    mongoDB := os.Getenv("MONGO_DB")
    if mongoURI == "" || mongoDB == "" {
        log.Fatal("MONGO_URI or MONGO_DB not set in environment")
    }

    // MongoDB
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("MongoDB connection error: ", err)
    }
    if err := client.Ping(context.Background(), nil); err != nil {
        log.Fatal("MongoDB ping error: ", err)
    }
    productsCollection := client.Database(mongoDB).Collection("products")

    // Usecase
    productUC := usecase.NewProductUsecase(productsCollection, log)

     authConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal("Failed to connect to AuthService: ", err)
    }
    defer authConn.Close()
    authClient := pb.NewAuthServiceClient(authConn)
    
    // gRPC server
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatal("Failed to listen: ", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterProductServiceServer(grpcServer,  deliverygrpc.NewProductHandler(productUC, authClient))

    log.Info("ProductService gRPC server started on :50052")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatal("Failed to serve: ", err)
    }
}