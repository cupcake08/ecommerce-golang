package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cupcake08/ecommerce-go/database"
	"github.com/cupcake08/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Validator *validator.Validate = validator.New()
var UserCollection *mongo.Collection = database.UserData()

func HashPassword(password string) string {
	return password
}

func VerifyPassword(password string, hash string) (bool, string) {

}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user := &models.User{}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := Validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already exists",
			})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Phone already exists",
			})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()

		token, refresh_token := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *&user.ID)
		user.Token = &token
		user.Refresh_Token = &refresh_token
		user.User_Cart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, inserErr := UserCollection.InsertOne(ctx, user)

		if inserErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": inserErr.Error(),
			})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user := &models.User{}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			log.Fatal("email or password is incorrect")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var foundUser models.User
		isValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !isValid {
			log.Fatal("email or password is incorrect")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": msg,
			})
			return
		}

		token, refresh_token := generate.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_Id)
		defer cancel()

		generate.UpdateAllTokens(token, refresh_token, foundUser.User_Id)

		c.JSON(http.StatusFound, gin.H{
			"message": "User logged in successfully",
		})

	}
}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

}

func SearchProductByQuery() gin.HandlerFunc {

}
