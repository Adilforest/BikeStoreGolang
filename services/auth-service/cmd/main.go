package main

import (
	"BikeStoreGolang/services/auth-service/internal/di"
	"BikeStoreGolang/services/auth-service/internal/repository/postgres"
	"log"
	"net"

	pb "BikeStoreGolang/services/auth-service/proto/github.com/adilforest/BikeStoreGolang/services/auth-service/proto/authpb"

	"google.golang.org/grpc"
)

func main() {
	// 1. Подключение к БД
	db, err := postgres.NewDB("postgres://user:db@db:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Define the JWT secret
	jwtSecret := "your-secret-key"

	// 2. Инициализация зависимостей
	authHandler := di.InitializeAuthService(db, jwtSecret)

	// 3. Запуск gRPC-сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	srv := grpc.NewServer()
	pb.RegisterAuthServiceServer(srv, authHandler)

	log.Println("Auth service running on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
