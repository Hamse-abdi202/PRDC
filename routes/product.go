package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ilwaad12/product-project/handlers"
)

func RegisterRoutes(r *gin.Engine) {

	apiGroup := r.Group("/api") // optional
	ProductHandler := handlers.RegisterProductHandler()

	productGroup := apiGroup.Group("/products")
	{
		productGroup.POST("/create", ProductHandler.Create)
		productGroup.PUT("/:id", ProductHandler.Update)
		productGroup.DELETE("/:id", ProductHandler.SoftDelete)
		productGroup.POST("/:id/restore", ProductHandler.Restore)

		productGroup.GET("/list", ProductHandler.List)
		productGroup.GET("/:id", ProductHandler.GetByID)
	}
}
