package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var (
	ErrInvalidPath	= errors.New("invalid path")
	ErrInvalidId	= errors.New("invalid id")
)

func WriteJSON(w http.ResponseWriter, httpStatusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	json.NewEncoder(w).Encode(data)
}


func WriteError(w http.ResponseWriter, httpStatusCode int, message string) {
	WriteJSON(w, httpStatusCode, map[string]string{"error": message})
}


func RequiredMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return false
	}
	return true
}


func RequiredQueryParam(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		WriteError(w, http.StatusBadRequest, fmt.Sprintf("%s query param is required", key))
		return "", false
	}

	return value, true
}


func ParseIDsFromPath(r *http.Request, expectedParts, idStartIndex int) ([]int, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != expectedParts {
		return nil, ErrInvalidPath
	}

	IDs := make([]int, expectedParts - idStartIndex)
	for i := idStartIndex; i < expectedParts; i++ {
		id, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, ErrInvalidId
		}
		IDs[i - idStartIndex] = id
	}

	return IDs, nil
}


func AssertResponse(t *testing.T, w *httptest.ResponseRecorder, wantStatus int, wantSubstring string) {
	t.Helper()

	if w.Code != wantStatus {
		t.Errorf("expected status %d, got %d", wantStatus, w.Code)
	}

	respBody := w.Body.String()
	if !strings.Contains(respBody, wantSubstring) {
		t.Errorf("expected body to contain %q, got %s", wantSubstring, respBody)
	}
}