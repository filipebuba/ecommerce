package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/filipebuba/ecommerce-yt/database"
	"github.com/filipebuba/ecommerce-yt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var prodCollection *mongo.Collection = database.ProductData(database.Client, "Products")

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.GenerateFromPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func Signup() gin.HandlerFunc {

	validate = validator.New(validator.WithRequiredStructEnabled())

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone no. is already in use"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Update_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, user.First_Name, user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Datails = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := userCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "tre user did not get created"})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		PasswordValid, msg := VerifyPassword(*user.Password, *founduser.Password)

		defer cancel()

		if !PasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}

}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	}

}

func SearchProductByQuery() gin.HandlerFunc {

}
