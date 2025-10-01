package handler

import (
	"BookVault-API/helper"
	"BookVault-API/service"
	"net/http"
	"strconv"
)

type CartHandler struct {
	service		service.CartService
}

func NewCartHandler(s service.CartService) *CartHandler {
	return &CartHandler{service: s}
}


func (c *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPost) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 4, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, bookID := uint(IDs[0]), uint(IDs[1])

	q, ok := helper.RequiredQueryParam(w, r, "quantity")
	if !ok {
		return
	}
	
	quantity, err := strconv.Atoi(q)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid quantity")
		return
	}

	if err := c.service.AddToCart(userID, bookID, quantity); err != nil {
		switch err {
		case service.ErrBookNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Book added to cart!"})
}


func (c *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodDelete) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(IDs[0])

	if err := c.service.ClearCart(userID); err != nil {
		switch err {
		case service.ErrCartNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Cart has been cleared!"})
}


func (c *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodDelete) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 4, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, bookID := uint(IDs[0]), uint(IDs[1])

	if err := c.service.RemoveFromCart(userID, bookID); err != nil {
		switch err {
		case service.ErrCartNotFound, service.ErrCartBookNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Book removed from cart!"})
}


func (c *CartHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPatch) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 4, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, bookID := uint(IDs[0]), uint(IDs[1])

	q, ok := helper.RequiredQueryParam(w, r, "quantity")
	if !ok {
		return
	}

	quantity, err := strconv.Atoi(q)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid quantity")
		return
	}

	if err := c.service.UpdateQuantity(userID, bookID, quantity); err != nil {
		switch err {
		case service.ErrCartNotFound, service.ErrCartBookNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Book quantity updated!"})
}


func (c *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 2, 1)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	cartID := uint(IDs[0])

	cart, err := c.service.GetCart(cartID)
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, cart)
}