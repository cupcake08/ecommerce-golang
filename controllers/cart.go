package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/cupcake08/ecommerce-go/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewApplication(userCollection, productCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
	}
}
func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is empty"))
			return
		}
		userQueryID := c.Query("userId")
		if userQueryID == "" {
			log.Println("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Product ID is not valid")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is not valid"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.AddProductToCart(ctx, app.productCollection, app.userCollection, productID, userQueryID); err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "successfully added to cart")
	}
}

func (app *Application) RemoveFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("Product ID is empty"))
			return
		}
		userQueryID := ctx.Query("userId")
		if userQueryID == "" {
			log.Println("User ID is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Product ID is not valid")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("Product ID is not valid"))
			return
		}
		ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.RemoveItemFromCart(ct, app.productCollection, app.userCollection, productID, userQueryID); err != nil {
			log.Println(err)
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "successfully removed from cart")
	}
}

func (app *Application) InstantCheckout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("Product ID is empty"))
			return
		}
		userQueryID := ctx.Query("userId")
		if userQueryID == "" {
			log.Println("User ID is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println("Product ID is not valid")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("Product ID is not valid"))
			return
		}
		ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.InstantBuyer(ct, app.productCollection, app.userCollection, productID, userQueryID); err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "successfully placed the order")
	}
}

func GetItemFromCart() gin.HandlerFunc {

}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Println("User ID is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.BuyItemFromCart(ctx, app.productCollection, app.userCollection, userQueryID); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "successfully bought from cart")
	}
}
