// Package http contains definition of the
// standard HTTP Gonyan Stream.
package http

import (
	"bytes"
	"fmt"
	"net/http"
)

// HTTPStream defines the standard Gonyan Stream for HTTP and HTTPS requests.
type HTTPStream struct {
	method      string
	url         string
	useHTTPS    bool
	prepareBody func([]byte) ([]byte, error)
}

// NewHTTPStream creates a new HTTP stream and sets its webhook URL.
func NewHTTPStream(url string) *HTTPStream {
	return &HTTPStream{
		method:      http.MethodGet,
		url:         url,
		useHTTPS:    false,
		prepareBody: nil,
	}
}

// SetMethod allows to define the HTTP method to be used by the stream.
// By default it'll perform a POST request but other methods are supported.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set`` functions can be chained together.
func (h *HTTPStream) SetMethod(method string) *HTTPStream {
	h.method = method
	return h
}

// SetCustomBodyPrepareFunction allows to define a custom body preparation
// function in order to apply custom operation on the body prior its
// transmission to the HTTP endpoint.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set`` functions can be chained together.
func (h *HTTPStream) SetCustomBodyPrepareFunction(f func([]byte) ([]byte, error)) *HTTPStream {
	h.prepareBody = f
	return h
}

// DisableHTTPS will set the internal flag for HTTPS to `false`
// thus disabling it; note that HTTPS is disabled by default.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set`` functions can be chained together.
func (h *HTTPStream) DisableHTTPS() *HTTPStream {
	h.useHTTPS = false
	return h
}

// EnableHTTPS will set the internal flag for HTTPS to `true` thus enabling it.
// HTTPS is disabled by default.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set`` functions can be chained together.
func (h *HTTPStream) EnableHTTPS() *HTTPStream {
	h.useHTTPS = true
	return h
}

// Write function defined to implement the Stream interface.
// The function prepares the body and fires the HTTP/HTTPS request
// using (optionally provided) headers and GET query parameters.
// Note: the actual HTTP request is performed inside a simple goroutine, it
// will be optimised in the future.
func (h *HTTPStream) Write(messageBytes []byte) (int, error) {
	body := messageBytes
	if h.prepareBody != nil {
		var err error
		body, err = h.prepareBody(messageBytes)
		if err != nil {
			return 0, fmt.Errorf("custom body prepare failed due to: %s", err.Error())
		}
	}

	go func() {
		// TODO: Handle request in a better way.
		if err := h.fireRequest(body); err != nil {
			fmt.Printf("[Gonyan] [HTTPStream] request firing failed due to: %s.\nRequest body: %+v", err.Error(), body)
			return
		}
	}()

	return len(body), nil
}

// fireRequest function will create and execute the actual HTTP request putting
// together all setup information, headers etc.
func (h *HTTPStream) fireRequest(body []byte) error {
	request, err := http.NewRequest(h.method, h.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("request creation failed due to: %s", err.Error())
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request execution failed due to: %s", err.Error())
	}
	defer response.Body.Close()

	return nil
}
