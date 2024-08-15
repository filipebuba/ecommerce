package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *string            `json:"price"`
	Rating       *string            `json:"rating"`
	Image        *string            `json:"image"`
}