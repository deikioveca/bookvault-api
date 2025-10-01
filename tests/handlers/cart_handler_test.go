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

func initCartTestHandler(t *testing.T) (*gorm.DB, service.UserService, service.CartService, service.BookService, *handler.CartHandler) {
	testDB := db.SetupTestDB(t)
	userService := service.NewUserService(testDB)
	cartService := service.NewCartService(testDB)
	bookService := service.NewBookService(testDB)
	cartHandler := handler.NewCartHandler(cartService)

	return testDB, userService, cartService, bookService, cartHandler
}

func createTestCartBook(t *testing.T, testDB *gorm.DB, bookService service.BookService, title, author string, price float32) model.Book {
	t.Helper()

	bookReq := model.BookRequest{
		Title:       title,
		Author:      author,
		Description: "Test book description",
		Price:       &price,
	}
	err := bookService.CreateBook(&bookReq)
	if err != nil {
		t.Fatalf("failed to create test book: %v", err)
	}

	var book model.Book
	if err := testDB.Where("title = ?", title).First(&book).Error; err != nil {
		t.Fatalf("failed to fetch test book: %v", err)
	}
	return book
}


func TestAddToCartHandler(t *testing.T) {
	testDB, userService, _, bookService, cartHandler := initCartTestHandler(t)

	user := createTestUser(t, testDB, userService, "cartUser1", "cartUser1@gmail.com")
	book := createTestCartBook(t, testDB, bookService, "Crime and Punishment", "Dostoevsky", 30.00)

	testCases := []struct{
		testName		string
		urlPath			string
		wantStatus		int
		wantRespBody	string
	}{
		{"test add to cart succeed", fmt.Sprintf("/cart/add/%d/%d?quantity=1", user.ID, book.ID), http.StatusCreated, "Book added to cart!"},
		{"test book not found", fmt.Sprintf("/cart/add/%d/%d?quantity=1", user.ID, 9999), http.StatusNotFound, service.ErrBookNotFound.Error()},
		{"test invalid quantity", fmt.Sprintf("/cart/add/%d/%d?quantity=invalidQuantity", user.ID, book.ID), http.StatusBadRequest, "invalid quantity"},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.urlPath, nil)
			w := httptest.NewRecorder()

			cartHandler.AddToCart(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestClearCartHandler(t *testing.T) {
	testDB, userService, cartService, bookService, cartHandler := initCartTestHandler(t)

	user := createTestUser(t, testDB, userService, "cartUser2", "cartUser2@gmail.com")
	book := createTestCartBook(t, testDB, bookService, "John Stainbeck", "East of Eden", 30.00)
	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct{
		testName		string
		userID			uint
		wantStatus		int
		wantRespBody	string
	}{
		{"test clear existing cart", user.ID, http.StatusOK, "Cart has been cleared!"},
		{"test clear non-existent cart", 9999, http.StatusNotFound, service.ErrCartNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			url := fmt.Sprintf("/cart/clear/%d", test.userID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()

			cartHandler.ClearCart(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestRemoveFromCart(t *testing.T) {
	testDB,userService,  cartService, bookService, cartHandler := initCartTestHandler(t)

	user := createTestUser(t,testDB, userService, "cartUser3", "cartUser3@gmail.com")
	book := createTestCartBook(t, testDB, bookService, "The Gambler", "Fyodor Dostoevsky", 15.00)
	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct {
		testName       	string
		userID     		uint
		bookID     		uint
		wantStatus 		int
		wantRespBody	string
	}{
		{"test remove existing book", user.ID, book.ID, http.StatusOK, "Book removed from cart!"},
		{"test remove from non-existent cart", 9999, book.ID, http.StatusNotFound, service.ErrCartNotFound.Error()},
		{"test remove non-existent book", user.ID, 9999, http.StatusNotFound, service.ErrCartBookNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			url := fmt.Sprintf("/cart/remove/%d/%d", test.userID, test.bookID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()

			cartHandler.RemoveFromCart(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestUpdateQuantityHandler(t *testing.T){
	testDB, userService, cartService, bookService, cartHandler := initCartTestHandler(t)

	user := createTestUser(t,testDB, userService, "cartUser4", "cartUser4@gmail.com")
	book := createTestCartBook(t, testDB, bookService, "Rudin", "Ivan Turgenev", 10.00)
	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct{
		testName		string
		urlPath			string
		wantStatus		int
		wantRespBody	string
	}{
		{"test update existing book", fmt.Sprintf("/cart/update/%d/%d?quantity=2", user.ID, book.ID), http.StatusOK, "Book quantity updated!"},
		{"test update missing cart", fmt.Sprintf("/cart/update/%d/%d?quantity=2", 9999, book.ID), http.StatusNotFound, service.ErrCartNotFound.Error()},
		{"test update missing book", fmt.Sprintf("/cart/update/%d/%d?quantity=2", user.ID, 9999), http.StatusNotFound, service.ErrCartBookNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, test.urlPath, nil)
			w := httptest.NewRecorder()

			cartHandler.UpdateQuantity(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestGetCartHandler(t *testing.T){
	testDB, userService,  cartService, bookService, cartHandler := initCartTestHandler(t)

	user := createTestUser(t,testDB, userService, "cartUser5", "cartUser5@gmail.com")
	book := createTestCartBook(t, testDB, bookService, "The Picture of Dorian Gray", "Oscar Wilde", 20.00)
	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct {
		testName       	string
		cartID     		uint
		wantStatus 		int
		wantRespBody	string
	}{
		{"test get existing cart", 1, http.StatusOK, "Oscar Wilde"},
		{"test get non-existent cart", 9999, http.StatusNotFound, service.ErrCartNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			url := fmt.Sprintf("/cart/%d", test.cartID)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			cartHandler.GetCart(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}