package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type errorListResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// WriteResponse writes the response if it's not an error.
func WriteResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	responseStruct := &response{Status: "ok"}

	if data != nil {
		responseStruct.Data = data
	}

	json.NewEncoder(w).Encode(responseStruct)
}

// WriteErrorResponse writes an error response and sets a status.
func WriteErrorResponse(
	w http.ResponseWriter, status int, err error, message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseStruct := &errorResponse{Status: "error", Message: message}

	if err != nil {
		responseStruct.Details = err.Error()
	}

	json.NewEncoder(w).Encode(responseStruct)
}

// WriteErrorListResponse writes an error response with a list of strings
// and sets a status.
func WriteErrorListResponse(
	w http.ResponseWriter, status int, errList error, message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errMessages := []string{}

	for _, err := range errList.(validator.ValidationErrors) {
		errMsgString := fmt.Sprintf(
			"Property '%s', failed requirement: %s %s",
			err.Field(), err.Tag(), err.Param(),
		)
		errMessages = append(errMessages, strings.Trim(errMsgString, " "))
	}

	responseStruct := &errorListResponse{"error", message, errMessages}
	json.NewEncoder(w).Encode(responseStruct)
}
