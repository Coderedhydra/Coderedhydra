// jules/core/http_client.go
package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar" // For session management
	"net/url"
	"strings"
	"time" // For default timeout
)

// BodyType defines the type of the HTTP request body.
type BodyType string

const (
	BodyTypeNone           BodyType = "none"
	BodyTypeFormURLEncoded BodyType = "form-urlencoded"
	BodyTypeJSON           BodyType = "json"
	BodyTypeXML            BodyType = "xml" // Placeholder
	BodyTypeRawString      BodyType = "raw_string"
	BodyTypeMultipart      BodyType = "multipart" // Placeholder for future
)

// CustomRequestDetails encapsulates all necessary information to build a complex HTTP request.
type CustomRequestDetails struct {
	URL             string
	Method          string // e.g., "GET", "POST", "PUT"
	Headers         map[string]string
	QueryParameters url.Values // Use url.Values for easier query string construction and multiple values per key.
	BodyType        BodyType
	FormData        url.Values        // For application/x-www-form-urlencoded
	JSONBody        interface{}       // For application/json; can be map[string]interface{} or a struct
	RawBody         []byte            // For other body types like text/plain, application/xml
	// TODO: Add field for multipart form data if needed
}

// InterceptionHook is a function type that can be used to modify requests/responses.
type InterceptionHook func(req *http.Request, res *http.Response) (*http.Request, *http.Response)

// JulesHTTPClient is a custom HTTP client with interception capabilities.
type JulesHTTPClient struct {
	Client http.Client
	Hooks  []InterceptionHook
}

// NewJulesHTTPClient creates a new instance of JulesHTTPClient.
// It now initializes with a cookie jar for session handling.
func NewJulesHTTPClient(timeoutSeconds int) (*JulesHTTPClient, error) {
	jar, err := cookiejar.New(nil) // nil uses default options
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	if timeoutSeconds <= 0 {
		timeoutSeconds = 15 // Default timeout
	}

	return &JulesHTTPClient{
		Client: http.Client{
			Jar:     jar,
			Timeout: time.Duration(timeoutSeconds) * time.Second,
			// TODO: Add proxy settings from config if needed for the client itself
			// TODO: Add redirect policy (e.g., check redirects but don't follow for scanner)
		},
	}, nil
}

// AddHook adds an interception hook to the client.
func (c *JulesHTTPClient) AddHook(hook InterceptionHook) {
	c.Hooks = append(c.Hooks, hook)
}

// SendCustomRequest constructs and sends an HTTP request based on provided details.
func (c *JulesHTTPClient) SendCustomRequest(details CustomRequestDetails) (*http.Response, []byte, error) {
	// 1. Construct URL with Query Parameters
	targetURL, err := url.Parse(details.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid target URL: %w", err)
	}
	if details.QueryParameters != nil {
		targetURL.RawQuery = details.QueryParameters.Encode()
	}

	// 2. Prepare Body
	var reqBody io.Reader
	var contentType string

	switch details.BodyType {
	case BodyTypeFormURLEncoded:
		if details.FormData != nil {
			reqBody = strings.NewReader(details.FormData.Encode())
			contentType = "application/x-www-form-urlencoded"
		}
	case BodyTypeJSON:
		if details.JSONBody != nil {
			jsonBytes, err := json.Marshal(details.JSONBody)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal JSON body: %w", err)
			}
			reqBody = bytes.NewReader(jsonBytes)
			contentType = "application/json; charset=utf-8"
		}
	case BodyTypeRawString:
		if details.RawBody != nil {
			reqBody = bytes.NewReader(details.RawBody)
			// contentType should be set in details.Headers for RawString, or a default like text/plain
			if ct, ok := details.Headers["Content-Type"]; ok {
				contentType = ct
			} else if ct, ok := details.Headers["content-type"]; ok {
                contentType = ct
            } else {
				contentType = "text/plain; charset=utf-8" // Default for raw string
			}
		}
	case BodyTypeXML:
		// Similar to RawString, expect Content-Type in headers
		if details.RawBody != nil {
			reqBody = bytes.NewReader(details.RawBody)
			if ct, ok := details.Headers["Content-Type"]; ok {
				contentType = ct
			} else if ct, ok := details.Headers["content-type"]; ok {
                contentType = ct
            } else {
				contentType = "application/xml; charset=utf-8"
			}
		}
	case BodyTypeNone:
		// No body
	default:
		// No body or unknown body type
	}

	// 3. Create Request
	method := strings.ToUpper(details.Method)
	if method == "" {
		method = "GET" // Default to GET
	}
	req, err := http.NewRequest(method, targetURL.String(), reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	// 4. Set Headers
	// Set Content-Type from body preparation first if it was determined
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	// Then apply custom headers, potentially overriding Content-Type if explicitly set
	for key, value := range details.Headers {
		req.Header.Set(key, value) // Use Set to override, Add to add multiple with same key
	}
	
	// Default User-Agent (can be overridden by details.Headers)
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Jules Security Toolkit/0.1")
	}


	// --- Placeholder for applying request hooks (might be less common for scanner requests) ---
	// var originalReqForHooks *http.Request = req
	// for _, hook := range c.Hooks {
	//  req, _ = hook(req, nil) // Hook modifies request
	// }
	// --- End Hook Placeholder ---

	fmt.Printf("CORE_HTTP_CLIENT: Sending %s request to %s\n", req.Method, req.URL.String())
	// if reqBody != nil && details.BodyType != BodyTypeFormURLEncoded { // Form data can be long
	//  // Careful with logging body, could be large or sensitive.
	//  // For debugging, one might log a snippet or if it's JSON/XML.
	// }


	// 5. Execute Request
	httpResp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("HTTP request execution failed: %w", err)
	}
	// defer httpResp.Body.Close() // Caller should close the body

	// 6. Read response body (optional here, but often useful for scanner)
	// It's often better for the caller to read and close the body.
	// However, for convenience in a scanner, reading it here and returning bytes can be useful.
	// This client will read and return the body, and also close it.
	var responseBodyBytes []byte
	if httpResp.Body != nil {
		responseBodyBytes, err = io.ReadAll(httpResp.Body)
		if err != nil {
			// Still return the response object, but with an error for body reading
			defer httpResp.Body.Close() // Ensure body is closed even on read error
			return httpResp, nil, fmt.Errorf("failed to read response body: %w", err)
		}
		defer httpResp.Body.Close() // Close after successful reading
	}
	
	// --- Placeholder for applying response hooks ---
	// for _, hook := range c.Hooks {
	//  _, httpResp = hook(originalReqForHooks, httpResp) // Hook modifies response
	// }
	// --- End Hook Placeholder ---

	return httpResp, responseBodyBytes, nil
}


// Do sends a standard HTTP request. Retained for general purpose or proxy use.
// The SendCustomRequest is more suitable for scanner's needs.
func (c *JulesHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	var originalReq *http.Request = req // Keep a copy for hooks if they modify req
	var res *http.Response

	// Apply request hooks (placeholder, actual modification needs careful handling of request body)
	// This is a simplified representation. Real hooks might need more context
	// or operate on request copies to avoid issues with read bodies.
	for _, hook := range c.Hooks {
		// Example: hook modifies originalReq in place or returns a new one
		// For this placeholder, assume it can modify req if it's a pointer, or originalReq
		// and then the modified version is used.
		// If hook returns a new request, it should be: req, _ = hook(req, nil)
		_, _ = hook(originalReq, nil) 
	}

	fmt.Printf("CORE_HTTP_CLIENT (Do): Sending request: %s %s\n", originalReq.Method, originalReq.URL.String())
	res, err = c.Client.Do(originalReq) // Use originalReq if hooks don't modify req in place
	if err != nil {
		return nil, err
	}

	// Apply response hooks (placeholder)
	for _, hook := range c.Hooks {
		// Example: hook modifies res in place or returns a new one
		_, res = hook(originalReq, res) 
	}

	return res, nil
}
