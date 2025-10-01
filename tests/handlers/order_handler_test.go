package handlers

import (
	"BookVault-API/handler"
	"BookVault-API/helper"
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"
)

func initOrderTestHandler(t *testing.T) (*gorm.DB, service.OrderService, service.CartService, *handler.OrderHandler) {
	testDB := db.SetupTestDB(t)
	cartService := service.NewCartService(testDB)
	orderService := service.NewOrderService(testDB)
	orderHandler := handler.NewOrderHandler(orderService)

	return testDB, orderService, cartService, orderHandler
}

func createTestOrderBook(t *testing.T, testDB *gorm.DB, title, author string, price float32) model.Book {
	t.Helper()

	book := model.Book{Title: title, Author: author, Description: "desc", Price: price}
	if err := testDB.Create(&book).Error; err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	return book
}

func createTestOrderUser(t *testing.T, testDB *gorm.DB, username, email string) model.User {
	t.Helper()

	user := model.User{Username: username, Email: email, Password: "pass"}
	if err := testDB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	return user
}

func addBookToCart(t *testing.T, cartService service.CartService, userID, bookID uint, quantity int) {
	t.Helper()

	if err := cartService.AddToCart(userID, bookID, quantity); err != nil {
		t.Fatalf("failed to add book to cart: %v", err)
	}
}


func TestCreateOrderHandler(t *testing.T) {
	testDB, _, cartService, orderHandler := initOrderTestHandler(t)
	user := createTestOrderUser(t, testDB, "orderUser1", "orderUser1@gmail.com")
	book := createTestOrderBook(t, testDB, "Book 1", "Author 1", 10)

	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{"test create order success", fmt.Sprintf("/order/create/%d?address=Addr1", user.ID), http.StatusCreated, "Book 1"},
		{"test empty cart", fmt.Sprintf("/order/create/%d?address=Addr1", 9999), http.StatusBadRequest, service.ErrEmptyCart.Error()},
		{"test missing address", fmt.Sprintf("/order/create/%d", user.ID), http.StatusBadRequest, "address query param is required"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.urlPath, nil)
			w := httptest.NewRecorder()

			orderHandler.CreateOrder(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantBody)
		})
	}
}


func TestCancelOrderHandler(t *testing.T) {
	testDB, orderService, cartService, orderHandler := initOrderTestHandler(t)

	user := createTestOrderUser(t, testDB, "orderUser2", "orderUser2@gmail.com")
	book := createTestOrderBook(t, testDB, "Book 2", "Author 2", 15)

	addBookToCart(t, cartService, user.ID, book.ID, 1)

	order, _ := orderService.CreateOrder(user.ID, "Addr2")

	_ = orderService.UpdateStatus(order.ID, "pending") 

	testCases := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{"test cancel order success", fmt.Sprintf("/order/cancel/%d", order.ID), http.StatusOK, "Order cancelled!"},
		{"test not found", "/order/cancel/9999", http.StatusNotFound, service.ErrOrderNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, test.urlPath, nil)
			w := httptest.NewRecorder()

			orderHandler.CancelOrder(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantBody)
		})
	}
}


func TestGetOrderHandler(t *testing.T) {
	testDB, orderService, cartService, orderHandler := initOrderTestHandler(t)

	user := createTestOrderUser(t, testDB, "orderUser3", "orderUser3@gmail.com")
	book := createTestOrderBook(t, testDB, "Book 3", "Author 3", 20)

	addBookToCart(t, cartService, user.ID, book.ID, 1)

	order, _ := orderService.CreateOrder(user.ID, "Addr3")

	testCases := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{"test get existing order", fmt.Sprintf("/order/%d", order.ID), http.StatusOK, "Book 3"},
		{"test get non-existing order", "/order/9999", http.StatusNotFound, service.ErrOrderNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.urlPath, nil)
			w := httptest.NewRecorder()

			orderHandler.GetOrder(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantBody)
		})
	}
}


func TestGetUserOrdersHandler(t *testing.T) {
	testDB, orderService, cartService, orderHandler := initOrderTestHandler(t)

	user := createTestOrderUser(t, testDB, "orderUser4", "orderUser4@gmail.com")
	book1 := createTestOrderBook(t, testDB, "Book 4A", "Author 4", 10)
	book2 := createTestOrderBook(t, testDB, "Book 4B", "Author 4", 15)

	addBookToCart(t, cartService, user.ID, book1.ID, 1)
	addBookToCart(t, cartService, user.ID, book2.ID, 1)

	_, _ = orderService.CreateOrder(user.ID, "Addr4")
	_, _ = orderService.CreateOrder(user.ID, "Addr4")

	testCases := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{"test existing user orders", fmt.Sprintf("/order/user/%d", user.ID), http.StatusOK, "Book 4A"},
		{"test non-existing user", "/order/user/9999", http.StatusOK, ""},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.urlPath, nil)
			w := httptest.NewRecorder()

			orderHandler.GetUserOrders(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantBody)
		})
	}
}


func TestGetOrdersByStatusHandler(t *testing.T) {
	testDB, orderService, cartService, orderHandler := initOrderTestHandler(t)

	user := createTestOrderUser(t, testDB, "orderUser5", "orderUser5@gmail.com")
	book := createTestOrderBook(t, testDB, "Book 5", "Author 5", 25)

	addBookToCart(t, cartService, user.ID, book.ID, 1)

	_, _ = orderService.CreateOrder(user.ID, "Addr5")
	_ = orderService.UpdateStatus(1, "shipped")

	testCases := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{"test existing status", "/order/status?status=shipped", http.StatusOK, "Book 5"},
		{"test non-existing status", "/order/status?status=cancelled", http.StatusOK, ""},
		{"test missing status param", "/order/status", http.StatusBadRequest, "status query param is required"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.urlPath, nil)
			w := httptest.NewRecorder()

			orderHandler.GetOrdersByStatus(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantBody)
		})
	}
}


func TestUpdateStatusHandler(t *testing.T) {
	testDB, orderService, cartService, orderHandler := initOrderTestHandler(t)

	user := createTestOrderUser(t, testDB, "orderUser6", "orderUser6@gmail.com")
	book := createTestOrderBook(t, testDB, "Book 6", "Author 6", 30)

	addBookToCart(t, cartService, user.ID, book.ID, 1)

	order, _ := orderService.CreateOrder(user.ID, "Addr6")

	testCases := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{"test update existing order", fmt.Sprintf("/order/update/%d?status=shipped", order.ID), http.StatusOK, "Order status updated!"},
		{"test update non-existing order", "/order/update/9999?status=shipped", http.StatusNotFound, service.ErrOrderNotFound.Error()},
		{"test missing status param", fmt.Sprintf("/order/update/%d", order.ID), http.StatusBadRequest, "status query param is required"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, test.urlPath, nil)
			w := httptest.NewRecorder()

			orderHandler.UpdateStatus(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantBody)
		})
	}
}