package services

import (
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"errors"
	"testing"

	"gorm.io/gorm"
)

func initOrderTestServices(t *testing.T) (*gorm.DB, service.OrderService, service.BookService, service.CartService) {
	testDB := db.SetupTestDB(t)
	orderService := service.NewOrderService(testDB)
	bookService := service.NewBookService(testDB)
	cartService := service.NewCartService(testDB)

	return testDB, orderService, bookService, cartService
}

func createTestOrderUser(t *testing.T, testDB *gorm.DB, username, email string) model.User {
	user := model.User{Username: username, Password: "1234", Email: email}
	if err := testDB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	return user
}

func createTestOrderBook(t *testing.T, testDB *gorm.DB, title, author string, price float32) model.Book {
	book := model.Book{Title: title, Author: author, Price: price}
	if err := testDB.Create(&book).Error; err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	return book
}

func addBookToCart(t *testing.T, cartService service.CartService, userID, bookID uint, quantity int) {
	if err := cartService.AddToCart(userID, bookID, quantity); err != nil {
		t.Fatalf("failed to add book to cart: %v", err)
	}
}


func TestCreateOrder(t *testing.T) {
	testDB, orderService, _, cartService := initOrderTestServices(t)

	user := createTestOrderUser(t, testDB, "orderUser1", "orderUser1@gmail.com")
	book := createTestOrderBook(t, testDB, "Crime and Punishment", "Dostoevsky", 30)

	addBookToCart(t, cartService, user.ID, book.ID, 2)

	testCases := []struct {
		name        string
		userID      uint
		address     string
		wantErr     error
		wantTotal   float32
		wantBookQty int
	}{
		{"test successfully create order", user.ID, "123 Main St", nil, 60.0, 2},
		{"test empty cart", 9999, "123 Main St", service.ErrEmptyCart, 0, 0},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := orderService.CreateOrder(test.userID, test.address)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("expected error %v, got %v", test.wantErr, err)
			}
			if err == nil {
				if resp.Total != test.wantTotal {
					t.Errorf("expected total %.2f, got %.2f", test.wantTotal, resp.Total)
				}
				if len(resp.Books) == 0 || resp.Books[0].Quantity != test.wantBookQty {
					t.Errorf("expected book quantity %d, got %d", test.wantBookQty, resp.Books[0].Quantity)
				}
			}
		})
	}
}


func TestCancelOrder(t *testing.T) {
	testDB, orderService, _, cartService := initOrderTestServices(t)

	user := createTestOrderUser(t, testDB, "cancelUser", "cancelUser@gmail.com")
	book := createTestOrderBook(t, testDB, "The Idiot", "Dostoevsky", 20)

	addBookToCart(t, cartService, user.ID, book.ID, 1)

	orderResp, _ := orderService.CreateOrder(user.ID, "Addr")
	orderService.UpdateStatus(orderResp.ID, "completed")

	testCases := []struct {
		name    string
		orderID uint
		wantErr error
	}{
		{"test cancel pending order", orderResp.ID, service.ErrOrderCancel},
		{"test cancel non-existing order", 9999, service.ErrOrderNotFound},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := orderService.CancelOrder(test.orderID)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
		})
	}
}


func TestGetOrder(t *testing.T) {
	testDB, orderService, _, cartService := initOrderTestServices(t)

	user := createTestOrderUser(t, testDB, "getOrderUser", "getOrderUser@gmail.com")
	book := createTestOrderBook(t, testDB, "Братя Карамазови", "Фьодор Достоевски", 10)

	addBookToCart(t, cartService, user.ID, book.ID, 1)

	orderResp, _ := orderService.CreateOrder(user.ID, "Addr")

	testCases := []struct {
		name      string
		orderID   uint
		wantErr   error
		wantTitle string
	}{
		{"test get existing order", orderResp.ID, nil, "Братя Карамазови"},
		{"test get non-existing order", 9999, service.ErrOrderNotFound, ""},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := orderService.GetOrder(test.orderID)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
			if err == nil && resp.Books[0].Title != test.wantTitle {
				t.Errorf("expected book %s, got %s", test.wantTitle, resp.Books[0].Title)
			}
		})
	}
}


func TestUpdateStatus(t *testing.T) {
	testDB, orderService, _, cartService := initOrderTestServices(t)

	user := createTestOrderUser(t, testDB, "updateUser", "updateUser@gmail.com")
	book := createTestOrderBook(t, testDB, "Бесове", "Фьодор Михайлович Достоевски", 12)

	addBookToCart(t, cartService, user.ID, book.ID, 1)

	orderResp, _ := orderService.CreateOrder(user.ID, "Addr")

	testCases := []struct {
		name    string
		orderID uint
		status  string
		wantErr error
	}{
		{"test update status", orderResp.ID, "shipped", nil},
		{"test update non-existing order", 9999, "shipped", service.ErrOrderNotFound},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := orderService.UpdateStatus(test.orderID, test.status)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
		})
	}
}


func TestGetUserOrders(t *testing.T) {
	testDB, orderService, _, cartService := initOrderTestServices(t)

	user1 := createTestOrderUser(t, testDB, "user1", "user1@gmail.com")
	user2 := createTestOrderUser(t, testDB, "user2", "user2@gmail.com")

	bookA := createTestOrderBook(t, testDB, "Book A", "Author A", 10)
	bookB := createTestOrderBook(t, testDB, "Book B", "Author B", 20)

	addBookToCart(t, cartService, user1.ID, bookA.ID, 1)
	addBookToCart(t, cartService, user1.ID, bookB.ID, 2)
	_, _ = orderService.CreateOrder(user1.ID, "Addr1")

	addBookToCart(t, cartService, user2.ID, bookB.ID, 3)
	_, _ = orderService.CreateOrder(user2.ID, "Addr2")

	testCases := []struct {
		name        string
		userID      uint
		wantErr     error
		wantOrders  int
		wantFirst   string
	}{
		{"test user1 has orders", user1.ID, nil, 1, "Book A"},
		{"test user2 has orders", user2.ID, nil, 1, "Book B"},
		{"test user with no orders", 9999, nil, 0, ""},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			orders, err := orderService.GetUserOrders(test.userID)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
			if len(orders) != test.wantOrders {
				t.Errorf("expected %d orders, got %d", test.wantOrders, len(orders))
			}
			if test.wantOrders > 0 && orders[0].Books[0].Title != test.wantFirst {
				t.Errorf("expected first book %s, got %s", test.wantFirst, orders[0].Books[0].Title)
			}
		})
	}
}


func TestGetOrdersByStatus(t *testing.T) {
	testDB, orderService, _, cartService := initOrderTestServices(t)

	user := createTestOrderUser(t, testDB, "statusUser", "statusUser@gmail.com")
	bookA := createTestOrderBook(t, testDB, "Book A", "Author A", 10)
	bookB := createTestOrderBook(t, testDB, "Book B", "Author B", 20)

	addBookToCart(t, cartService, user.ID, bookA.ID, 1)
	order1, _ := orderService.CreateOrder(user.ID, "Addr1")
	orderService.UpdateStatus(order1.ID, "pending")

	addBookToCart(t, cartService, user.ID, bookB.ID, 2)
	order2, _ := orderService.CreateOrder(user.ID, "Addr2")
	orderService.UpdateStatus(order2.ID, "shipped")

	testCases := []struct {
		name        string
		status      string
		wantErr     error
		wantCount   int
		wantFirst   string
	}{
		{"test get pending orders", "pending", nil, 1, "Book A"},
		{"test get shipped orders", "shipped", nil, 1, "Book B"},
		{"test get non-existing status", "cancelled", nil, 0, ""},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			orders, err := orderService.GetOrdersByStatus(test.status)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
			if len(orders) != test.wantCount {
				t.Errorf("expected %d orders, got %d", test.wantCount, len(orders))
			}
			if test.wantCount > 0 && orders[0].Books[0].Title != test.wantFirst {
				t.Errorf("expected first book %s, got %s", test.wantFirst, orders[0].Books[0].Title)
			}
		})
	}
}