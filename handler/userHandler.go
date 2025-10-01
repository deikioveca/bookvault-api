package handler

import (
	"BookVault-API/helper"
	"BookVault-API/model"
	"BookVault-API/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}


func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPost) {
		return
	}

	var registerRequest model.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := u.service.Register(&registerRequest); err != nil {
		switch err {
		case service.ErrUsernameExist, service.ErrEmailExist, service.ErrEmptyFields:
			helper.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Registration completed!"})
}


func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPost) {
		return
	}

	var loginRequest model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := u.service.Login(&loginRequest)
	if err != nil {
		switch err {
		case service.ErrEmptyFields:
			helper.WriteError(w, http.StatusBadRequest, err.Error())
		case service.ErrInvalidCredentials:
			helper.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}


func (u *UserHandler) CreateDetails(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPost) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(IDs[0])
	
	var userDetailsRequest model.UserDetailsRequest

	if err := json.NewDecoder(r.Body).Decode(&userDetailsRequest); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := u.service.CreateDetails(userID, &userDetailsRequest); err != nil {
		switch err {
		case service.ErrUserNotFound, service.ErrEmptyFields:
			helper.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User details created!"})
}


func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(IDs[0])

	userDetails, err := u.service.GetUserByID(userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, userDetails)
}