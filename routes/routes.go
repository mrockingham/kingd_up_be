package routes

import (
	"kingdup/api"
	"kingdup/handlers/auth"
	order "kingdup/handlers/order" // Use an alias

	myorder "kingdup/handlers/order"
	"kingdup/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	router := r

	// Enable CORS
	router.Use(cors.Default())

	// Test
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Sync products
	router.GET("/api/sync-products", api.SyncProductsHandler)

	// Get all products
	router.GET("/api/products", api.GetProductsHandler)

	// Get product by ID
	router.GET("/api/products/:id", api.GetProductHandler)

	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/", order.CreateOrderHandler(db))
		orderGroup.GET("/me", middleware.JWTMiddleware(), myorder.GetUserOrdersHandler(db))
		orderGroup.GET("/:id", middleware.JWTMiddleware(), order.GetOrderByIDHandler(db))

		// Add more as needed
	}

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", auth.RegisterHandler(db))
		authGroup.POST("/login", auth.LoginHandler(db))
	}
	userGroup := router.Group("/user")
	userGroup.Use(middleware.JWTMiddleware())
	{
		userGroup.GET("/me", auth.MeHandler(db))
	}

	return router
}
