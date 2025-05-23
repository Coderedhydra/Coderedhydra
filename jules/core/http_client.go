// jules/core/http_client.go
package core

import (
	"net/http"
	"fmt"
)

// InterceptionHook is a function type that can be used to modify requests/responses.
type InterceptionHook func(req *http.Request, res *http.Response) (*http.Request, *http.Response)

// JulesHTTPClient is a custom HTTP client with interception capabilities.
type JulesHTTPClient struct {
	Client http.Client
	Hooks  []InterceptionHook
}

// NewJulesHTTPClient creates a new instance of JulesHTTPClient.
func NewJulesHTTPClient() *JulesHTTPClient {
	return &JulesHTTPClient{
		Client: http.Client{},
		// TODO: Add default timeout, proxy settings from config
	}
}

// AddHook adds an interception hook to the client.
func (c *JulesHTTPClient) AddHook(hook InterceptionHook) {
	c.Hooks = append(c.Hooks, hook)
}

// Do sends an HTTP request and allows hooks to modify it and its response.
func (c *JulesHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	var originalReq *http.Request = req
	var res *http.Response

	// Apply request hooks (placeholder, actual modification needs careful handling of request body)
	for _, hook := range c.Hooks {
		// This is a simplified representation. Real hooks might need more context
		// or operate on request copies to avoid issues with read bodies.
		_, _ = hook(originalReq, nil) // Example: hook modifies originalReq in place
	}

	fmt.Printf("Sending request: %s %s\n", originalReq.Method, originalReq.URL.String())
	res, err = c.Client.Do(originalReq)
	if err != nil {
		return nil, err
	}

	// Apply response hooks (placeholder)
	for _, hook := range c.Hooks {
		_, res = hook(originalReq, res) // Example: hook modifies res in place
	}

	return res, nil
}
