package domain

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
    RoleCustomer Role = "customer"
    RoleAdmin    Role = "admin"
)

type User struct {
    ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name             string             `bson:"name" json:"name"`
    Email            string             `bson:"email" json:"email"`
    PasswordHash     string             `bson:"password_hash" json:"-"`
    Role             Role               `bson:"role" json:"role"`
    CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
    IsActive         bool               `bson:"is_active" json:"is_active"`
    ActivationToken  string             `bson:"activation_token,omitempty" json:"-"`
    ActivationExpires time.Time         `bson:"activation_expires,omitempty" json:"-"`
    ResetToken       string             `bson:"reset_token,omitempty" json:"-"`
    ResetExpires     time.Time          `bson:"reset_expires,omitempty" json:"-"`
}