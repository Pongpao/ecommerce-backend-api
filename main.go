package main

// @title E-Commerce API
// @version 1.0
// @description Backend API for E-Commerce system
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
import (
	"github.com/gin-gonic/gin"

	"project-e-commerce/config"
	"project-e-commerce/middleware"
	"project-e-commerce/models"
	"project-e-commerce/routes"
	_ "project-e-commerce/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()

	config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderStatusHistory{},
		&models.Payment{},
		&models.Review{},
	)

	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	routes.SetupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := config.GetEnv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
