package di

import (
	"BikeStoreGolang/services/auth-service/internal/delivery/grpc"
	"BikeStoreGolang/services/auth-service/internal/repository/postgres"
	"BikeStoreGolang/services/auth-service/internal/usecase"
	"database/sql"
)

func InitializeAuthService(db *sql.DB, jwtSecret string) *grpc.AuthHandler {
    userRepo := postgres.NewUserRepo(db)
    sessionUC := usecase.NewSessionUsecase(jwtSecret)
    authUC := usecase.NewAuthUsecase(userRepo, sessionUC)
    return grpc.NewAuthHandler(authUC)
}