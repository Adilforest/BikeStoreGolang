package service

import (
	"BikeStoreGolang/api-gateway/proto/gen"
	"context"
	"io"

	"google.golang.org/grpc"
)

type ProductService interface {
	ListProducts(ctx context.Context, filter *gen.ProductFilter) ([]*gen.ProductResponse, error)
	CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.ProductResponse, error)
    GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.ProductResponse, error)
    UpdateProduct(ctx context.Context, req *gen.UpdateProductRequest) (*gen.ProductResponse, error)
    DeleteProduct(ctx context.Context, req *gen.DeleteProductRequest) (*grpc.CallOption, error)
	
}

func (p *productService) CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.ProductResponse, error) {
    panic("unimplemented")
}

func (p *productService) GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.ProductResponse, error) {
    panic("unimplemented")
}

func (p *productService) UpdateProduct(ctx context.Context, req *gen.UpdateProductRequest) (*gen.ProductResponse, error) {
    panic("unimplemented")
}

func (p *productService) DeleteProduct(ctx context.Context, req *gen.DeleteProductRequest) (*grpc.CallOption, error) {
    panic("unimplemented")
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
