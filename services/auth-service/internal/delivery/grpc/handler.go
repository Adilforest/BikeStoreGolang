package grpc

import (
	"context"
	"time"

	"BikeStoreGolang/services/auth-service/internal/usecase"
	pb "BikeStoreGolang/services/auth-service/proto/github.com/adilforest/BikeStoreGolang/services/auth-service/proto/authpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user, err := h.uc.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := h.uc.SessionUC.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.AuthResponse{
		UserId: user.ID,
		Token:  token,
		Role:   string(user.Role),
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, token, err := h.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		UserId: user.ID,
		Token:  token,
		Role:   string(user.Role),
	}, nil
}

func (h *AuthHandler) GetUserByID(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.uc.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h *AuthHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.uc.DeleteUserByID(ctx, req.Id)
	if err != nil {
		return &pb.DeleteUserResponse{Success: false}, err
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}