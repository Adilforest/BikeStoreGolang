package grpc

import (
	"context"
	"fmt"
	"time"

	"BikeStoreGolang/services/auth-service/internal/domain"
	"BikeStoreGolang/services/auth-service/internal/usecase"
	pb "BikeStoreGolang/services/auth-service/proto/github.com/adilforest/BikeStoreGolang/services/auth-service/proto/authpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
	pb.UnimplementedAuthServiceServer
}

func convertToUserResponse(user *domain.User) *pb.UserResponse {
    return &pb.UserResponse{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      string(user.Role),
        IsActive:  user.IsActive,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
        UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
    }
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
    fmt.Printf("gRPC Register request: %+v\n", req)
    
    user, err := h.uc.Register(ctx, req.Name, req.Email, req.Password, req.Role)
    if err != nil {
        return nil, err
    }

    fmt.Printf("User after registration: %+v\n", user)
    
    token, err := h.uc.SessionUC.GenerateToken(user.ID, user.Role)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to generate token")
    }

    fmt.Printf("Generated token for role: %s\n", user.Role)
    
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

func (h *AuthHandler) GetUserByID(ctx context.Context, req *pb.AdminGetUserRequest) (*pb.UserResponse, error) {
    user, err := h.uc.AdminGetUser(ctx, req.AdminId, req.UserId)
    if err != nil {
        return nil, status.Error(convertErrorToCode(err), err.Error())
    }
    return convertToUserResponse(user), nil
}

func (h *AuthHandler) GetAllUsers(ctx context.Context, req *pb.AdminGetAllUsersRequest) (*pb.UsersListResponse, error) {
    fmt.Printf("Checking admin access for admin_id: %s\n", req.AdminId)

    // Проверяем, что admin_id существует и имеет роль администратора
    admin, err := h.uc.AdminGetUser(ctx, req.AdminId, req.AdminId) // Проверяем только admin_id
    if err != nil {
        fmt.Printf("Error fetching admin: %v\n", err)
        return nil, status.Error(codes.PermissionDenied, "admin access required")
    }

    if admin.Role != domain.RoleAdmin {
        fmt.Printf("User is not an admin: %s\n", admin.Role)
        return nil, status.Error(codes.PermissionDenied, "admin access required")
    }

    // Получаем пользователей с пагинацией
    users, total, err := h.uc.AdminGetAllUsers(ctx, int(req.Page), int(req.Limit))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
    }

    // Формируем ответ
    response := &pb.UsersListResponse{
        TotalCount: int32(total),
    }

    for _, user := range users {
        response.Users = append(response.Users, convertToUserResponse(user))
    }

    // Логируем действие администратора
    h.uc.AdminLogger.LogAction(req.AdminId, "list users", "", 
        fmt.Sprintf("Page: %d, Limit: %d", req.Page, req.Limit))
    return response, nil
}

func (h *AuthHandler) UpdateUser(ctx context.Context, req *pb.AdminUpdateUserRequest) (*pb.UserResponse, error) {
    // Проверяем, что admin_id и user_id указаны
    if req.AdminId == "" || req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "admin_id and user_id are required")
    }

    // Создаем запрос для обновления
    updateReq := &domain.AdminUpdateRequest{
        AdminID:  req.AdminId,
        UserID:   req.Id,
    }

    // Обрабатываем optional поля protobuf
    if req.Name != nil {
        updateReq.Name = req.Name
    }
    if req.Email != nil {
        updateReq.Email = req.Email
    }
    if req.Role != nil {
        updateReq.Role = req.Role
    }
    if req.IsActive != nil {
        updateReq.IsActive = req.IsActive
    }

    // Вызываем бизнес-логику для обновления пользователя
    user, err := h.uc.AdminUpdateUser(ctx, updateReq)
    if err != nil {
        return nil, status.Error(convertErrorToCode(err), err.Error())
    }

    // Возвращаем обновленного пользователя
    return convertToUserResponse(user), nil
}

func (h *AuthHandler) AdminDeleteUser(ctx context.Context, req *pb.AdminDeleteRequest) (*pb.DeleteUserResponse, error) {
    if err := h.uc.AdminDeleteUser(ctx, req.AdminId, req.UserId); err != nil {
        return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
    }
    return &pb.DeleteUserResponse{Success: true}, nil
}

func convertErrorToCode(err error) codes.Code {
    switch err {
    case domain.ErrUserNotFound:
        return codes.NotFound
    case domain.ErrAdminRequired:
        return codes.PermissionDenied
    case domain.ErrEmailExists:
        return codes.AlreadyExists
    default:
        return codes.Internal
    }
}
