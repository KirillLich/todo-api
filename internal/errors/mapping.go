package errors

import (
	"errors"
	"net/http"
)

var ErrEmptyTitle = errors.New("empty title of todo")
var ErrNotFound = errors.New("todo not found")
var ErrInvalidId = errors.New("invalid id parametr")

func MapErrorHttpStatus(err error) (int, string) {
	if errors.Is(err, ErrEmptyTitle) {
		return http.StatusBadRequest, "empty title"
	}
	if errors.Is(err, ErrInvalidId) {
		return http.StatusBadRequest, "invalid id"
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, "not found"
	}
	return http.StatusInternalServerError, "internal server error"
}
