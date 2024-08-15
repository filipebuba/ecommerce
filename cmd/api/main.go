package main

import (
	"log"
	"os"

	"github.com/filipebuba/ecommerce-yt/cmd/routes"
	controllers "github.com/filipebuba/ecommerce-yt/internal/adapters/handler"
	database "github.com/filipebuba/ecommerce-yt/internal/adapters/repository/mongo"
	"github.com/filipebuba/ecommerce-yt/internal/core/service"
	"github.com/filipebuba/ecommerce-yt/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Obtém a porta do ambiente ou define como 8080 se não estiver definida
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	conn := database.DBSet()

	prodCollection := database.ProductData(database.Client, "products")
	userCollection := database.UserData(database.Client, "users")

	repo := database.NewMongoRepository(conn)
	service := service.NewService(repo)

	// Cria uma nova aplicação com as coleções de produtos e usuários

	app := controllers.NewApplication(service)
	// Cria um novo roteador Gin
	router := gin.New()
	router.Use(gin.Logger())

	// Define as rotas de usuário
	routes.UserRoutes(router)

	// Adiciona middleware de autenticação
	router.Use(middleware.Authentication())

	// Define as rotas da aplicação
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	// Inicia o servidor na porta especificada
	log.Fatal(router.Run(":" + port))
}
