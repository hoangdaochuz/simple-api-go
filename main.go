package main

import (
	"fmt"

	"example.com/go-api/docs"
	"example.com/go-api/internal/config"
	"example.com/go-api/internal/database"
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
	database.DB.AutoMigrate(&repository.Event{})

	server := gin.Default()
	docs.SwaggerInfo.Title = "Swagger Learn GO API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.BasePath = "/"
	

	server.GET("/events", transport.GetAllEventsHandler)
	server.GET("/events/:id",transport.GetEventById)
	server.POST("/events",transport.CreateEventHandler)
	server.PUT("events/:id",transport.UpdateEventHandler)
	server.DELETE("events/:id",transport.DeleteEventHandler)
	server.GET("/events/search", transport.SearchEventsHandler)


	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	

	server.Run(":" + config.ServerPort)
	fmt.Println("App is running on port ", config.ServerPort)
}