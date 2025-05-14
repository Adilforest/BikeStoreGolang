package client

import (
    "BikeStoreGolang/api-gateway/proto/gen"
    "google.golang.org/grpc"
)

type ProductClient struct {
    Client gen.ProductServiceClient
}

func NewProductClient(conn *grpc.ClientConn) *ProductClient {
    return &ProductClient{
        Client: gen.NewProductServiceClient(conn),
    }
}