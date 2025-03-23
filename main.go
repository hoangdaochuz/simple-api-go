package main

import (
	"fmt"

	"example.com/go-api/docs"
	"example.com/go-api/internal/config"
	"example.com/go-api/internal/database"
	"example.com/go-api/internal/middlewares"
	"example.com/go-api/internal/repository"
	"example.com/go-api/internal/transport"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {


	fmt.Println("Hello, World!")
	config := config.LoadConfig()
	fmt.Printf("DBHost: %s\n", config.DBHost)

	database.Connect(config)

	// Migrate the schema
	database.DB.AutoMigrate(&repository.Event{},&repository.User{}, &repository.Product{})


	server := gin.Default()
	docs.SwaggerInfo.Title = "Swagger Learn GO API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.BasePath = "/"
	

	// EVENT ROUTES
	server.GET("/events", transport.GetAllEventsHandler)
	server.GET("/events/:id",transport.GetEventById)
	server.POST("/events",transport.CreateEventHandler)
	server.PUT("events/:id",transport.UpdateEventHandler)
	server.DELETE("events/:id",transport.DeleteEventHandler)
	server.GET("/events/search", transport.SearchEventsHandler)

	// PRODUCTS ROUTES
	server.GET("/products", transport.GetAllProductsHandler)
	server.GET("/products/:id",transport.GetProductById)
	server.POST("/products",transport.CreateProductHandler)
	server.PUT("products/:id",transport.UpdateProductHandler)
	server.DELETE("products/:id",transport.DeleteProductHandler)
	server.GET("/products/recently", transport.GetRecentlyProductsHandler)

	// AUTH ROUTES
	server.POST("auth/register", transport.UserRegisterHandler)
	server.POST("auth/login", transport.UserLoginHandler)
	server.GET("auth/refresh", transport.UserRefreshTokenHandler)
	server.GET("auth/logout", transport.UserLogoutHandler)

	// USER ROUTES
	server.GET("users/my-profile",middlewares.CheckAuth ,transport.GetUserMyProfileHandler)


	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	

	server.Run(":" + config.ServerPort)
	fmt.Println("App is running on port ", config.ServerPort)
}