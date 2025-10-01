package handler

import (
	"BookVault-API/helper"
	"BookVault-API/service"
	"net/http"
)

type OrderHandler struct {
	service	service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}


func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPost) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(IDs[0])

	address, ok := helper.RequiredQueryParam(w, r, "address")
	if !ok {
		return
	}

	order, err := o.service.CreateOrder(userID, address)
	if err != nil {
		switch err {
		case service.ErrEmptyCart:
			helper.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusCreated, order)
}


func (o *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPatch) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderID := uint(IDs[0])

	if err := o.service.CancelOrder(orderID); err != nil {
		switch err {
		case service.ErrOrderNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		case service.ErrOrderCancel:
			helper.WriteError(w, http.StatusForbidden, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Order cancelled!"})
}


func (o *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 2, 1)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderID := uint(IDs[0])

	order, err := o.service.GetOrder(orderID)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, order)
}


func (o *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := uint(IDs[0])

	userOrders, err := o.service.GetUserOrders(userID)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.WriteJSON(w, http.StatusOK, userOrders)
}


func (o *OrderHandler) GetOrdersByStatus(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	orderStatus, ok := helper.RequiredQueryParam(w, r, "status")
	if !ok {
		return
	}

	orders, err := o.service.GetOrdersByStatus(orderStatus)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.WriteJSON(w, http.StatusOK, orders)
}


func (o *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPatch) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderID := uint(IDs[0])

	status, ok := helper.RequiredQueryParam(w, r, "status")
	if !ok {
		return
	}

	if err := o.service.UpdateStatus(orderID, status); err != nil {
		switch err {
		case service.ErrOrderNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Order status updated!"})
}