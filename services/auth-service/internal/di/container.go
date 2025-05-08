package di

import (
	"BikeStoreGolang/services/auth-service/internal/delivery/grpc"
	"BikeStoreGolang/services/auth-service/internal/repository/postgres"
	"BikeStoreGolang/services/auth-service/internal/usecase"
	"database/sql"
)

func InitializeAuthService(db *sql.DB, jwtSecret string) *grpc.AuthHandler {
	// Создаем репозиторий пользователей
	userRepo := postgres.NewUserRepo(db)

	// Создаем SessionUsecase для работы с JWT
	sessionUC := usecase.NewSessionUsecase(jwtSecret)

	// Определяем "passwordPepper" (можно загрузить из конфигурации)
	passwordPepper := "your-password-pepper"

	// Создаем AdminActionLogger (заглушка или реальная реализация)
	adminLogger := usecase.NewAdminLogger()

	// Создаем AuthUsecase, передавая необходимые зависимости
	authUC := usecase.NewAuthUsecase(userRepo, sessionUC, adminLogger, passwordPepper)

	// Возвращаем обработчик gRPC
	return grpc.NewAuthHandler(authUC)
}
