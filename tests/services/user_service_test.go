package services

import (
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"errors"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func initUserTestServices(t *testing.T) (*gorm.DB, service.UserService) {
	testDB := db.SetupTestDB(t)
	return testDB, service.NewUserService(testDB)
}

func createTestUser(t *testing.T, username, email string) model.User {
	t.Helper()
	testDB, userService := initUserTestServices(t)

	err := userService.Register(&model.RegisterRequest{
		Username: username,
		Password: "1234",
		Email: email,
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	var user model.User
	if err := testDB.Where("username = ?", username).First(&user).Error; err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}

	return user
}


func TestRegister(t *testing.T) {
	_, userService := initUserTestServices(t)

	testCases := []struct{
		testName			string
		registerRequest		model.RegisterRequest
		wantErr				error
	}{
		{"test empty fields", model.RegisterRequest{Username: "", Password: "", Email: ""}, service.ErrEmptyFields},
		{"test registration succeeds", model.RegisterRequest{Username: "user1", Password: "1234", Email: "user1@gmail.com"}, nil},
		{"test duplicate username", model.RegisterRequest{Username: "user1", Password: "12345", Email: "user2@gmail.com"}, service.ErrUsernameExist},
		{"test duplicate email", model.RegisterRequest{Username: "user3", Password: "123456", Email: "user1@gmail.com"}, service.ErrEmailExist},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			err := userService.Register(&test.registerRequest)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("expected %v, got %v", test.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("expected success, got %v", err)
			}
		})
	}
}


func TestLogin(t *testing.T) {
	_, userService := initUserTestServices(t)

	user := createTestUser(t, "deikioveca", "deikioveca@gmail.com")

	testCases := []struct{
		testName		string
		loginRequest	model.LoginRequest
		wantErr			error
		valid			bool
	}{
		{"empty fields", model.LoginRequest{Username: "", Password: ""}, service.ErrEmptyFields, false},
		{"test wrong username", model.LoginRequest{Username: "wrongUser", Password: "1234"}, service.ErrInvalidCredentials, false},
		{"test wrong password", model.LoginRequest{Username: user.Username, Password: "wrongPassword"}, service.ErrInvalidCredentials, false},
		{"test login succeeded", model.LoginRequest{Username: user.Username, Password: "1234"}, nil, true},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			token, err := userService.Login(&test.loginRequest)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("expected %v, got %v", test.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected success, got %v", err)
			}
			if test.valid && token == "" {
				t.Errorf("expected non-empty token")
			}
		})
	}
}


func TestCreateDetails(t *testing.T) {
	_, userService := initUserTestServices(t)

	user := createTestUser(t, "detailsUser", "detailsUser@gmail.com")

	testCases := []struct{
		testName			string
		userID				uint
		userDetailsRequest	model.UserDetailsRequest
		wantErr				error
	}{
		{"test user not found", 9999, model.UserDetailsRequest{FullName: "FullName", PhoneNumber: "PhoneNumber"}, service.ErrUserNotFound},
		{"test empty fields", user.ID, model.UserDetailsRequest{FullName: "", PhoneNumber: ""}, service.ErrEmptyFields},
		{"test create details succeeded", user.ID, model.UserDetailsRequest{FullName: "Nikolay Ivanov Nikolaev", PhoneNumber: "0899373708"}, nil},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			err := userService.CreateDetails(test.userID, &test.userDetailsRequest)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("expected %v, got %v", test.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected success, got %v", err)
			}
		})
	}
}


func TestGetUserByID(t *testing.T) {
	_, userService := initUserTestServices(t)

	user := createTestUser(t, "getUser", "getUser@gmail.com")

	if err := userService.CreateDetails(user.ID, &model.UserDetailsRequest{
		FullName: "Nikolay Ivanov Nikolaev",
		PhoneNumber: "0899373708",
	}); err != nil {
		t.Fatalf("failed to create user details: %v", err)
	}

	testCases := []struct{
		testName		string
		userID			uint
		wantErr			error
		userResponse	*model.UserResponse
	}{
		{"test user not found", 9999, service.ErrUserNotFound, nil},
		{"test get user succeeded", user.ID, nil, &model.UserResponse{Username: "getUser", Email: "getUser@gmail.com", PhoneNumber: "0899373708", FullName: "Nikolay Ivanov Nikolaev"}},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			resp, err := userService.GetUserByID(test.userID)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("expected %v, got %v", test.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected success, got %v", err)
			}
			if !reflect.DeepEqual(resp, test.userResponse) {
				t.Errorf("expected %+v, got %+v", test.userResponse, resp)
			}
		})
	}
}