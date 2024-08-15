package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }

    fmt.Println("Successfully connected to MongoDB")
    return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
    if client == nil {
        log.Fatal("MongoDB client is nil")
    }
    return client.Database("Ecommerce").Collection(collectionName)
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
    if client == nil {
        log.Fatal("MongoDB client is nil")
    }
    return client.Database("Ecommerce").Collection(collectionName)
}