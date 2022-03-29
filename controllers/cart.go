package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/cupcake08/ecommerce-go/database"
	"github.com/cupcake08/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("User ID is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is empty"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("User ID is not valid")
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filledCart := models.User{}
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		if err := UserCollection.FindOne(ctx, filter).Decode(&filledCart); err != nil {
			log.Println("err")
			c.IndentedJSON(http.StatusInternalServerError, "Not found")
			return
		}
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		cursor, er := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if er != nil {
			log.Println(er)
		}
		listing := make([]bson.M, 0)

		if err = cursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			log.Println(json)
			c.IndentedJSON(http.StatusOK, json["total"])
			c.IndentedJSON(http.StatusOK, filledCart.User_Cart)
		}
		ctx.Done()
	}
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
