package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Get id from url
func GetUUID(r *http.Request) (uuid.UUID, error) {
	id := mux.Vars(r)["id"]
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}