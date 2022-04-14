package database

import (
	"context"
	"errors"
	"log"

	"github.com/cupcake08/ecommerce-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Error Variables

var (
	ErrorNotFound         = errors.New("product not found")
	ErrCantDecodeProduct  = errors.New("can't decode product")
	ErrUserIdNotValid     = errors.New("user Id is not valid")
	ErrCantUpdateUser     = errors.New("can't update user")
	ErrCantRemoveItemCart = errors.New("can't remove item from cart")
	ErrCantBuyCartItem    = errors.New("can't buy cart item")
	ErrCantGetItem        = errors.New("can't get item")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})

	if err != nil {
		log.Println(err)
		return ErrorNotFound
	}

	var productCart []models.ProductUser

	err = searchFromDB.All(ctx, &productCart)

	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	update := bson.D{{Key: "$push", Value: bson.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveItemFromCart(ctx context.Context, userCollection, prodCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
