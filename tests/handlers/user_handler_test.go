package handlers

import (
	"BookVault-API/handler"
	"BookVault-API/helper"
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"gorm.io/gorm"
)


func initUserTestHandler(t *testing.T) (*gorm.DB, service.UserService, *handler.UserHandler) {
	testDB := db.SetupTestDB(t)
	userService := service.NewUserService(testDB)
	userHandler := handler.NewUserHandler(userService)
	return testDB, userService, userHandler
}


func createTestUser(t *testing.T, testDB *gorm.DB, userService service.UserService, username, email string) model.User {
	t.Helper()

	err := userService.Register(&model.RegisterRequest{
		Username: username,
		Password: "1234",
		Email: email,
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	var user model.User

	if err := testDB.Where("username = ?", username).First(&user).Error; err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}

	return user
}


func TestRegisterHandler(t *testing.T) {
	_, _, userHandler := initUserTestHandler(t)

	testCases := []struct{
		name			string
		reqBody			model.RegisterRequest
		wantStatus		int
		wantRespBody 	string
	}{
		{"test empty fields", model.RegisterRequest{Username: "", Password: "", Email: ""}, http.StatusBadRequest, service.ErrEmptyFields.Error()},
		{"test valid registration", model.RegisterRequest{Username: "deykio", Password: "1234", Email: "deykio@gmail.com"}, http.StatusCreated, "Registration completed!"},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			w := httptest.NewRecorder()

			userHandler.Register(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestLoginHandler(t *testing.T) {
	testDB, userService, userHandler := initUserTestHandler(t)
	
	createTestUser(t, testDB, userService, "validUser", "validUser@gmail.com")

	testCases := []struct{
		name			string
		reqBody			model.LoginRequest
		wantStatus		int
		wantRespBody 	string
	}{
		{"test valid login", model.LoginRequest{Username: "validUser", Password: "1234"}, http.StatusOK, `"token"`},
		{"test wrong password", model.LoginRequest{Username: "validUser", Password: "wrongPassword"}, http.StatusUnauthorized, service.ErrInvalidCredentials.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			w := httptest.NewRecorder()

			userHandler.Login(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestCreateDetailsHandler(t *testing.T) {
	testDB, userService, userHandler := initUserTestHandler(t)
	
	user := createTestUser(t,testDB, userService, "validUser1001", "validUser1001@gmail.com")

	testCases := []struct{
		name			string
		userID			uint
		reqBody			model.UserDetailsRequest
		wantStatus		int
		wantRespBody	string
	}{
		{"test valid details", user.ID, model.UserDetailsRequest{FullName: "Nikolay Ivanov Nikolaev", PhoneNumber: "0899373708"}, http.StatusCreated, "User details created!"},
		{"test empty fields", user.ID, model.UserDetailsRequest{FullName: "", PhoneNumber: ""}, http.StatusBadRequest, service.ErrEmptyFields.Error()},
		{"test user not found", 9999, model.UserDetailsRequest{FullName: "Nikolay Ivanov Nikolaev", PhoneNumber: "0899373708"}, http.StatusBadRequest, service.ErrUserNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.reqBody)
			url := "/user/createDetails/" + strconv.Itoa(int(test.userID))
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			w := httptest.NewRecorder()

			userHandler.CreateDetails(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}


func TestGetUserByIDHandler(t *testing.T) {
	testDB, userService, userHandler := initUserTestHandler(t)
	
	user := createTestUser(t,testDB, userService, "validUser2002", "validUser2002@gmail.com")
	
	userService.CreateDetails(user.ID, &model.UserDetailsRequest{
		FullName: "Ivan Petrov Dimitrov",
		PhoneNumber: "088735084",
	})

	testCases := []struct{
		name 			string
		userID			uint
		wantStatus		int
		wantRespBody	string
	}{
		{"test valid user", user.ID, http.StatusOK, user.Username},
		{"test user not found", 9999, http.StatusNotFound, service.ErrUserNotFound.Error()},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			url := "/user/getById/" + strconv.Itoa(int(test.userID))
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			userHandler.GetUserByID(w, req)

			helper.AssertResponse(t, w, test.wantStatus, test.wantRespBody)
		})
	}
}