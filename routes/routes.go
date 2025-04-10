package routes

import (
	"kingdup/api"
	"kingdup/handlers/auth"
	order "kingdup/handlers/order" // Use an alias
	"kingdup/handlers/payment"

	"kingdup/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {
	router := r

	// Unwrap *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		panic("❌ Failed to get raw DB from GORM: " + err.Error())
	}

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://kingdup.com", "https://www.kingdup.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		orderGroup.GET("/me", middleware.JWTMiddleware(), order.GetUserOrdersHandler(db))
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

	paymentGroup := router.Group("/payment")
	{
		paymentGroup.POST("/checkout", payment.CreateCheckoutHandler(sqlDB))
	}

	webhookGroup := router.Group("/webhook")
	{
		webhookGroup.POST("/stripe", payment.StripeWebhookHandler(sqlDB))
	}

	return router
}
