package handler

import (
	"BookVault-API/helper"
	"net/http"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}


func (h *HomeHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	helper.WriteError(w, http.StatusNotFound, "endpoint not found")
}


func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Welcome user!"})
}


func (h *HomeHandler) AdminHome(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Welcome admin!"})
}