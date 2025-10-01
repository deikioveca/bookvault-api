package handlers

import (
	"BookVault-API/handler"
	"BookVault-API/helper"
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"
)

func float32Ptr(f float32) *float32 {
	return &f
}

func initBookTestHandler(t *testing.T) (*gorm.DB, service.BookService, *handler.BookHandler) {
	testDB := db.SetupTestDB(t)
	bookService := service.NewBookService(testDB)
	bookHandler := handler.NewBookHandler(bookService)
	return testDB, bookService, bookHandler
}

func createTestBook(t *testing.T,testDB *gorm.DB, bookService service.BookService, title, author string, price float32) model.Book {
	t.Helper()

	bookReq := &model.BookRequest{
		Title: title,
		Author: author,
		Description: "Test book description",
		Price: &price,
	}

	err := bookService.CreateBook(bookReq)
	if err != nil {
		t.Fatalf("failed to create test book: %v", err)
	}

	var book model.Book
	if err := testDB.Where("title = ?", title).First(&book).Error; err != nil {
		t.Fatalf("failed to fetch book: %v", err)
	}

	return book
}


func TestCreateBookHandler(t *testing.T) {
	_, _, bookHandler := initBookTestHandler(t)

	testCases := []struct {
		name			string
		reqBody			model.BookRequest
		wantStatus		int
		wantRespBody	string
	}{
		{"test empty fields", model.BookRequest{Title: "", Author: "", Description: "", Price: nil}, http.StatusBadRequest, service.ErrEmptyFields.Error()},
		{"test valid book", model.BookRequest{Title: "The Idiot", Author: "Dostoevsky", Description: "One of the 5 major books from Dostoevsky", Price: float32Ptr(20.00)}, http.StatusCreated, "Book created!"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/book/create", bytes.NewReader(body))
			w := httptest.NewRecorder()

			bookHandler.CreateBook(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestGetByTitleHandler(t *testing.T) {
	testDB, bookService, bookHandler := initBookTestHandler(t)

	createTestBook(t, testDB, bookService, "The Idiot", "Dostoevsky", 20.00)

	testCases := []struct{
		name			string
		urlPath			string
		wantStatus		int
		wantRespBody	string
	}{
		{"test book not found", "/book?title=Demons", http.StatusNotFound, service.ErrBookNotFound.Error()},
		{"test valid title", "/book?title=The%20Idiot", http.StatusOK, "The Idiot"},
		{"test missing query param title", "/book", http.StatusBadRequest, "title query param is required"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.urlPath, nil)
			w := httptest.NewRecorder()

			bookHandler.GetByTitle(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestGetBooksHandler(t *testing.T) {
	testDB, bookService, bookHandler := initBookTestHandler(t)
	
	testCases := []struct{
		testName		string
		prepareDB		func()
		wantStatus		int
		wantRespBody	string
	}{
		{"test no books", func(){}, http.StatusNotFound, service.ErrNoBooks.Error()},
		{"test books exist", func(){createTestBook(t, testDB, bookService, "The Idiot", "Dostoevsky", 20.00)}, http.StatusOK, "The Idiot"},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			test.prepareDB()
			req := httptest.NewRequest(http.MethodGet, "/book/all", nil)
			w := httptest.NewRecorder()

			bookHandler.GetBooks(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestGetBooksByAuthorHandler(t *testing.T) {
	testDB, bookService, bookHandler := initBookTestHandler(t)

	createTestBook(t, testDB, bookService, "The Idiot", "Dostoevsky", 20.00)

	testCases := []struct{
		 name			string
		 urlPath		string
		 wantStatus		int
		 wantRespBody	string
	}{
		{"test no books found", "/book/author?author=Stainbeck", http.StatusOK, ""},
		{"test valid author", "/book/author?author=Dostoevsky", http.StatusOK, "The Idiot"},
		{"test missing query param author", "/book/author", http.StatusBadRequest, "author query param is required"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.urlPath, nil)
			w := httptest.NewRecorder()

			bookHandler.GetBooksByAuthor(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestUpdateStockHandler(t *testing.T) {
	testDB, bookService, bookHandler := initBookTestHandler(t)

	book := createTestBook(t, testDB, bookService, "The Idiot", "Dostoevsky", 20.00)

	testCases := []struct{
		name			string
		urlPath			string
		wantStatus		int
		wantRespBody	string
	}{
		{"test book not found", "/book/updateStock/9999", http.StatusNotFound, service.ErrBookNotFound.Error()},
		{"test book stock updated", fmt.Sprintf("/book/updateStock/%d", book.ID), http.StatusOK, "Book stock updated!"},
		{"test missing id in path", "/book/updateStock/", http.StatusBadRequest, helper.ErrInvalidPath.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, test.urlPath, nil)
			w := httptest.NewRecorder()

			bookHandler.UpdateStock(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}