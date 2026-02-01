package utils

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

// Default timeout for HTTP requests
const defaultTimeout = 30 * time.Second

// HTTPClient wraps fasthttp.Client with convenient methods
type HTTPClient struct {
	client  *fasthttp.Client
	baseURL string
	headers map[string]string
}

// NewHTTPClient creates a new HTTPClient with optional base URL
func NewHTTPClient(baseURL ...string) *HTTPClient {
	c := &HTTPClient{
		client: &fasthttp.Client{
			ReadTimeout:  defaultTimeout,
			WriteTimeout: defaultTimeout,
		},
		headers: make(map[string]string),
	}
	if len(baseURL) > 0 {
		c.baseURL = baseURL[0]
	}
	return c
}

// SetTimeout sets read and write timeout
func (c *HTTPClient) SetTimeout(timeout time.Duration) *HTTPClient {
	c.client.ReadTimeout = timeout
	c.client.WriteTimeout = timeout
	return c
}

// SetHeader sets a default header for all requests
func (c *HTTPClient) SetHeader(key, value string) *HTTPClient {
	c.headers[key] = value
	return c
}

// SetBearerToken sets Authorization header with Bearer token
func (c *HTTPClient) SetBearerToken(token string) *HTTPClient {
	c.headers["Authorization"] = "Bearer " + token
	return c
}

// HTTPResponse wraps the response data
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    *fasthttp.ResponseHeader
}

// buildURL constructs full URL from base + path
func (c *HTTPClient) buildURL(url string) string {
	if c.baseURL != "" && len(url) > 0 && url[0] == '/' {
		return c.baseURL + url
	}
	return url
}

// applyHeaders applies default headers to request
func (c *HTTPClient) applyHeaders(req *fasthttp.Request) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}

// do executes an HTTP request
func (c *HTTPClient) do(method, url string, body []byte, contentType string) (*HTTPResponse, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.buildURL(url))
	req.Header.SetMethod(method)
	c.applyHeaders(req)

	if body != nil {
		req.SetBody(body)
		if contentType != "" {
			req.Header.SetContentType(contentType)
		}
	}

	if err := c.client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("http %s %s: %w", method, url, err)
	}

	// Copy body since resp will be released
	bodyBytes := make([]byte, len(resp.Body()))
	copy(bodyBytes, resp.Body())

	return &HTTPResponse{
		StatusCode: resp.StatusCode(),
		Body:       bodyBytes,
	}, nil
}

// Get sends a GET request
func (c *HTTPClient) Get(url string) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodGet, url, nil, "")
}

// Post sends a POST request with JSON body
func (c *HTTPClient) Post(url string, body []byte) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodPost, url, body, "application/json")
}

// PostForm sends a POST request with form data
func (c *HTTPClient) PostForm(url string, body []byte) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodPost, url, body, "application/x-www-form-urlencoded")
}

// Put sends a PUT request with JSON body
func (c *HTTPClient) Put(url string, body []byte) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodPut, url, body, "application/json")
}

// Patch sends a PATCH request with JSON body
func (c *HTTPClient) Patch(url string, body []byte) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodPatch, url, body, "application/json")
}

// Delete sends a DELETE request
func (c *HTTPClient) Delete(url string) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodDelete, url, nil, "")
}

// DeleteWithBody sends a DELETE request with JSON body
func (c *HTTPClient) DeleteWithBody(url string, body []byte) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodDelete, url, body, "application/json")
}

// Head sends a HEAD request
func (c *HTTPClient) Head(url string) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodHead, url, nil, "")
}

// Options sends an OPTIONS request
func (c *HTTPClient) Options(url string) (*HTTPResponse, error) {
	return c.do(fasthttp.MethodOptions, url, nil, "")
}

// --- Standalone functions for quick use ---

// Get sends a simple GET request
func Get(url string) (*HTTPResponse, error) {
	return NewHTTPClient().Get(url)
}

// Post sends a simple POST request with JSON body
func Post(url string, body []byte) (*HTTPResponse, error) {
	return NewHTTPClient().Post(url, body)
}

// Put sends a simple PUT request with JSON body
func Put(url string, body []byte) (*HTTPResponse, error) {
	return NewHTTPClient().Put(url, body)
}

// Patch sends a simple PATCH request with JSON body
func Patch(url string, body []byte) (*HTTPResponse, error) {
	return NewHTTPClient().Patch(url, body)
}

// Delete sends a simple DELETE request
func Delete(url string) (*HTTPResponse, error) {
	return NewHTTPClient().Delete(url)
}
