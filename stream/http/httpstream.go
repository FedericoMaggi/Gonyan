// Package http contains definition of the
// standard HTTP Gonyan Stream.
package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// Stream defines the standard Gonyan Stream for HTTP and HTTPS requests.
type Stream struct {
	method      string                       // HTTP used method;
	url         string                       // The URL webhook;
	useHTTPS    bool                         // Flag to activate TLS/SSL;
	prepareBody func([]byte) ([]byte, error) // Function executed on body before transmission;
	headers     map[string]string            // HTTP headers container;
	queryParams map[string]string            // GET query parameter container.
}

// NewStream creates a new HTTP stream and sets its webhook URL.
func NewStream(url string) *Stream {
	return &Stream{
		method:      http.MethodPost,
		url:         url,
		useHTTPS:    false,
		prepareBody: nil,
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
	}
}

// SetMethod allows to define the HTTP method to be used by the stream.
// By default it'll perform a POST request but other methods are supported.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) SetMethod(method string) *Stream {
	h.method = method
	return h
}

// SetCustomBodyPrepareFunction allows to define a custom body preparation
// function in order to apply custom operation on the body prior its
// transmission to the HTTP endpoint.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) SetCustomBodyPrepareFunction(f func([]byte) ([]byte, error)) *Stream {
	h.prepareBody = f
	return h
}

// DisableHTTPS will set the internal flag for HTTPS to `false`
// thus disabling it; note that HTTPS is disabled by default.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) DisableHTTPS() *Stream {
	h.useHTTPS = false
	return h
}

// EnableHTTPS will set the internal flag for HTTPS to `true` thus enabling it.
// HTTPS is disabled by default.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) EnableHTTPS() *Stream {
	h.useHTTPS = true
	return h
}

// SetAllHeaders sets provided key-value pairs for later usage as HTTP headers.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) SetAllHeaders(header map[string]string) *Stream {
	h.headers = header
	return h
}

// SetHeader sets provided key-value pair for later usage as HTTP headers.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) SetHeader(key, value string) *Stream {
	h.headers[key] = value
	return h
}

// RemoveHeader removes provided key-value pair from HTTP headers.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) RemoveHeader(key string) {
	delete(h.headers, key)
}

// SetAllQueryParams sets provided key-value pairs for later usage
// as HTTP Query parameters.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) SetAllQueryParams(queryParams map[string]string) *Stream {
	h.queryParams = queryParams
	return h
}

// SetQueryParam sets provided key-value pair for later usage
// as HTTP Query parameters.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) SetQueryParam(key, value string) *Stream {
	h.queryParams[key] = value
	return h
}

// RemoveQueryParam removes provided key-value pair from HTTP Query parameters.
// Note: the method will return the same instance of the invoked structure
// so that multiple `Set` functions can be chained together.
func (h *Stream) RemoveQueryParam(key string) {
	delete(h.queryParams, key)
}

// Write function defined to implement the Stream interface.
// The function prepares the body and fires the HTTP/HTTPS request
// using (optionally provided) headers and GET query parameters.
// Note: the actual HTTP request is performed inside a simple goroutine, it
// will be optimised in the future.
func (h *Stream) Write(messageBytes []byte) (int, error) {
	body := messageBytes
	if h.prepareBody != nil {
		var err error
		body, err = h.prepareBody(messageBytes)
		if err != nil {
			return 0, fmt.Errorf("custom body prepare failed due to: %s", err.Error())
		}
	}

	go func(body []byte) {
		// TODO: Handle request in a better way.
		if err := h.fireRequest(body); err != nil {
			fmt.Printf("[Gonyan] [Stream] request firing failed due to: %s.\nRequest body: %+v", err.Error(), body)
			return
		}
	}(body)

	return len(body), nil
}

// fireRequest function will create and execute the actual HTTP request putting
// together all setup information, headers etc.
// The expected input is the previously prepared body (if a prepare function is
// provided).
func (h *Stream) fireRequest(preparedBody []byte) error {
	targetURL := h.url

	// Prepare GET Query parameters for GET method.
	if h.method == http.MethodGet {
		getParams := url.Values{}
		for key, val := range h.queryParams {
			getParams.Add(key, val)
		}
		targetURL += "?" + getParams.Encode()
	}

	request, err := http.NewRequest(h.method, targetURL, bytes.NewBuffer(preparedBody))
	if err != nil {
		return fmt.Errorf("request creation failed due to: %s", err.Error())
	}

	// Add all headers to request.
	for key, val := range h.headers {
		request.Header.Add(key, val)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request execution failed due to: %s", err.Error())
	}
	defer response.Body.Close()

	return nil
}
