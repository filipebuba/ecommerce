package service

import (
	"context"

	"github.com/filipebuba/ecommerce-yt/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type itemServiceImpl struct {
	repo ports.ItemRepository
}

func NewService(repo ports.ItemRepository) ports.ItemService{
	return &itemServiceImpl{
		repo: repo,
	}
}

func (s *itemServiceImpl) AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	return s.repo.AddProductToCart(ctx, productCollection, userCollection, productID, userID) 
}

func (s *itemServiceImpl) RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	return s.repo.RemoveCartItem(ctx, prodCollection, userCollection, productID, userID)
}

func (s *itemServiceImpl) BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error{
	return s.repo.BuyItemFromCart(ctx, userCollection, userID)
}

func (s *itemServiceImpl) InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, userID string, productID primitive.ObjectID) error {
	return s.repo.InstantBuyer(ctx, prodCollection, userCollection, userID, productID)
}