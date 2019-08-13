package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type handlerTestCase struct {
	reqMethod, reqRoute string
	reqBody             *string
	expectedStatus      int
	expectedBody        string
	handlerFunc         http.HandlerFunc
}

func testRoute(t *testing.T, test handlerTestCase) {
	var bodyReader io.Reader = nil
	if test.reqBody != nil {
		bodyReader = bytes.NewBufferString(*test.reqBody)
	}

	req, err := http.NewRequest(
		test.reqMethod, test.reqRoute, bodyReader)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(test.handlerFunc)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != test.expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, test.expectedStatus)
	}

	if ctype := res.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("content typ eheader does not match: got %v want %v",
			ctype, "application/json")
	}

	if body := strings.TrimSpace(res.Body.String()); body != test.expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, test.expectedBody)
	}
}
