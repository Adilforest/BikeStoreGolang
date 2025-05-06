package domain

import (
    "time"
)

type Role string

const (
    RoleCustomer Role = "customer"
    RoleAdmin    Role = "admin"
)

type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    PasswordHash string
    Role      Role      `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    IsActive  bool      `json:"is_active"`
}