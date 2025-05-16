package grpc

import (
	"context"
	"strings"

	"BikeStoreGolang/services/auth-service/internal/usecase"
	pb "BikeStoreGolang/services/auth-service/proto/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
    // Извлекаем токен из metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "no metadata provided")
    }
    authHeaders := md["authorization"]
    if len(authHeaders) == 0 {
        return nil, status.Error(codes.Unauthenticated, "no authorization header")
    }
    token := strings.TrimPrefix(authHeaders[0], "Bearer ")
    userID, err := h.uc.ParseUserIDFromToken(token)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }
    return h.uc.GetMe(ctx, userID)
}
