package mongo

import (
    "context"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/joho/godotenv"
)

func LoadEnv(path string) error {
    return godotenv.Load(path)
}

func NewMongoClient() (*mongo.Client, string, error) {
    err := LoadEnv(".env")
    if err != nil {
        return nil, "", err
    }

    uri := os.Getenv("MONGO_URI")
    db := os.Getenv("MONGO_DB")
    if uri == "" || db == "" {
        return nil, "", ErrNoEnv
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, "", err
    }

    return client, db, nil
}

var ErrNoEnv = mongo.CommandError{Message: "MONGO_URI or MONGO_DB not set in environment"}