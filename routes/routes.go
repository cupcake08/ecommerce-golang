package routes

import (
	"github.com/cupcake08/ecommerce-go/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup", controllers.SignUp)
	incomingRoutes.POST("/user/login", controllers.Login)
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin)
	incomingRoutes.GET("/user/productview", controllers.SearchProduct)
	incomingRoutes.GET("/user/search", controllers.SearchProductByQuery)
}
