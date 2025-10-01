package handler

import (
	"BookVault-API/helper"
	"BookVault-API/model"
	"BookVault-API/service"
	"encoding/json"
	"net/http"
)

type ReviewHandler struct {
	service service.ReviewService
}

func NewReviewHandler(s service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: s}
}


func (r *ReviewHandler) AddReview(w http.ResponseWriter, req *http.Request) {
	if !helper.RequiredMethod(w, req, http.MethodPost) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(req, 4, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, bookID := uint(IDs[0]), uint(IDs[1])

	var reviewRequest model.ReviewRequest
	if err := json.NewDecoder(req.Body).Decode(&reviewRequest); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.service.AddReview(userID, bookID, reviewRequest); err != nil {
		switch err {
		case service.ErrEmptyReview:
			helper.WriteError(w, http.StatusBadRequest, err.Error())
		case service.ErrUserNotFound, service.ErrBookNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Review added!"})
}


func (r *ReviewHandler) GetReviewsByBook(w http.ResponseWriter, req *http.Request) {
	if !helper.RequiredMethod(w, req, http.MethodGet) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(req, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	bookID := uint(IDs[0])

	reviews, err := r.service.GetReviewsByBook(bookID)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.WriteJSON(w, http.StatusOK, reviews)
}


func (r *ReviewHandler) GetReviewsByUser(w http.ResponseWriter, req *http.Request) {
	if !helper.RequiredMethod(w, req, http.MethodGet) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(req, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(IDs[0])

	reviews, err := r.service.GetReviewsByUser(userID)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.WriteJSON(w, http.StatusOK, reviews)
}


func (r *ReviewHandler) UpdateReview(w http.ResponseWriter, req *http.Request) {
	if !helper.RequiredMethod(w, req, http.MethodPatch) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(req, 4, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, bookID := uint(IDs[0]), uint(IDs[1])

	var reviewRequest model.ReviewRequest
	if err := json.NewDecoder(req.Body).Decode(&reviewRequest); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := r.service.UpdateReview(userID, bookID, reviewRequest); err != nil {
		switch err {
		case service.ErrReviewNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Review updated!"})
}


func (r *ReviewHandler) DeleteReviewByID(w http.ResponseWriter, req *http.Request) {
	if !helper.RequiredMethod(w, req, http.MethodDelete) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(req, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	reviewID := uint(IDs[0])

	if err := r.service.DeleteReviewByID(reviewID); err != nil {
		switch err {
		case service.ErrReviewNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Review deleted!"})
}