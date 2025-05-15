package service

import (
	"BikeStoreGolang/api-gateway/proto/gen"
	"context"
	"io"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductService interface {
	ListProducts(ctx context.Context, filter *gen.ProductFilter) ([]*gen.ProductResponse, error)
	CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.ProductResponse, error)
	GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.ProductResponse, error)
	UpdateProduct(ctx context.Context, req *gen.UpdateProductRequest) (*gen.ProductResponse, error)
	DeleteProduct(ctx context.Context, req *gen.DeleteProductRequest) (*emptypb.Empty, error)
}

type productService struct {
	client gen.ProductServiceClient
}

func NewProductService(client gen.ProductServiceClient) ProductService {
	return &productService{client: client}
}

func (s *productService) ListProducts(ctx context.Context, filter *gen.ProductFilter) ([]*gen.ProductResponse, error) {
	stream, err := s.client.ListProducts(ctx, filter)
	if err != nil {
		return nil, err
	}
	var products []*gen.ProductResponse
	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *productService) CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.ProductResponse, error) {
	return s.client.CreateProduct(ctx, req)
}

func (s *productService) GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.ProductResponse, error) {
	return s.client.GetProduct(ctx, req)
}

func (s *productService) UpdateProduct(ctx context.Context, req *gen.UpdateProductRequest) (*gen.ProductResponse, error) {
	return s.client.UpdateProduct(ctx, req)
}

func (s *productService) DeleteProduct(ctx context.Context, req *gen.DeleteProductRequest) (*emptypb.Empty, error) {
	return s.client.DeleteProduct(ctx, req)
}
