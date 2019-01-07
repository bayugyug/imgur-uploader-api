package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	dummyCredentials = ` {
        "client_id": "80473df3ff0641d",
        "client_secret": "f27a1350cb03410cbfc4fea0069201f2cb6cb93c",
        "code": "236fc06f1bea341b6a1baac9d6d366cabef47c6e",
        "bearer": {
            "access_token": "f1d9c3c2db6676118ce1ff7611b78480289556a7",
            "token_type": "bearer",
            "refresh_token": "d86b3dcc9169bffb5bd142a004b2cbaca6cc123c",
            "expiry": "2029-01-04T09:51:04.701868765+08:00"
        }
	}`

	dummyFdata = `{
                "urls": [
                        "https://farm3.staticflickr.com/2879/11234651086_681b3c2c00_b_d.jpg",
                        "https://farm4.staticflickr.com/3790/11244125445_3c2f32cd83_k_d.jpg"
                        ]
                }`

	testRawData = []struct {
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
			FormData: dummyFdata,
		},
	}
)

func TestSomeCredentials(t *testing.T) {
	//try
	ok, cfg := formatConfig("", dummyCredentials)
	if !ok || cfg == nil {
		t.Error("Credentials format invalid!")
	}
	t.Log("Credentials Ok")
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
