package grpc

import (
	"context"

	"BikeStoreGolang/services/product-service/internal/usecase"
	pb "BikeStoreGolang/services/product-service/proto/gen"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	uc *usecase.ProductUsecase
}

func NewProductHandler(uc *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	// role, ok := ctx.Value("role").(string)
    // if !ok || role != "admin" {
    //     return nil, status.Error(codes.PermissionDenied, "only admin can create products")
    // }

	return h.uc.CreateProduct(ctx, req)
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
    // role, ok := ctx.Value("role").(string)
    // if !ok || role != "admin" {
    //     return nil, status.Error(codes.PermissionDenied, "only admin can update products")
    // }

    return h.uc.UpdateProduct(ctx, req)
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
    // role, ok := ctx.Value("role").(string)
    // if !ok || role != "admin" {
    //     return nil, status.Error(codes.PermissionDenied, "only admin can delete products")
    // }
    return h.uc.DeleteProduct(ctx, req)
}

func (h *ProductHandler) ChangeProductStock(ctx context.Context, req *pb.ChangeStockRequest) (*pb.ProductResponse, error) {
    // role, ok := ctx.Value("role").(string)
    // if !ok || role != "admin" {
    //     return nil, status.Error(codes.PermissionDenied, "only admin can change product stock")
    // }
    return h.uc.ChangeProductStock(ctx, req)
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	return h.uc.GetProduct(ctx, req)
}

func (h *ProductHandler) ListProducts(req *pb.ProductFilter, stream pb.ProductService_ListProductsServer) error {
	return h.uc.ListProducts(stream.Context(), req, stream)
}

func (h *ProductHandler) SearchProducts(req *pb.SearchRequest, stream pb.ProductService_SearchProductsServer) error {
    return h.uc.SearchProducts(req, stream)
}