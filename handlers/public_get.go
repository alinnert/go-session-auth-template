package handlers

import (
	"net/http"
)

// GetPublic GET /public
func GetPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, "I'm a public route.")
	}
}
