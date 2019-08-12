package handlers

import (
	"auth-server/services/validator"
	"net/http"
)

// handleValidationErrors handles validation errors.
// If one of the validated values is not valid this function writes a response.
// In that case the function returns true.
// If this function returns true, the handler should return and end the request.
// If this function returns false, the handler should countinue to handle the
// request as normal.
func handleValidationErrors(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch err := err.(type) {
	case validator.ValidationError:
		WriteErrorResponse(w, http.StatusBadRequest, err,
			"Request body is not valid.")
		return true
	default:
		WriteErrorResponse(w, http.StatusInternalServerError, err,
			"Error while validating request body.")
		return true
	}
}
