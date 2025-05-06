package usecase

import (
	"BikeStoreGolang/services/auth-service/internal/domain"
	"context"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// AuthUsecase реализует бизнес-логику аутентификации
type AuthUsecase struct {
	userRepo domain.UserRepository
}

// NewAuthUsecase создает новый экземпляр AuthUsecase
func NewAuthUsecase(userRepo domain.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
	}
}

// validateEmail проверяет корректность email
func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Register регистрирует нового пользователя
func (u *AuthUsecase) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	if !validateEmail(email) {
		return nil, errors.New("invalid email format")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         domain.RoleCustomer,
		IsActive:     true,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID возвращает пользователя по ID
func (u *AuthUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

// GetAllUsers возвращает всех пользователей
func (u *AuthUsecase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return u.userRepo.GetAll(ctx)
}

// UpdateUser обновляет данные пользователя
func (u *AuthUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return u.userRepo.Update(ctx, user)
}

// DeleteUserByID удаляет пользователя по ID
func (u *AuthUsecase) DeleteUserByID(ctx context.Context, id string) error {
	return u.userRepo.DeleteByID(ctx, id)
}

// DeleteAllUsers удаляет всех пользователей
func (u *AuthUsecase) DeleteAllUsers(ctx context.Context) error {
	return u.userRepo.DeleteAll(ctx)
}
