package database

import (
	"BookVault-API/model"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	host 		:= 	os.Getenv("DB_HOST")
	port 		:= 	os.Getenv("DB_PORT")
	user 		:= 	os.Getenv("DB_USER")
	password 	:= 	os.Getenv("DB_PASSWORD")
	dbName 		:= 	os.Getenv("DB_NAME")
	sslMode 	:= 	os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",host, port, user, password, dbName, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	db.Migrator().DropTable(&model.User{}, &model.UserDetails{}, &model.Book{}, &model.Cart{}, &model.CartBook{}, &model.Order{}, &model.OrderBook{}, &model.Review{})
	db.AutoMigrate(&model.User{}, &model.UserDetails{}, &model.Book{}, &model.Cart{}, &model.CartBook{}, &model.Order{}, &model.OrderBook{}, &model.Review{})

	return db
}