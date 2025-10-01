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

func initReviewTestHandler(t *testing.T) (*gorm.DB, service.ReviewService, service.UserService, service.BookService, *handler.ReviewHandler) {
	testDB := db.SetupTestDB(t)
	reviewService := service.NewReviewService(testDB)
	userService := service.NewUserService(testDB)
	bookService := service.NewBookService(testDB)
	reviewHandler := handler.NewReviewHandler(reviewService)

	return testDB, reviewService, userService, bookService, reviewHandler
}

func createTestReviewBook(t *testing.T, testDB *gorm.DB, title, author string, price float32) model.Book {
	t.Helper()

	book := model.Book{Title: title, Author: author, Description: "desc", Price: price}
	if err := testDB.Create(&book).Error; err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	return book
}

func createTestReviewUser(t *testing.T, testDB *gorm.DB, username, email string) model.User {
	t.Helper()

	user := model.User{Username: username, Email: email, Password: "pass"}
	if err := testDB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	return user
}


func TestAddReviewHandler(t *testing.T) {
	testDB, _, _, _, reviewHandler := initReviewTestHandler(t)

	user := createTestReviewBook(t, testDB, "The Humilated and Insulted", "Dostoevsky", 20.00)
	book := createTestReviewUser(t, testDB, "testReviewUser1", "testReviewUser1@gmail.com")

	testCases := []struct {
		name		string
		urlPath		string
		reqBody		model.ReviewRequest
		wantStatus	int
		wantResp	string
	}{
		{"test empty text", fmt.Sprintf("/review/add/%d/%d", user.ID, book.ID), model.ReviewRequest{Text: ""}, http.StatusBadRequest, service.ErrEmptyReview.Error()},
		{"test book not found", fmt.Sprintf("/review/add/%d/%d", user.ID, 9999), model.ReviewRequest{Text: "Great book!"}, http.StatusNotFound, service.ErrBookNotFound.Error()},
		{"test user not found", fmt.Sprintf("/review/add/%d/%d", 9999, book.ID), model.ReviewRequest{Text: "Great book!"}, http.StatusNotFound, service.ErrUserNotFound.Error()},
		{"test review added successfully", fmt.Sprintf("/review/add/%d/%d", user.ID, book.ID), model.ReviewRequest{Text: "Great book!"}, http.StatusCreated, "Review added!"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.reqBody)
			req := httptest.NewRequest(http.MethodPost, test.urlPath, bytes.NewReader(body))
			w := httptest.NewRecorder()

			reviewHandler.AddReview(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantResp)
		})
	}
}


func TestGetReviewsByBookHandler(t *testing.T) {
	testDB, reviewService, _, _, reviewHandler := initReviewTestHandler(t)

	user := createTestReviewBook(t, testDB, "The Cousine Bette", "Balzac", 20.00)
	book := createTestReviewUser(t, testDB, "testReviewUser2", "testReviewUser2@gmail.com")

	err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "Excellent book!"})
	if err != nil {
		t.Fatalf("failed to add book review: %v", err)
	}

	testCases := []struct{
		name		string
		urlPath		string
		wantStatus	int
		wantResp	string
	}{
		{"test valid book with reviews", fmt.Sprintf("/review/get/%d", book.ID), http.StatusOK, "Excellent book!"},
		{"test book with no reviews", "/review/get/9999", http.StatusOK, "[]"},
		{"test invalid id in path", "/review/get/invalidId", http.StatusBadRequest, helper.ErrInvalidId.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.urlPath, nil)
			w := httptest.NewRecorder()

			reviewHandler.GetReviewsByBook(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantResp)
		})
	}
}


func TestGetReviewsByUserHandler(t *testing.T) {
	testDB, reviewService, _, _, reviewHandler := initReviewTestHandler(t)

	user := createTestReviewUser(t, testDB, "testUser", "testUser@gmail.com")
	book := createTestReviewBook(t, testDB, "The Idiot", "Dostoevsky", 15.00)

	err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "Very deep and psychological"})
	if err != nil {
		t.Fatalf("failed to add user review: %v", err)
	}

	testCases := []struct{
		name		string
		urlPath		string
		wantStatus	int
		wantResp	string
	}{
		{"test valid user with reviews", fmt.Sprintf("/review/getByUser/%d", user.ID), http.StatusOK, "Very deep and psychological"},
		{"test user with no reviews", "/review/getByUser/9999", http.StatusOK, "[]"},
		{"test invalid id in path", "/review/getByUser/invalidId", http.StatusBadRequest, helper.ErrInvalidId.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.urlPath, nil)
			w := httptest.NewRecorder()

			reviewHandler.GetReviewsByUser(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantResp)
		})
	}
}


func TestUpdateReviewHandler(t *testing.T) {
	testDB, reviewService, _, _, reviewHandler := initReviewTestHandler(t)

	user := createTestReviewUser(t, testDB, "testUser", "testUser@gmail.com")
	book := createTestReviewBook(t, testDB, "Demons", "Dostoevsky", 15.00)

	err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "Very deep and psychological"})
	if err != nil {
		t.Fatalf("failed to add user review: %v", err)
	}

	testCases := []struct{
		name			string
		urlPath			string
		reviewRequest	model.ReviewRequest
		wantStatus		int
		wantResp		string
	}{
		{"test valid update", fmt.Sprintf("/review/update/%d/%d", user.ID, book.ID), model.ReviewRequest{Text: "The hardest book to read by Dostoevsky"}, http.StatusOK, "Review updated!"},
		{"test review not found", "/review/update/9999/9999", model.ReviewRequest{Text: "The hardest book to read by Dostoevsky"}, http.StatusNotFound, service.ErrReviewNotFound.Error()},
		{"test invalid id in path", "/review/update/9999/invalidId", model.ReviewRequest{Text: "The hardest book to read by Dostoevsky"}, http.StatusBadRequest, helper.ErrInvalidId.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.reviewRequest)
			req := httptest.NewRequest(http.MethodPatch, test.urlPath, bytes.NewReader(body))
			w := httptest.NewRecorder()

			reviewHandler.UpdateReview(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantResp)
		})
	}
}


func TestDeleteReviewByIDHandler(t *testing.T){
	testDB, reviewService, _, _, reviewHandler := initReviewTestHandler(t)

	user := createTestReviewUser(t, testDB, "testUser", "testUser@gmail.com")
	book := createTestReviewBook(t, testDB, "Ana Karenina", "Tolstoy", 15.00)

	err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "The best book by Tolstoy"})
	if err != nil {
		t.Fatalf("failed to add user review: %v", err)
	}

	var review model.Review
	if err := testDB.Where("user_id = ? AND book_id = ?", user.ID, book.ID).First(&review).Error; err != nil {
		t.Fatalf("failed to fetch review: %v", err)
	}

	testCases := []struct{
		name		string
		urlPath		string
		wantStatus	int
		wantResp	string
	}{
		{"test delete existing review", fmt.Sprintf("/review/delete/%d", review.ID), http.StatusOK, "Review deleted!"},
		{"test delete non-existing review", "/review/delete/9999", http.StatusNotFound, service.ErrReviewNotFound.Error()},
		{"test invalid id in path", "/review/delete/invalidId", http.StatusBadRequest, helper.ErrInvalidId.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, test.urlPath, nil)
			w := httptest.NewRecorder()

			reviewHandler.DeleteReviewByID(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantResp)
		})
	}
}