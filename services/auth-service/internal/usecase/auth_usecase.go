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

// AdminGetAllUsers возвращает всех пользователей с пагинацией
func (u *AuthUsecase) AdminGetAllUsers(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
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

// AdminUpdateUser обновляет данные пользователя
func (u *AuthUsecase) AdminUpdateUser(ctx context.Context, adminID string, user *domain.User) (*domain.User, error) {
	existingUser, err := u.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Сохраняем неизменяемые поля
	user.PasswordHash = existingUser.PasswordHash
	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	u.adminLogger.LogAction(adminID, "update user", user.ID, "")
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