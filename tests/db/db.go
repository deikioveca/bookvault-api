package db

import (
	"BookVault-API/model"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	err := godotenv.Load(filepath.Join("..", "..", ".env"))
	if err != nil {
		t.Fatalf("no env file found: %v", err)
	}

	host 		:= 	os.Getenv("DB_HOST")
	port 		:= 	os.Getenv("DB_PORT")
	user 		:= 	os.Getenv("DB_USER")
	password 	:= 	os.Getenv("DB_PASSWORD")
	dbName 		:= 	os.Getenv("DB_NAME_TEST")
	sslMode 	:= 	os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.Migrator().DropTable(&model.User{}, &model.UserDetails{}, &model.Book{}, &model.Cart{}, &model.CartBook{}, &model.Order{}, &model.OrderBook{}, &model.Review{})
	if err != nil {
		t.Fatalf("failed to drop tables: %v", err)
	}
	
	err = db.AutoMigrate(&model.User{}, &model.UserDetails{}, &model.Book{}, &model.Cart{}, &model.CartBook{}, &model.Order{}, &model.OrderBook{}, &model.Review{})
	if err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	if err := db.Exec("TRUNCATE users, user_details, books, carts, cart_books, orders, order_books, reviews RESTART IDENTITY CASCADE;").Error; err != nil {
		t.Fatalf("failed to TRUNCATE db tables: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Exec("TRUNCATE users, user_details, books, carts, cart_books, orders, order_books, reviews RESTART IDENTITY CASCADE;").Error; err != nil {
			t.Fatalf("failed to TRUNCATE db tables: %v", err)
		}
	})

	return db
}