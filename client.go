package gobriva

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	// Base URLs for production and sandbox
	productionBaseURL = "https://partner.api.bri.co.id"
	sandboxBaseURL    = "https://sandbox.partner.api.bri.co.id"

	// Default timeout
	defaultTimeout = 30 * time.Second

	// Maximum number of bytes of body to include in logs (avoid huge logs)
	maxLogBodySize = 8 * 1024 // 8 KiB
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Authenticator interface for handling authentication
type Authenticator interface {
	Authenticate(ctx context.Context) error
	IsAuthenticated() bool
	EnsureAuthenticated(ctx context.Context) error
}

// Config holds the client configuration
type Config struct {
	PartnerID     string
	ClientID      string
	ClientSecret  string
	PrivateKey    string
	ChannelID     string
	IsSandbox     bool
	Timeout       time.Duration
	Debug         bool          // Enable debug logging for HTTP requests/responses
	Logger        *slog.Logger  // Optional: custom slog.Logger; if provided the client will use it (no global changes)
	HTTPClient    HTTPClient    // Optional: custom HTTP client for testing
	Authenticator Authenticator // Optional: custom authenticator for testing
}

// Client represents the BRI Virtual Account API client
type Client struct {
	httpClient   HTTPClient
	auth         Authenticator
	baseURL      string
	partnerID    string
	clientID     string
	clientSecret string
	privateKey   string
	channelID    string
	isSandbox    bool
	debug        bool
	logger       *slog.Logger
	accessToken  string
	tokenExpiry  time.Time
}

// NewClient creates a new BRI Virtual Account API client
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	// Use provided HTTP client or create default
	var httpClient HTTPClient
	if config.HTTPClient != nil {
		httpClient = config.HTTPClient
	} else {
		// Skip TLS verification for sandbox
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: config.IsSandbox},
		}
		httpClient = &http.Client{
			Transport: tr,
			Timeout:   config.Timeout,
		}
	}

	baseURL := productionBaseURL
	if config.IsSandbox {
		baseURL = sandboxBaseURL
	}

	client := &Client{
		httpClient:   httpClient,
		baseURL:      baseURL,
		partnerID:    config.PartnerID,
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		privateKey:   config.PrivateKey,
		channelID:    config.ChannelID,
		isSandbox:    config.IsSandbox,
		debug:        config.Debug,
	}

	// If a custom logger is provided, use it locally (do NOT change global slog.Default).
	// Otherwise, if Debug is enabled, create a local default logger so debug messages
	// are printed without affecting global application logger.
	if config.Logger != nil {
		client.logger = config.Logger
	} else if config.Debug {
		client.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	// Use provided authenticator or create default
	if config.Authenticator != nil {
		client.auth = config.Authenticator
	} else {
		client.auth = &DefaultAuthenticator{client: client}
	}

	return client
}

// DefaultAuthenticator implements the Authenticator interface
type DefaultAuthenticator struct {
	client *Client
}

// Authenticate performs OAuth2 authentication to get access token
func (a *DefaultAuthenticator) Authenticate(ctx context.Context) error {
	return a.client.authenticate(ctx)
}

// IsAuthenticated checks if the client has a valid access token
func (a *DefaultAuthenticator) IsAuthenticated() bool {
	return a.client.accessToken != "" && time.Now().Before(a.client.tokenExpiry)
}

// EnsureAuthenticated ensures the client has a valid access token
func (a *DefaultAuthenticator) EnsureAuthenticated(ctx context.Context) error {
	if !a.IsAuthenticated() {
		return a.Authenticate(ctx)
	}
	return nil
}

// generateExternalID generates a random 9-digit external ID
func (c *Client) generateExternalID() string {
	return fmt.Sprintf("%09d", rand.Intn(999999999))
}

// generateTimestamp generates current timestamp in ISO 8601 format
func (c *Client) generateTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")
}

// calculateSignature calculates HMAC-SHA512 signature for API requests
func (c *Client) calculateSignature(httpMethod, requestPath, requestBody string) (string, error) {
	// Parse request body if present
	var bodyStr string
	if httpMethod != "GET" && requestBody != "" {
		bodyStr = requestBody
	}

	// Create hash of request body using SHA256
	var payloadHash string
	if bodyStr != "" {
		payloadHash = fmt.Sprintf("%x", sha256.Sum256([]byte(bodyStr)))
	} else {
		payloadHash = fmt.Sprintf("%x", sha256.Sum256([]byte("")))
	}

	// Create signature payload
	timestamp := c.generateTimestamp()
	payload := fmt.Sprintf("%s:%s:%s:%s:%s",
		httpMethod, requestPath, c.accessToken, payloadHash, timestamp)

	// Calculate HMAC-SHA512
	h := hmac.New(sha512.New, []byte(c.clientSecret))
	h.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature, nil
}

// makeRequest makes an HTTP request with proper authentication
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Serialize body if present
	var bodyBytes []byte
	var bodyStr string
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
		bodyStr = string(bodyBytes)
	}

	// Calculate signature
	signature, err := c.calculateSignature(method, path, bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate signature: %w", err)
	}

	// Create request
	fullURL := c.baseURL + path
	var reqBytes io.Reader
	if bodyBytes != nil {
		reqBytes = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	timestamp := c.generateTimestamp()
	externalID := c.generateExternalID()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-PARTNER-ID", c.partnerID)
	req.Header.Set("X-EXTERNAL-ID", externalID)
	req.Header.Set("CHANNEL-ID", c.channelID)
	req.Header.Set("X-SIGNATURE", signature)
	req.Header.Set("X-TIMESTAMP", timestamp)

	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	// Debug logging - structured request (method/url/headers/body)
	if c.debug {
		// Prepare headers map copy
		headersMap := map[string][]string{}
		for k, v := range req.Header {
			headersMap[k] = append([]string(nil), v...)
		}

		// Prepare body for logging, truncate if too large
		var bodyForLog string
		if bodyBytes != nil {
			if len(bodyBytes) > maxLogBodySize {
				bodyForLog = string(bodyBytes[:maxLogBodySize]) + "... (truncated)"
			} else {
				bodyForLog = string(bodyBytes)
			}
		}

		if c.logger != nil {
			c.logger.Debug("HTTP Request",
				"method", req.Method,
				"url", req.URL.String(),
				"headers", headersMap,
				"body", bodyForLog,
			)
		} else {
			slog.Debug("HTTP Request",
				"method", req.Method,
				"url", req.URL.String(),
				"headers", headersMap,
				"body", bodyForLog,
			)
		}
	}

	// Make request with timing
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(start)
	if err != nil {
		return nil, err
	}

	// Debug logging - structured response (status/headers/body/duration)
	if c.debug {
		// Read and restore response body so caller can still read it
		respBodyBytes, _ := io.ReadAll(resp.Body)
		// replace the body so it can be read by the caller
		resp.Body = io.NopCloser(bytes.NewBuffer(respBodyBytes))

		// Prepare headers map copy
		respHeaders := map[string][]string{}
		for k, v := range resp.Header {
			respHeaders[k] = append([]string(nil), v...)
		}

		// Prepare body for logging, truncate if too large
		var respBodyForLog string
		if len(respBodyBytes) > 0 {
			if len(respBodyBytes) > maxLogBodySize {
				respBodyForLog = string(respBodyBytes[:maxLogBodySize]) + "... (truncated)"
			} else {
				respBodyForLog = string(respBodyBytes)
			}
		}

		if c.logger != nil {
			c.logger.Debug("HTTP Response",
				"status", resp.Status,
				"statusCode", resp.StatusCode,
				"headers", respHeaders,
				"body", respBodyForLog,
				"duration", duration.String(),
			)
		} else {
			slog.Debug("HTTP Response",
				"status", resp.Status,
				"statusCode", resp.StatusCode,
				"headers", respHeaders,
				"body", respBodyForLog,
				"duration", duration.String(),
			)
		}
	}

	return resp, nil
}

// parseErrorResponse parses an error response from the API
func (c *Client) parseErrorResponse(respBody []byte, httpStatusCode int) *StructuredBRIAPIResponse {
	var errorResp ErrorResponse
	json.Unmarshal(respBody, &errorResp)
	return &StructuredBRIAPIResponse{
		ResponseCode:    errorResp.ResponseCode,
		ResponseMessage: errorResp.ResponseMessage,
		HTTPStatusCode:  httpStatusCode,
		Timestamp:       time.Now(),
	}
}

// AuthResponse represents the OAuth2 token response
type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   string `json:"expiresIn"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}
