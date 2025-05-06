package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

type TokenRepository interface {
    SaveToken(ctx context.Context, userID string, token string) error
    GetToken(ctx context.Context, userID string) (string, error)
    DeleteToken(ctx context.Context, userID string) error
    ValidateToken(ctx context.Context, token string) (bool, error)
}