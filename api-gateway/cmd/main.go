package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "BikeStoreGolang/api-gateway/proto/gen"
	

	"google.golang.org/grpc"
)

func main() {
	// Подключение к auth-service по gRPC
	authConn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к auth-service: %v", err)
	}
	defer authConn.Close()
	authClient := pb.NewAuthServiceClient(authConn)

	// Подключение к product-service по gRPC
	productConn, err := grpc.Dial(os.Getenv("PRODUCT_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к product-service: %v", err)
	}
	defer productConn.Close()
	productClient := pb.NewProductServiceClient(productConn)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var reqBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &reqBody); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := authClient.Login(ctx, &pb.LoginRequest{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		})
		if err != nil {
			http.Error(w, "gRPC error: "+err.Error(), http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Если фильтры не нужны, передайте пустой фильтр
		stream, err := productClient.ListProducts(ctx, &pb.ProductFilter{})
		if err != nil {
			http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var products []*pb.ProductResponse
		for {
			product, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "Stream error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})
}
