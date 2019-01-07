package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testRawData = []struct {
	Method   string
	URL      string
	Expected string
	FormData string
	Handler  func(w http.ResponseWriter, r *http.Request)
}{
	{
		Method:   "GET",
		URL:      "/v1/api/images",
		Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
		Handler:  GetAllImages,
	},
	{
		Method:   "GET",
		URL:      "/v1/api/images/upload/invalid-content-job-id",
		Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
		Handler:  GetOneImage,
	},
	{
		Method:   "POST",
		URL:      "/v1/api/images/upload",
		Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
		Handler:  UploadImage,
	},
}

func TestSomeHandler(t *testing.T) {

	t.Log("Sanity checking ....")

	var body io.Reader
	for _, rdata := range testRawData {
		body = nil
		if rdata.Method == "POST" {
			body = bytes.NewBufferString(rdata.FormData)
		}

		//t.Log(rdata.Method, rdata.URL)

		req, err := http.NewRequest(rdata.Method, rdata.URL, body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(rdata.Handler)

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
