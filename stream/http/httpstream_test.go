package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

// TestConstructorURL verifies that the URL is properly copied.
func TestConstructorURL(t *testing.T) {
	s := NewStream("the-url.com")
	if s.url != "the-url.com" {
		t.Fatalf("Unexpected url found. Expected: `%s` - Found: `%s`.", "the-url.com", s.url)
	}
}

// TestSetMethod verifies that the SetMethod function
// works as expected.
func TestSetMethod(t *testing.T) {
	s := NewStream("")

	if s.method != http.MethodPost {
		t.Fatalf("Unexpected default method found: `%s`.", s.method)
	}

	s.SetMethod(http.MethodPut)
	if s.method != http.MethodPut {
		t.Fatalf("Unexpected request method found. Expected: `%s` - Found: `%s`.", http.MethodPut, s.method)
	}
}

// TestSetCustomBodyPrepareFunction verifies that the
// SetCustomBodyPrepareFunction function works as expected.
func TestSetCustomBodyPrepareFunction(t *testing.T) {
	s := NewStream("")
	if s.prepareBody != nil {
		t.Fatalf("Unexpected body prepare function found: %+v", reflect.ValueOf(s.prepareBody))
	}

	prepare := func(in []byte) ([]byte, error) {
		return in, nil
	}

	s.SetCustomBodyPrepareFunction(prepare)
	if reflect.ValueOf(s.prepareBody) != reflect.ValueOf(prepare) {
		t.Fatalf("Unexpected request method found. Expected: `%s` - Found: `%s`.", http.MethodPut, s.method)
	}
}

// TestEnableAndDisableHTTPS will verify default value for
// HTTPS usage, then verifies that the EnableHTTPS and
// DisableHTTPS functions properly work.
func TestEnableAndDisableHTTPS(t *testing.T) {
	s := NewStream("")
	if s.useHTTPS == true {
		t.Fatalf("HTTPS should be disabled by default!")
	}

	s.EnableHTTPS()
	if s.useHTTPS != true {
		t.Fatalf("HTTPS should have been enabled!")
	}

	s.DisableHTTPS()
	if s.useHTTPS != false {
		t.Fatalf("HTTPS should have been disabled!")
	}
}

// TestSetAllHeaders verifies that SetAllHeaders function works as expected.
func TestSetAllHeaders(t *testing.T) {
	s := NewStream("")

	headers := map[string]string{
		"hey":  "ho",
		"lets": "go",
	}
	s.SetAllHeaders(headers)
	if !reflect.DeepEqual(headers, s.headers) {
		t.Fatalf("Unexpected headers map. Expected: %+v - Found: %+v", headers, s.headers)
	}
}

// TestSetHeader verifies that SetHeader function works as expected.
func TestSetHeader(t *testing.T) {
	s := NewStream("")

	s.SetHeader("hey", "oh")
	val, ok := s.headers["hey"]
	if !ok {
		t.Fatalf("The key-pair should exist!")
	}
	if val != "oh" {
		t.Fatalf("Unexpected value found. Expected: `%s` - Found: `%s`.", "oh", val)
	}
}

// TestRemoveHeader verifies that RemoveHeader function works as expected.
func TestRemoveHeader(t *testing.T) {
	s := NewStream("")

	headers := map[string]string{
		"hey":  "ho",
		"lets": "go",
	}
	s.SetAllHeaders(headers)
	if !reflect.DeepEqual(headers, s.headers) {
		t.Fatalf("Unexpected headers map. Expected: %+v - Found: %+v", headers, s.headers)
	}

	s.RemoveHeader("NOT PRESENT")

	s.RemoveHeader("hey")
	val, ok := s.headers["lets"]
	if !ok {
		t.Fatalf("The key-pair should exist!")
	}
	if val != "go" {
		t.Fatalf("Unexpected value found. Expected: `%s` - Found: `%s`.", "go", val)
	}
}

// TestSetAllQueryParams verifies that SetAllQueryParams function works as expected.
func TestSetAllQueryParams(t *testing.T) {
	s := NewStream("")

	params := map[string]string{
		"hey":  "ho",
		"lets": "go",
	}
	s.SetAllQueryParams(params)
	if !reflect.DeepEqual(params, s.queryParams) {
		t.Fatalf("Unexpected params map. Expected: %+v - Found: %+v", params, s.queryParams)
	}
}

// TestSetQueryParams verifies that SetQueryParams function works as expected.
func TestSetQueryParams(t *testing.T) {
	s := NewStream("")

	s.SetQueryParam("hey", "oh")
	val, ok := s.queryParams["hey"]
	if !ok {
		t.Fatalf("The key-pair should exist!")
	}
	if val != "oh" {
		t.Fatalf("Unexpected value found. Expected: `%s` - Found: `%s`.", "oh", val)
	}
}

// TestRemoveQueryParam verifies that RemoveQueryParam function works as expected.
func TestRemoveQueryParam(t *testing.T) {
	s := NewStream("")

	params := map[string]string{
		"hey":  "ho",
		"lets": "go",
	}
	s.SetAllQueryParams(params)
	if !reflect.DeepEqual(params, s.queryParams) {
		t.Fatalf("Unexpected params map. Expected: %+v - Found: %+v", params, s.queryParams)
	}

	s.RemoveQueryParam("NOT PRESENT")

	s.RemoveQueryParam("hey")
	val, ok := s.queryParams["lets"]
	if !ok {
		t.Fatalf("The key-pair should exist!")
	}
	if val != "go" {
		t.Fatalf("Unexpected value found. Expected: `%s` - Found: `%s`.", "go", val)
	}
}

// TestFireRequest successfully fires an HTTP PUT request.
func TestFireRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPut {
			t.Fatalf("Unexpected request method. Expected: %s - Found: %s.", http.MethodPut, r.Method)
		}

		// Read body.
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading body: %v", err)
			return
		}
		defer r.Body.Close()

		if len(payload) != 3 {
			t.Fatalf("Unexpected number of bytes read. Expected: %d - Found: %d.", 3, len(payload))
		}
		if bytes.Compare(payload, []byte("hey")) != 0 {
			t.Fatalf("Unexpected payload bytes found. Expected: %+v - Found: %+v.", []byte("hey"), payload)
		}
		w.Write(nil)
	}))
	defer ts.Close()

	s := NewStream(ts.URL)
	s.DisableHTTPS()
	s.SetMethod(http.MethodPut)
	if err := s.fireRequest([]byte("hey")); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

// TestFireRequestVerifyHeadersWithSet successfully fires an HTTP PUT request
// and verify that headers have been set using SetHeader.
func TestFireRequestVerifyHeadersWithSet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			t.Fatalf("Unexpected request method. Expected: %s - Found: %s.", http.MethodPost, r.Method)
		}

		if val := r.Header.Get("hey"); val != "oh" {
			t.Fatalf("Unexpected request header value for key `%s`. Expected: `%s` - Found: `%s`.", "hey", "oh", val)
		}

		if val := r.Header.Get("lets"); val != "go" {
			t.Fatalf("Unexpected request header value for key `%s`. Expected: `%s` - Found: `%s`.", "lets", "go", val)
		}

		w.Write(nil)
	}))
	defer ts.Close()

	s := NewStream(ts.URL)
	s.DisableHTTPS()

	s.SetHeader("hey", "oh").SetHeader("lets", "go")
	if err := s.fireRequest([]byte{}); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

// TestFireRequestVerifyHeadersWithSetAll successfully fires an HTTP PUT request
// and verify that headers have been set using SetAllHeaders.
func TestFireRequestVerifyHeadersWithSetAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			t.Fatalf("Unexpected request method. Expected: %s - Found: %s.", http.MethodPost, r.Method)
		}

		if val := r.Header.Get("hey"); val != "oh" {
			t.Fatalf("Unexpected request header value for key `%s`. Expected: `%s` - Found: `%s`.", "hey", "oh", val)
		}

		if val := r.Header.Get("lets"); val != "go" {
			t.Fatalf("Unexpected request header value for key `%s`. Expected: `%s` - Found: `%s`.", "lets", "go", val)
		}

		if val := r.Header.Get("what"); val != "" {
			t.Fatalf("Header `%s` should have been removed. Fond value: `%s`.", "`what`", val)
		}

		w.Write(nil)
	}))
	defer ts.Close()

	s := NewStream(ts.URL)
	s.DisableHTTPS()

	s.SetAllHeaders(map[string]string{
		"hey":  "oh",
		"lets": "go",
		"what": "not",
	})
	s.RemoveHeader("what")
	if err := s.fireRequest([]byte{}); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

// TestFireRequestVerifyQueryParamsWithSet successfully fires an HTTP GET
// request and verify that GET query parameters have been
// set using SetQueryParam.
func TestFireRequestVerifyQueryParamsWithSet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			t.Fatalf("Unexpected request method. Expected: %s - Found: %s.", http.MethodGet, r.Method)
		}

		queryParams := r.URL.Query()

		if val := queryParams.Get("hey"); val != "oh" {
			t.Fatalf("Unexpected request query param value for key `%s`. Expected: `%s` - Found: `%s`.", "hey", "oh", val)
		}

		if val := queryParams.Get("lets"); val != "go" {
			t.Fatalf("Unexpected request query param value for key `%s`. Expected: `%s` - Found: `%s`.", "lets", "go", val)
		}

		w.Write(nil)
	}))
	defer ts.Close()

	s := NewStream(ts.URL)
	s.DisableHTTPS()
	s.SetMethod(http.MethodGet)

	s.SetQueryParam("hey", "oh").SetQueryParam("lets", "go")
	if err := s.fireRequest([]byte{}); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

// TestFireRequestVerifyQueryParamsWithSetAll successfully fires an HTTP GET
// request and verify that GET query parameters have been
// set using SetAllQueryParams.
func TestFireRequestVerifyQueryParamsWithSetAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			t.Fatalf("Unexpected request method. Expected: %s - Found: %s.", http.MethodGet, r.Method)
		}

		queryParams := r.URL.Query()
		if val := queryParams.Get("hey"); val != "oh" {
			t.Fatalf("Unexpected request query param value for key `%s`. Expected: `%s` - Found: `%s`.", "hey", "oh", val)
		}

		if val := queryParams.Get("lets"); val != "go" {
			t.Fatalf("Unexpected request query param value for key `%s`. Expected: `%s` - Found: `%s`.", "lets", "go", val)
		}

		if val := queryParams.Get("what"); val != "" {
			t.Fatalf("Query param `%s` should have been removed. Fond value: `%s`.", "`what`", val)
		}

		w.Write(nil)
	}))
	defer ts.Close()

	s := NewStream(ts.URL)
	s.DisableHTTPS()
	s.SetMethod(http.MethodGet)

	s.SetAllQueryParams(map[string]string{
		"hey":  "oh",
		"lets": "go",
		"what": "not",
	})
	s.RemoveQueryParam("what")
	if err := s.fireRequest([]byte{}); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}

// TestFireRequestFailure covers the possible error generated by
// fireRequest function.
func TestFireRequestFailure(t *testing.T) {
	s := NewStream("invalid-url.com")
	if err := s.fireRequest([]byte("hey")); err == nil {
		t.Fatalf("This request should have failed, found nil error instead.")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(nil)
	}))
	defer ts.Close()

	s = NewStream(ts.URL)
	s.SetMethod("NOT A METHOD")
	if err := s.fireRequest([]byte("hey")); err == nil {
		t.Fatalf("This request should have failed, found nil error instead.")
	}
}

// TestWriteSuccess successfully uses Write stream entrypoint.
func TestWriteSuccess(t *testing.T) {
	timedout := true
	mtx := &sync.Mutex{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPut {
			w.Write(nil)
			t.Fatalf("Unexpected request method. Expected: %s - Found: %s.", http.MethodPut, r.Method)
		}

		// Read body.
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write(nil)
			t.Fatalf("Error reading body: %v", err)
			return
		}
		defer r.Body.Close()

		if len(payload) != 6 {
			w.Write(nil)
			t.Fatalf("Unexpected number of bytes read. Expected: %d - Found: %d.", 6, len(payload))
		}
		if bytes.Compare(payload, []byte("hey ho")) != 0 {
			w.Write(nil)
			t.Fatalf("Unexpected payload bytes found. Expected: %+v - Found: %+v.", []byte("hey ho"), payload)
		}
		mtx.Lock()
		timedout = false
		mtx.Unlock()
		w.Write(nil)
	}))
	defer ts.Close()

	s := NewStream(ts.URL)
	s.DisableHTTPS()
	s.SetMethod(http.MethodPut)

	s.SetCustomBodyPrepareFunction(func(in []byte) ([]byte, error) {
		stringed := string(in)
		return []byte(stringed + " ho"), nil
	})

	nbytes, err := s.Write([]byte("hey"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if nbytes != 6 {
		t.Fatalf("Unexpected number of written bytes. Expected: %d - Found: %d.", 6, nbytes)
	}

	// Sleep a bit to let the handler process the request.
	time.Sleep(5 * time.Second)

	mtx.Lock()
	willfail := timedout
	mtx.Unlock()
	if willfail {
		t.Fatalf("No successful request was received by the test handler and it timed out.")
	}
}

// TestWritePrepareBodyFails verifies what happens when the provided
// body prepare function fails.
func TestWritePrepareBodyFails(t *testing.T) {
	s := NewStream("unneded.url")
	s.DisableHTTPS()
	s.SetMethod(http.MethodPut)

	s.SetCustomBodyPrepareFunction(func(in []byte) ([]byte, error) {
		return nil, fmt.Errorf("an-error")
	})

	nbytes, err := s.Write([]byte("hey"))
	if err == nil {
		t.Fatalf("An error was expected!")
	}
	if err.Error() != "custom body prepare failed due to: an-error" {
		t.Fatalf("Unexpected error value found. Expected: `%s` - Found: `%s`.", "custom body prepare failed due to: an-error", err.Error())
	}
	if nbytes != 0 {
		t.Fatalf("Unexpected number of written bytes. Expected: %d - Found: %d.", 0, nbytes)
	}
}
