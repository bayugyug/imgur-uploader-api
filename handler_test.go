package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

type ctxKey struct {
	name string
}

func (k ctxKey) String() string {
	return k.name
}

var (
	api              ApiHandler
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
		Ctx      context.Context
		Status   int
		FormData string
		Handler  func(w http.ResponseWriter, r *http.Request)
	}{
		{
			Method:   "GET",
			URL:      "/v1/api/images",
			Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
			Handler:  api.GetAllImages,
			Status:   200,
			Ctx:      context.WithValue(context.Background(), ctxKey{""}, ""),
		},
		{
			Method:   "GET",
			URL:      "/v1/api/images/upload/invalid-content-job-id",
			Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
			Handler:  api.GetOneImage,
			Status:   200,
			Ctx:      context.WithValue(context.Background(), ctxKey{""}, ""),
		},
		{
			Method:   "POST",
			URL:      "/v1/api/images/upload",
			Expected: `{"code":203,"message":"Non-Authoritative Information"}`,
			Handler:  api.UploadImage,
			FormData: dummyFdata,
			Status:   200,
			Ctx:      context.WithValue(context.Background(), ctxKey{""}, ""),
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
		req = req.WithContext(rdata.Ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(rdata.Handler)
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != rdata.Status {
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

func TestMockRoutes(t *testing.T) {

	var (
		mockUploadById = `{"id":"551046a9-7448-466d-ad9e-4e870a8e4a59-20190107-130152","created":"2019-01-07T13:01:52+08:00","finished":"2019-01-07T13:02:04+08:00","status":"complete","uploaded":{"complete":["https://i.imgur.com/MWfMA9e.jpg","https://i.imgur.com/thUkhoz.jpg"],"pending":null,"failed":null}}`
		mockUploadOk   = `{"jobId":"bbeb5ebd-401c-4ea3-ae9c-b7de8ee3482b-20190107-110300"}`
		mockUserAuth   = `{"code": 202,"message": "Accepted"}`
		mockUploadList = `"uploaded": ["https://i.imgur.com/8yc2oCz.jpg","https://i.imgur.com/u6GIFQA.jpg"]`
	)

	getAllImages := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockUploadList))
	})
	getOneImageById := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = chi.URLParam(r, "id")
		w.Write([]byte(mockUploadById))
	})
	uploadImages := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var images ParamImageURLS
		err := json.NewDecoder(r.Body).Decode(&images)
		if err != nil {
			return
		}
		//just in case :-)
		defer r.Body.Close()
		w.Write([]byte(mockUploadOk))
	})
	getUserAuth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = chi.URLParam(r, "code")
		w.Write([]byte(mockUserAuth))
	})
	r := chi.NewRouter()
	r.Get("/v1/api/images", getAllImages)
	r.Get("/v1/api/images/upload/{id}", getOneImageById)
	r.Post("/v1/api/images/upload", uploadImages)
	r.Get("/v1/api/credentials/{code}", getUserAuth)

	ts := httptest.NewServer(r)
	defer ts.Close()

	mockLists := []struct {
		Method   string
		URL      string
		Expected string
		Ctx      context.Context
		Body     io.Reader
	}{
		{
			Method:   "GET",
			URL:      "/v1/api/images",
			Expected: `"uploaded": ["https://i.imgur.com/8yc2oCz.jpg","https://i.imgur.com/u6GIFQA.jpg"]`,
			Body:     bytes.NewBufferString(``),
		},
		{
			Method:   "GET",
			URL:      "/v1/api/images/upload/551046a9-7448-466d-ad9e-4e870a8e4a59-20190107-130152",
			Expected: `{"id":"551046a9-7448-466d-ad9e-4e870a8e4a59-20190107-130152","created":"2019-01-07T13:01:52+08:00","finished":"2019-01-07T13:02:04+08:00","status":"complete","uploaded":{"complete":["https://i.imgur.com/MWfMA9e.jpg","https://i.imgur.com/thUkhoz.jpg"],"pending":null,"failed":null}}`,
			Body:     bytes.NewBufferString(``),
		},
		{
			Method:   "POST",
			URL:      "/v1/api/images/upload",
			Expected: `{"jobId":"bbeb5ebd-401c-4ea3-ae9c-b7de8ee3482b-20190107-110300"}`,
			Ctx:      context.WithValue(context.Background(), ctxKey{""}, ""),
			Body:     bytes.NewBufferString(`{"urls": ["https://farm3.staticflickr.com/2879/11234651086_681b3c2c00_b_d.jpg","https://farm4.staticflickr.com/3790/11244125445_3c2f32cd83_k_d.jpg"]}`),
		},
		{
			Method:   "GET",
			URL:      "/v1/api/credentials/236fc06f1bea341b6a1baac9d6d366cabef47c6e",
			Expected: `{"code": 202,"message": "Accepted"}`,
			Body:     bytes.NewBufferString(``),
		},
	}

	for _, rec := range mockLists {
		_, body := testRequest(t, ts, rec.Method, rec.URL, rec.Body)
		//t.Log(rec.URL, body)
		if body != rec.Expected {
			t.Fatalf("expected:%s got:%s", rec.Expected, body)
		}
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
