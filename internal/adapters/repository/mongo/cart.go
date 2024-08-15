package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/filipebuba/ecommerce-yt/internal/core/domain"
	"github.com/filipebuba/ecommerce-yt/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCartFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("cant find the product")
	ErrUserIdsNotValid    = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

type itemRepositoryImpl struct{
	conn *mongo.Client
}

func NewMongoRepository(conn *mongo.Client) ports.ItemRepository{
	return &itemRepositoryImpl{
		conn: conn,
	}
}

func (s *itemRepositoryImpl) AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := productCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCartFindProduct
	}
	var productCart []domain.ProductUser
	if err = searchfromdb.All(ctx, &productCart); err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser

	}
	return nil
}

func (s *itemRepositoryImpl) RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}
	return nil

}

func (s *itemRepositoryImpl) BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	// fetch the cart of user
	// find the cart total
	// create an order with the items
	// added order to the user Collection
	// addded items in the cart to the order list
	// empty up the cart

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdsNotValid
	}
	var getCartItems domain.User
	var orderCart domain.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]domain.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "Path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	currentresults, err := userCollection.Find(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M

	if err = currentresults.All(ctx, &getusercart); err != nil {
		log.Println(err)
		return ErrCantGetItem
	}
	var total_price int32

	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int32)
	}
	orderCart.Price = int(total_price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id"}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[]order_list": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateMany(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	usercart_empty := make([]domain.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
	_, err = userCollection.UpdateMany(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func (s *itemRepositoryImpl) InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, userID string, productID primitive.ObjectID) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdsNotValid
	}

	var product_details domain.ProductUser
	var order_details domain.Order

	order_details.Order_ID = primitive.NewObjectID()
	order_details.Ordered_At = time.Now()
	order_details.Order_Cart = make([]domain.ProductUser, 0)
	order_details.Payment_Method.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}

	order_details.Price = product_details.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order_details}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	return nil

}
