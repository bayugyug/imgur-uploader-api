package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testRawData = []struct {
	Method   string
	URL      string
	Expected string
}{
	{
		Method:   "GET",
		URL:      "/v1/api/images",
		Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
	},
	{
		Method: "GET",

		URL:      "/v1/api/images/upload/invalid-content-job-id",
		Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
	},
}

func TestSomeHandler(t *testing.T) {
	t.Log("Sanity checking ....")
	pUserCredentials = `{"client_id": "80473df3ff0641d", "client_secret": "f27a1350cb03410cbfc4fea0069201f2cb6cb93c"}`
	for _, rdata := range testRawData {
		req, err := http.NewRequest(rdata.Method, rdata.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAllImages)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		if strings.TrimSpace(rr.Body.String()) != rdata.Expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), rdata.Expected)
		}
	}
}
