package domain

import (
	"context"
	"errors"
)

type UserUpdate struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	GetAllWithPagination(ctx context.Context, page, limit int) ([]*User, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

type TokenRepository interface {
	SaveToken(ctx context.Context, userID, token string) error
	GetToken(ctx context.Context, userID string) (string, error)
	DeleteToken(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, token string) (bool, error)
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserInactive     = errors.New("user is inactive")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmailExists      = errors.New("email already exists")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidRole      = errors.New("invalid role")
	ErrSelfDeletion     = errors.New("self-deletion not allowed")
	ErrAdminDeletion    = errors.New("admin deletion requires special privileges")
)