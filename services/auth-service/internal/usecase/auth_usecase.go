package usecase

import (
	"BikeStoreGolang/services/auth-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo    domain.UserRepository
	SessionUC *SessionUsecase  
}

func NewAuthUsecase(userRepo domain.UserRepository, sessionUC *SessionUsecase) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  userRepo,
		SessionUC: sessionUC,
	}
}

// validateEmail проверяет корректность email
func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (u *AuthUsecase) Register(ctx context.Context, name, email, password, role string) (*domain.User, error) {
    fmt.Printf("Register request - name: %s, email: %s, role: %s\n", name, email, role)
    
    if !validateEmail(email) {
        return nil, errors.New("invalid email format")
    }
    if len(password) < 6 {
        return nil, errors.New("password must be at least 6 characters long")
    }

    // Проверка и нормализация роли
    var normalizedRole domain.Role
    switch role {
    case string(domain.RoleAdmin):
        normalizedRole = domain.RoleAdmin
        fmt.Println("Setting admin role")
    case string(domain.RoleCustomer):
        normalizedRole = domain.RoleCustomer
        fmt.Println("Setting customer role")
    default:
        normalizedRole = domain.RoleCustomer
        fmt.Println("Setting default customer role")
    }

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &domain.User{
        ID:           uuid.New().String(),
        Name:         name,
        Email:        email,
        PasswordHash: string(passwordHash),
        Role:         normalizedRole,
        IsActive:     true,
    }

    fmt.Printf("User before save: %+v\n", user)
    
    err = u.userRepo.Create(ctx, user)
    if err != nil {
        return nil, err
    }

    // Проверяем сохраненного пользователя
    savedUser, err := u.userRepo.GetByID(ctx, user.ID)
    if err != nil {
        fmt.Printf("Error fetching saved user: %v\n", err)
    } else {
        fmt.Printf("Saved user from DB: %+v\n", savedUser)
    }

    return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	if !user.IsActive {
		return nil, "", errors.New("user is inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid password")
	}

	// Используем SessionUsecase для генерации токена
	token, err := u.SessionUC.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
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
