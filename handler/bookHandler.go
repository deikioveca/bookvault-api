package handler

import (
	"BookVault-API/helper"
	"BookVault-API/model"
	"BookVault-API/service"
	"encoding/json"
	"net/http"
)

type BookHandler struct {
	service	service.BookService
}

func NewBookHandler(s service.BookService) *BookHandler {
	return &BookHandler{service: s}
}


func (b *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPost) {
		return
	}

	var bookRequest model.BookRequest

	if err := json.NewDecoder(r.Body).Decode(&bookRequest); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := b.service.CreateBook(&bookRequest); err != nil {
		switch err {
		case service.ErrEmptyFields:
			helper.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Book created!"})
}


func (b *BookHandler) GetByTitle(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	bookTitle, ok := helper.RequiredQueryParam(w, r, "title")
	if !ok {
		return
	}
	
	book, err := b.service.GetByTitle(bookTitle)
	if err != nil {
		switch err {
		case service.ErrBookNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, book)
}


func (b *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet) {
		return
	}

	books, err := b.service.GetBooks()
	if err != nil {
		switch err {
		case service.ErrNoBooks:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, books)
}


func (b *BookHandler) GetBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodGet){
		return
	}

	bookAuthor, ok := helper.RequiredQueryParam(w, r, "author")
	if !ok {
		return
	}

	books, err := b.service.GetBooksByAuthor(bookAuthor)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.WriteJSON(w, http.StatusOK, books)
}


func (b *BookHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	if !helper.RequiredMethod(w, r, http.MethodPatch) {
		return
	}

	IDs, err := helper.ParseIDsFromPath(r, 3, 2)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	bookID := uint(IDs[0])

	if err := b.service.UpdateStock(bookID); err != nil {
		switch err {
		case service.ErrBookNotFound:
			helper.WriteError(w, http.StatusNotFound, err.Error())
		default:
			helper.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Book stock updated!"})
}