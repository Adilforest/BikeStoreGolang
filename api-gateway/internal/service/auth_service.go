package service

import (
    "context"
    "BikeStoreGolang/api-gateway/proto/gen"
)

type AuthService interface {
    Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error)
	Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error)
	Activate(ctx context.Context, req *gen.ActivateRequest) (*gen.ActivateResponse, error)
	ForgotPassword(ctx context.Context, req *gen.ForgotPasswordRequest) (*gen.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req *gen.ResetPasswordRequest) (*gen.ResetPasswordResponse, error)
	RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error)
	GetMe(ctx context.Context, req *gen.GetMeRequest) (*gen.UserResponse, error)
	Logout(ctx context.Context, req *gen.LogoutRequest) (*gen.LogoutResponse, error)


}

type authService struct {
    client gen.AuthServiceClient
}

func NewAuthService(client gen.AuthServiceClient) AuthService {
    return &authService{client: client}
}

func (s *authService) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
    return s.client.Login(ctx, req)
}

func (s *authService) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
    return s.client.Register(ctx, req)
}
func (s *authService) Activate(ctx context.Context, req *gen.ActivateRequest) (*gen.ActivateResponse, error) {
    return s.client.Activate(ctx, req)
}
func (s *authService) ForgotPassword(ctx context.Context, req *gen.	ForgotPasswordRequest) (*gen.ForgotPasswordResponse, error) {
    return s.client.ForgotPassword(ctx, req)
}
func (s *authService) ResetPassword(ctx context.Context, req *gen.ResetPasswordRequest) (*gen.ResetPasswordResponse, error) {
	return s.client.ResetPassword(ctx, req)
}
func (s *authService) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error) {
    return s.client.RefreshToken(ctx, req)
}
func (s *authService) GetMe(ctx context.Context, req *gen.GetMeRequest) (*gen.UserResponse, error) {
    return s.client.GetMe(ctx, req)
}
func (s *authService) Logout(ctx context.Context, req *gen.LogoutRequest) (*gen.LogoutResponse, error) {
    return s.client.Logout(ctx, req)
}