package main

import (
	"log"
	"os"

	"github.com/cupcake08/ecommerce-go/controllers"
	"github.com/cupcake08/ecommerce-go/database"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	//without logger and recovery
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/addToCart", app.AddToCart())
	router.GET("/removeFromCart", app.RemoveFromCart())
	router.GET("/cartCheckout", app.CartCheckout())
	router.GET("/instantCheckout", app.InstantCheckout())

	log.Fatal(router.Run(":" + port))
}
