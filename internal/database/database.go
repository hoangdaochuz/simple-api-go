package database

import (
	"fmt"
	"log"

	"example.com/go-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



var DB *gorm.DB

func Connect(config *config.Config) {
	// Connect to database

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	 config.DBHost, config.DBPort, config.DBUser,config.DBPassword, config.DBName)
	
	var err error

	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{}) 
	if err != nil {
		log.Fatal("Could not connect to the database")
	}
	log.Println("Connected to the database")
}