package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestServeHTTP tests the CORS ServeHTTP function, ensuring the headers are set
// 	to the correct values
func TestServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	rr := httptest.NewRecorder()

	corsTest := CORS{testHandler}
	corsTest.ServeHTTP(rr, req)

	exp := "Access-Control-Allow-Origin"
	got := "*"

	if rr.Header().Get(exp) != got {
		t.Errorf("Incorrect Access-Control-Allow-Origin set\nExpected: %s\nGot: %s", got, rr.Header().Get(exp))

	}

	exp = "Access-Control-Allow-Methods"
	got = "GET, PUT, POST, PATCH, DELETE"

	if rr.Header().Get(exp) != got {
		t.Errorf("Incorrect Access-Control-Allow-Methods set\nExpected: %s\nGot: %s", got, rr.Header().Get(exp))

	}

	exp = "Access-Control-Allow-Headers"
	got = "Content-Type, Authorization"

	if rr.Header().Get(exp) != got {
		t.Errorf("Incorrect Access-Control-Allow-Headers set\nExpected: %s\nGot: %s", got, rr.Header().Get(exp))

	}

	exp = "Access-Control-Expose-Headers"
	got = "Authorization"

	if rr.Header().Get(exp) != got {
		t.Errorf("Incorrect Access-Control-Expose-Headers set\nExpected: %s\nGot: %s", got, rr.Header().Get(exp))

	}

	exp = "Access-Control-Max-Age"
	got = "600"

	if rr.Header().Get(exp) != got {
		t.Errorf("Incorrect Access-Control-Max-Age set\nExpected: %s\nGot: %s", got, rr.Header().Get(exp))

	}

}
