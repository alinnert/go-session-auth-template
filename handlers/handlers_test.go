package handlers_test

import (
	"auth-server/globals"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type handlerTestCase struct {
	reqMethod, reqRoute string
	reqBody             *string
	contextMap          map[globals.ContextKey]interface{}
	expectedStatus      int
	expectedBody        string
	handlerFunc         http.HandlerFunc
}

func applyContext(
	req *http.Request,
	items map[globals.ContextKey]interface{},
) *http.Request {
	ctx := req.Context()

	for key, item := range items {
		ctx = context.WithValue(ctx, key, item)
	}

	return req.WithContext(ctx)
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

	if test.contextMap != nil {
		req = applyContext(req, test.contextMap)
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
