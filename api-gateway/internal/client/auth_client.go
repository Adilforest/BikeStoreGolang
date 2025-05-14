package client

import (
    "BikeStoreGolang/api-gateway/proto/gen"
    "google.golang.org/grpc"
)

type AuthClient struct {
    Client gen.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
    return &AuthClient{
        Client: gen.NewAuthServiceClient(conn),
    }
}