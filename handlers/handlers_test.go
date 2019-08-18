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

	"github.com/dgraph-io/badger"
)

type testRouteOptions struct {
	name, reqMethod, reqRoute string
	reqBody                   string
	contextMap                map[globals.ContextKey]interface{}
	expectedStatus            int
	expectedBody              string
	handlerFunc               http.HandlerFunc
}

// #region Helper functions
func flushDb(t *testing.T, db *badger.DB) {
	err := db.Update(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchValues = false
		iter := txn.NewIterator(options)
		defer iter.Close()

		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			err := txn.Delete(item.Key())
			if err != nil {
				t.Fatal(err)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
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

// #endregion Helper functions

func testRoute(
	t *testing.T, test *testRouteOptions,
) *httptest.ResponseRecorder {
	// #region Setup test
	var bodyReader io.Reader = nil
	if test.reqBody != "" {
		bodyReader = bytes.NewBufferString(test.reqBody)
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
	// #endregion Setup test

	// #region Check result
	if test.expectedStatus != 0 {
		if status := res.Code; status != test.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, test.expectedStatus)
		}
	}

	if ctype := res.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("content typ eheader does not match: got %v want %v",
			ctype, "application/json")
	}

	if test.expectedBody != "" {
		if body := strings.TrimSpace(res.Body.String()); body != test.expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, test.expectedBody)
		}
	}
	// #endregion Check result

	return res
}
