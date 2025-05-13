package grpc

import (
	"context"

	"BikeStoreGolang/services/auth-service/internal/usecase"
	pb "BikeStoreGolang/services/auth-service/proto/gen"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return h.uc.Register(ctx, req)
}

func (h *AuthHandler) Activate(ctx context.Context, req *pb.ActivateRequest) (*pb.ActivateResponse, error) {
	return h.uc.Activate(ctx, req)
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return h.uc.Login(ctx, req)
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return h.uc.Logout(ctx, req)
}

func (h *AuthHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	return h.uc.ForgotPassword(ctx, req)
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	return h.uc.ResetPassword(ctx, req)
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return h.uc.RefreshToken(ctx, req)
}

func (h *AuthHandler) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.UserResponse, error) {
	userID := ""
	return h.uc.GetMe(ctx, userID)
}
