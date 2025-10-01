package service

import (
	"BookVault-API/jwt"
	"BookVault-API/model"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUsernameExist 		= errors.New("user with this username already exist")
	ErrEmailExist    		= errors.New("user with this email already exist")
	ErrInvalidCredentials 	= errors.New("invalid credentials")
	ErrEmptyFields			= errors.New("fields cannot be empty")
	ErrCreatingToken		= errors.New("could not create token")
	ErrUserNotFound			= errors.New("user not found")
	ErrUserDetailsNotFound	= errors.New("user details not found")
	ErrHashPassword			= errors.New("couldn't generate hash password")
)

type UserService interface {
	UsersExist() bool

	Register(registerRequest *model.RegisterRequest) error

	Login(loginRequest *model.LoginRequest) (string, error)

	CreateDetails(userID uint, userDetailsRequest *model.UserDetailsRequest) error
	
	GetUserByID(userID uint) (*model.UserResponse, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}


func (u *userService) UsersExist() bool {
	var count int64
	u.db.Model(&model.User{}).Count(&count)
	return count > 0
}


func (u *userService) Register(registerRequest *model.RegisterRequest) error {
	var user model.User

	if registerRequest.Username == "" || registerRequest.Password == "" || registerRequest.Email == "" {
		return ErrEmptyFields
	}

	if err := u.db.Where("username = ?", registerRequest.Username).First(&user).Error; err == nil {
		return ErrUsernameExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := u.db.Where("email = ?", registerRequest.Email).First(&user).Error; err == nil {
		return ErrEmailExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrHashPassword
	}

	role := "user"
	if !u.UsersExist() {
		role = "admin"
	}

	user.Username 	= registerRequest.Username
	user.Password 	= string(hashedPassword)
	user.Email 		= registerRequest.Email
	user.Role 		= role

	return u.db.Create(&user).Error
}


func (u *userService) Login(loginRequest *model.LoginRequest) (string, error) {
	var user model.User

	if loginRequest.Username == "" || loginRequest.Password == "" {
		return "", ErrEmptyFields
	}
	
	if err := u.db.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(loginRequest.Username, user.Role, 24 * time.Hour)
	if err != nil {
		return "", ErrCreatingToken
	}

	return token, nil
}


func (u *userService) CreateDetails(userID uint, userDetailsRequest *model.UserDetailsRequest) error {
	var user model.User

	if err := u.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if userDetailsRequest.FullName == "" || userDetailsRequest.PhoneNumber == "" {
		return ErrEmptyFields
	}

	userDetails := &model.UserDetails{
		UserID: 		userID,
		FullName: 		userDetailsRequest.FullName,
		PhoneNumber: 	userDetailsRequest.PhoneNumber,
	}
	return u.db.Create(&userDetails).Error
}


func (u *userService) GetUserByID(userID uint) (*model.UserResponse, error) {
	var user model.User
	
	if err := u.db.Preload("Details").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userResponse := model.UserResponse{
		Username: 		user.Username,
		Email: 			user.Email,
		PhoneNumber: 	user.Details.PhoneNumber,
		FullName: 		user.Details.FullName,
	}
	return &userResponse, nil
}