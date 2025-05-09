package usecase

import (
	"BikeStoreGolang/services/auth-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	AdminLogger AdminActionLogger
	userRepo       domain.UserRepository
	SessionUC      *SessionUsecase
	adminLogger    AdminActionLogger
	passwordPepper string
}

type AdminActionLogger interface {
	LogAction(adminID, action, targetID, details string)
}

type ConsoleAdminLogger struct{}

func (c *ConsoleAdminLogger) LogAction(adminID, action, targetID, details string) {
	fmt.Printf("AdminID: %s performed Action: %s on TargetID: %s with Details: %s\n", adminID, action, targetID, details)
}

func NewAdminLogger() AdminActionLogger {
	return &ConsoleAdminLogger{}
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	sessionUC *SessionUsecase,
	adminLogger AdminActionLogger,
	passwordPepper string,
) *AuthUsecase {
	return &AuthUsecase{
		AdminLogger: adminLogger,
		userRepo:       userRepo,
		SessionUC:      sessionUC,
		adminLogger:    adminLogger,
		passwordPepper: passwordPepper,
	}
}

// validateEmail проверяет корректность email
func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// validateRegistrationInput проверяет входные данные для регистрации
func validateRegistrationInput(name, email, password string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if !validateEmail(email) {
		return errors.New("invalid email format")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

// Register регистрирует нового пользователя
func (u *AuthUsecase) Register(ctx context.Context, name, email, password, role string) (*domain.User, error) {
	if err := validateRegistrationInput(name, email, password); err != nil {
		return nil, err
	}

	// Проверяем, существует ли email
	if _, err := u.userRepo.GetByEmail(ctx, email); err == nil {
		return nil, domain.ErrEmailExists
	}

	normalizedRole := normalizeRole(role)
	return u.createUser(ctx, name, email, password, normalizedRole)
}

// normalizeRole нормализует роль пользователя
func normalizeRole(role string) domain.Role {
	switch domain.Role(role) {
	case domain.RoleAdmin:
		return domain.RoleAdmin
	default:
		return domain.RoleCustomer
	}
}

// Login выполняет аутентификацию пользователя
func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	fmt.Println("Login attempt with email:", email)

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		fmt.Println("Login failed: user not found for email:", email)
		return nil, "", errors.New("user not found")
	}

	if !user.IsActive {
		fmt.Println("Login failed: user is inactive for email:", email)
		return nil, "", errors.New("user is inactive")
	}

	pepperedPassword := password + u.passwordPepper
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pepperedPassword))
	if err != nil {
		fmt.Println("Login failed: invalid password for email:", email)
		return nil, "", errors.New("invalid password")
	}

	// Используем SessionUsecase для генерации токена
	token, err := u.SessionUC.GenerateToken(user.ID, user.Role)
	if err != nil {
		fmt.Println("Login failed: failed to generate token for user ID:", user.ID)
		return nil, "", errors.New("failed to generate token")
	}

	fmt.Println("Login successful for user ID:", user.ID)
	return user, token, nil
}

// createUser создает нового пользователя
func (u *AuthUsecase) createUser(ctx context.Context, name, email, password string, role domain.Role) (*domain.User, error) {
	pepperedPassword := password + u.passwordPepper
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pepperedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	return user, nil
}

func (u *AuthUsecase) AdminGetUser(ctx context.Context, adminID, userID string) (*domain.User, error) {
    // Проверяем, что запрашивающий - администратор
    admin, err := u.userRepo.GetByID(ctx, adminID)
    if err != nil || admin.Role != domain.RoleAdmin {
        return nil, domain.ErrAdminRequired
    }

    user, err := u.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }

    u.adminLogger.LogAction(adminID, "get user", userID, "")
    return user, nil
}

func (u *AuthUsecase) AdminGetAllUsers(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
    page, limit = normalizePagination(page, limit)

    users, err := u.userRepo.GetAllWithPagination(ctx, page, limit)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get users: %w", err)
    }

    fmt.Printf("Fetched users: %+v\n", users)

    total, err := u.userRepo.Count(ctx)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", err)
    }

    return users, total, nil
}

// Для AdminGetAllUsers с пагинацией
func (u *AuthUsecase) GetUsersWithPagination(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
    page, limit = normalizePagination(page, limit)
    
    users, err := u.userRepo.GetAllWithPagination(ctx, page, limit)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get users: %w", err)
    }
    
    total, err := u.userRepo.Count(ctx)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", err)
    }
    
    return users, total, nil
}

// normalizePagination нормализует параметры пагинации
func normalizePagination(page, limit int) (int, int) {
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20
    }
    return page, limit
}

// AdminUpdateUser обновляет данные пользователя (админ)
func (u *AuthUsecase) AdminUpdateUser(ctx context.Context, req *domain.AdminUpdateRequest) (*domain.User, error) {
    // 1. Проверка прав администратора
    admin, err := u.userRepo.GetByID(ctx, req.AdminID)
    if err != nil || admin.Role != domain.RoleAdmin {
        return nil, domain.ErrAdminRequired
    }

    // 2. Получаем пользователя для обновления
    user, err := u.userRepo.GetByID(ctx, req.UserID)
    if err != nil {
        return nil, domain.ErrUserNotFound
    }

    // 3. Проверка уникальности email если он изменяется
    if req.Email != nil && *req.Email != user.Email {
        if _, err := u.userRepo.GetByEmail(ctx, *req.Email); err == nil {
            return nil, domain.ErrEmailExists
        }
        user.Email = *req.Email
    }

    // 4. Обновляем остальные поля
    if req.Name != nil {
        user.Name = *req.Name
    }
    if req.Role != nil {
        user.Role = domain.Role(*req.Role)
    }
    if req.IsActive != nil {
        user.IsActive = *req.IsActive
    }

    user.UpdatedAt = time.Now()

    // 5. Сохраняем изменения
    if err := u.userRepo.Update(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    u.adminLogger.LogAction(req.AdminID, "update user", req.UserID, 
        fmt.Sprintf("Updated fields: Name:%t, Email:%t, Role:%t, Active:%t",
            req.Name != nil, req.Email != nil, req.Role != nil, req.IsActive != nil))
    
    return user, nil
}

// GetUserProfile возвращает профиль пользователя
func (u *AuthUsecase) GetUserProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Убираем хэш пароля
	user.PasswordHash = ""
	return user, nil
}

// Для AdminDeleteUser
func (u *AuthUsecase) AdminDeleteUser(ctx context.Context, adminID, userID string) error {
    if err := u.userRepo.DeleteByID(ctx, userID); err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    u.adminLogger.LogAction(adminID, "delete user", userID, "")
    return nil
}