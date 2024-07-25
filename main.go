package main

import (
	"log"
	"os"

	"github.com/filipebuba/ecommerce-yt/controllers"
	"github.com/filipebuba/ecommerce-yt/database"
	"github.com/filipebuba/ecommerce-yt/middleware"
	"github.com/filipebuba/ecommerce-yt/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuuFrom())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
