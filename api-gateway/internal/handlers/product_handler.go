package handlers

import (
	"BikeStoreGolang/api-gateway/internal/service"
	"BikeStoreGolang/api-gateway/proto/gen"
    "google.golang.org/grpc/metadata"
    //"BikeStoreGolang/api-gateway/internal/logger"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type ProductHandler struct {
    Service service.ProductService
}
type ProductService interface {
	ListProducts(ctx context.Context, filter *gen.ProductFilter) ([]*gen.ProductResponse, error)
	CreateProduct(ctx context.Context, req *gen.CreateProductRequest) (*gen.ProductResponse, error)
	GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.ProductResponse, error)
	UpdateProduct(ctx context.Context, req *gen.UpdateProductRequest) (*gen.ProductResponse, error)
	DeleteProduct(ctx context.Context, req *gen.DeleteProductRequest) (*grpc.CallOption, error)

}

func NewProductHandler(s service.ProductService) *ProductHandler {
    return &ProductHandler{Service: s}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    products, err := h.Service.ListProducts(ctx, &gen.ProductFilter{})
    if err != nil {
        http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

// POST /product
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req gen.CreateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    token := r.Header.Get("Authorization")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if token != "" {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
    }
    resp, err := h.Service.CreateProduct(ctx, &req)
    if err != nil {
        http.Error(w, "gRPC error: "+err.Error(), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

// GET, PUT, DELETE /product/{id}
func (h *ProductHandler) ProductCRUD(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/product/")
    if id == "" {
        http.Error(w, "Product ID required", http.StatusBadRequest)
        return
    }
     token := r.Header.Get("Authorization")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if token != "" {
        ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
    }

    switch r.Method {
    case http.MethodGet:
        req := &gen.GetProductRequest{Id: id}
        resp, err := h.Service.GetProduct(ctx, req)
        if err != nil {
            http.Error(w, "gRPC error: "+err.Error(), http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    case http.MethodPut:
        var req gen.UpdateProductRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        req.Id = id
        resp, err := h.Service.UpdateProduct(ctx, &req)
        if err != nil {
            http.Error(w, "gRPC error: "+err.Error(), http.StatusBadRequest)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    case http.MethodDelete:
        _, err := h.Service.DeleteProduct(ctx, &gen.DeleteProductRequest{Id: id})
        if err != nil {
            http.Error(w, "gRPC error: "+err.Error(), http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}