package routes

import (
	"project-e-commerce/controllers"
	"project-e-commerce/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProductByID)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/reviews", controllers.CreateReview)
	}

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/profile", controllers.Profile)
		admin.POST("/products", controllers.CreateProduct)
		admin.PUT("/products/:id", controllers.UpdateProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)
		admin.GET("/dashboard", controllers.GetDashboard)
		admin.GET("/orders", controllers.GetAllOrders)
		admin.PATCH("/orders/status/:id", controllers.UpdateOrderStatus)
	}
	cart := r.Group("/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.POST("/items", controllers.AddToCart)
		cart.PUT("/items", controllers.UpdateCartItem)
		cart.GET("/", controllers.GetCart)

	}
	orders := r.Group("/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("/checkout", controllers.Checkout)
		orders.GET("/", controllers.GetMyOrders)
		orders.GET("/:id", controllers.GetOrderDetail)
		orders.POST("/:id/pay", controllers.PayOrder)
		orders.PUT("/:id/cancel", controllers.CancelOrder)
	}
}
