package tokens

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/filipebuba/ecommerce-yt/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "UserS")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email, firt_name, last_name, uid string) (signedtoken string, segnedrefreshtoken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_Name: firt_name,
		Last_name:  last_name,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	claims.StandardClaims.ExpiresAt = time.Now().Local().Add(time.Hour * time.Duration(1)).Unix()
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err

}

func ValidateToken(SignedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(SignedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token has expired"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedtoken string, signedrefreshtoken string, uid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updated_at", Value: updated_at})

	upsert := true

	filter := bson.M{"use_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	UserData.UpdateOne(ctx, filter, bson.D{
		{"$set", updateobj}},
		&opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}

}
