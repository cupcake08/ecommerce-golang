package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cupcake08/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("user_id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "invalid search query",
			})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			log.Println("invalid user id")
			c.IndentedJSON(http.StatusInternalServerError, "invalid user id")
		}

		addresses := models.Address{}

		addresses.Address_ID = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address_details"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		cursor, e := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, grouping})
		if e != nil {
			log.Println(e)
			c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
		}

		addressinfo := make([]bson.M, 0)
		if err = cursor.All(ctx, &addressinfo); err != nil {
			panic(err)
		}

		var size int32
		for _, addess_no := range addressinfo {
			size = addess_no["count"].(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address_details", Value: addresses}}}}
			_, e := UserCollection.UpdateOne(ctx, filter, update)
			if e != nil {
				log.Println(e)
				c.IndentedJSON(http.StatusInternalServerError, "Internal server error")
			}
		} else {
			c.IndentedJSON(400, "you can't add more than 2 address")
		}
		ctx.Done()
	}

}

func EditHomeAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("user_id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid search query",
			})
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("invalid user_id")
			c.IndentedJSON(http.StatusBadRequest, "invalid user id")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, er := UserCollection.UpdateOne(ctx, filter, update)
		if er != nil {
			c.IndentedJSON(http.StatusBadRequest, "wrong command")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, "address deleted")
	}
}
