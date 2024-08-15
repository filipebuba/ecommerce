package ports

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItemRepository interface {
	AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error
	RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error
	BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error
	InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, userID string, productID primitive.ObjectID) error
}