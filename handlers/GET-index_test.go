package handlers_test

import (
	"auth-server/handlers"
	"net/http"
	"testing"
)

func TestGetIndex(t *testing.T) {
	tests := []testRouteOptions{
		{
			name:           "Get index page (success)",
			reqMethod:      "GET",
			reqRoute:       "/",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok","data":"Demo Auth Server is running."}`,
			handlerFunc:    handlers.GetIndex(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testRoute(t, &test)
		})
	}
}
