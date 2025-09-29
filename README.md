# Go BRI Virtual Account Adapter

A Go library for integrating with BRI's BRIVA WS SNAP BI API.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Features

- **Complete BRIVA WS API implementation** - Full coverage of BRI Virtual Account operations
- **OAuth2 authentication with RSA signing** - Secure authentication using asymmetric cryptography
- **HMAC-SHA512 request signing** - Cryptographic request integrity verification
- **Production and sandbox environments** - Environment-specific configuration support
- **Interface-based design for testability** - Dependency injection for comprehensive testing
- **Comprehensive test suite** - 75+ unit tests with 81% code coverage
- **Debug mode** for HTTP request/response logging with timing measurements
- **Enhanced error handling** with detailed suggestions and severity levels
- **Comprehensive response code definitions** - Built-in BRI API response code handling

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Security](#security)
- [Development](#development)
- [Changelog](#changelog)
- [License](#license)

## Installation

```bash
go get gobriva
```

### Requirements

- Go 1.21 or later
- Valid BRI API credentials (Partner ID, Client ID, Client Secret, Private Key)
- Channel ID for your application

## Quick Start

```go
package main

import (
	"context"
	"log"

	"gobriva"
)

func main() {
	client := gobriva.NewClient(gobriva.Config{
		PartnerID:    "your-partner-id",
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		PrivateKey:   "your-private-key-pem",
		ChannelID:    "your-channel-id",
		IsSandbox:    true,
		Debug:        true,
	})

	// Create a virtual account
	req := &gobriva.CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "John Doe",
		TotalAmount: gobriva.Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "trx-123456",
	}

	resp, err := client.CreateVirtualAccount(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created VA: %s", resp.VirtualAccountData.VirtualAccountNo)
}
```

## Architecture

### Design Principles

The library follows SOLID principles and Go best practices:

- **Interface-based design** - All dependencies use interfaces for easy testing and mocking
- **Dependency injection** - HTTP client and authenticator can be injected for testing
- **Structured error handling** - Custom error types with detailed categorization
- **Context support** - All operations support context for cancellation and timeouts

### Core Components

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Application   │────│     Client       │────│   HTTP Client   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │  Authenticator   │
                       └──────────────────┘
```

- **Client**: Main API client with all BRIVA operations
- **Authenticator**: Handles OAuth2 authentication and token management
- **HTTPClient**: Interface for making HTTP requests (default: Go's net/http)
- **Response Codes**: Structured error handling with BRI-specific response codes

### Security Architecture

- **RSA-SHA256** for OAuth2 authentication
- **HMAC-SHA512** for request signing
- **Timestamp-based signatures** to prevent replay attacks
- **External ID generation** for request deduplication

## Configuration

### Client Configuration

```go
type Config struct {
	PartnerID     string
	ClientID      string
	ClientSecret  string
	PrivateKey    string
	ChannelID     string
	IsSandbox     bool
	Timeout       time.Duration
	Debug         bool          // Enable debug logging for HTTP requests/responses
	HTTPClient    HTTPClient    // Optional: custom HTTP client for testing
	Authenticator Authenticator // Optional: custom authenticator for testing
}
```

### Environment Variables

```bash
export BRI_PARTNER_ID="your-partner-id"
export BRI_CLIENT_ID="your-client-id"
export BRI_CLIENT_SECRET="your-client-secret"
export BRI_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\n..."
export BRI_CHANNEL_ID="your-channel-id"
export BRI_SANDBOX="true"
```

### Private Key Format

The private key should be in PEM format:

```pem
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC...
-----END PRIVATE KEY-----
```

## API Reference

### Virtual Account Operations

#### CreateVirtualAccount

Creates a new virtual account.

```go
func (c *Client) CreateVirtualAccount(ctx context.Context, req *CreateVirtualAccountRequest) (*CreateVirtualAccountResponse, error)
```

**Request:**

```go
type CreateVirtualAccountRequest struct {
	PartnerServiceID   string         `json:"partnerServiceId"`
	CustomerNo         string         `json:"customerNo"`
	VirtualAccountNo   string         `json:"virtualAccountNo"`
	VirtualAccountName string         `json:"virtualAccountName"`
	TotalAmount        Amount         `json:"totalAmount"`
	ExpiredDate        string         `json:"expiredDate"`
	TrxID              string         `json:"trxId"`
	AdditionalInfo     AdditionalInfo `json:"additionalInfo,omitempty"`
}
```

#### UpdateVirtualAccount

Updates an existing virtual account.

```go
func (c *Client) UpdateVirtualAccount(ctx context.Context, req *UpdateVirtualAccountRequest) (*UpdateVirtualAccountResponse, error)
```

#### UpdateVirtualAccountStatus

Updates the status of a virtual account.

```go
func (c *Client) UpdateVirtualAccountStatus(ctx context.Context, req *UpdateVirtualAccountStatusRequest) (*UpdateVirtualAccountStatusResponse, error)
```

#### InquiryVirtualAccount

Retrieves information about a virtual account.

```go
func (c *Client) InquiryVirtualAccount(ctx context.Context, req *InquiryVirtualAccountRequest) (*InquiryVirtualAccountResponse, error)
```

#### InquiryVirtualAccountStatus

Retrieves the status of a virtual account.

```go
func (c *Client) InquiryVirtualAccountStatus(ctx context.Context, req *InquiryVirtualAccountStatusRequest) (*InquiryVirtualAccountStatusResponse, error)
```

#### DeleteVirtualAccount

Deletes a virtual account.

```go
func (c *Client) DeleteVirtualAccount(ctx context.Context, req *DeleteVirtualAccountRequest) (*DeleteVirtualAccountResponse, error)
```

#### GetVirtualAccountReport

Retrieves transaction reports for virtual accounts.

```go
func (c *Client) GetVirtualAccountReport(ctx context.Context, req *VirtualAccountReportRequest) (*VirtualAccountReportResponse, error)
```

### Common Types

#### Amount

```go
type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}
```

#### AdditionalInfo

```go
type AdditionalInfo struct {
	Description string `json:"description,omitempty"`
}
```

## Error Handling

The library provides comprehensive error handling with structured error types.

### Error Types

#### StructuredBRIAPIResponse

```go
type StructuredBRIAPIResponse struct {
	HTTPStatusCode    int
	ResponseCode      string
	ResponseMessage   string
	ResponseDefinition *BRIVAResponseDefinition
}
```

### Error Categories

- `Success` - Successful operations (2xx)
- `BadRequest` - Invalid request data (4xx)
- `Unauthorized` - Authentication failures (401)
- `Forbidden` - Permission issues (403)
- `NotFound` - Resource not found (404)
- `MethodNotAllowed` - HTTP method not allowed (405)
- `Conflict` - Resource conflicts (409)
- `InternalServerError` - Server errors (5xx)
- `BadGateway` - Gateway errors (502)
- `ServiceUnavailable` - Service unavailable (503)
- `Pending` - Unknown response codes requiring manual verification

### Error Methods

```go
func (e *StructuredBRIAPIResponse) GetCategory() ErrorCategory
func (e *StructuredBRIAPIResponse) GetTimestamp() time.Time
func (e *StructuredBRIAPIResponse) Error() string 
func (e *StructuredBRIAPIResponse) IsPending() bool
func (e *StructuredBRIAPIResponse) IsSuccess() bool
func (e *StructuredBRIAPIResponse) IsClientError() bool
```

### Error Handling Example

```go
resp, err := client.CreateVirtualAccount(context.Background(), req)
if err != nil {
	if briErr, ok := err.(*gobriva.StructuredBRIAPIResponse); ok {
		log.Printf("BRI Error [%s]: %s", briErr.ResponseCode, briErr.ResponseMessage)
		log.Printf("Category: %s", briErr.GetCategory())

		if briErr.IsPending() {
			log.Printf("Unknown response code - manual verification required")
		}
	}
	return err
}
```

## Testing

### Test Coverage

Current test coverage: **81.0%**

### Testing Strategy

The library uses multiple testing approaches:

1. **Unit Tests** - Test individual functions and methods
2. **Integration Tests** - Test with real HTTP servers using httptest
3. **Mock Tests** - Test with injected mock dependencies
4. **Negative Tests** - Test error scenarios and edge cases

### Test Bed Pattern

The library uses a test bed pattern for integration testing:

```go
func TestCreateVirtualAccount(t *testing.T) {
	tb := setupTestBed(t, func (w http.ResponseWriter, r *http.Request) {
		// Mock BRI API responses
		w.WriteHeader(200)
		w.Write([]byte(`{"responseCode":"2002700","responseMessage":"Success"}`))
	})
	defer tb.teardown()

	client := tb.client
	// Test implementation
}
```

### Running Tests

```bash
# Run all tests
go test -v

# Run with coverage
go test -cover

# Run specific test
go test -run TestCreateVirtualAccount

# Run with race detection
go test -race
```

### Mock Testing

```go
type MockHTTPClient struct {
	DoFunc func (*http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
```

## Security

### Authentication Security

- **RSA-SHA256** signatures for OAuth2 token requests
- **HMAC-SHA512** for API request integrity
- **Timestamp-based signatures** with 5-minute validity windows
- **Unique external IDs** for request deduplication

### Request Signing Process

1. Generate timestamp and external ID
2. Create signature string: `HTTP_METHOD + ":" + REQUEST_PATH + ":" + REQUEST_BODY + ":" + TIMESTAMP`
3. Sign with HMAC-SHA512 using client secret
4. Include signature in `X-SIGNATURE` header

### Best Practices

- Store private keys securely (never in source code)
- Use environment variables for sensitive configuration
- Rotate credentials regularly
- Enable debug mode only in development
- Validate all input data before API calls

## Performance

### HTTP Client Configuration

The library uses Go's default HTTP client. For production use, consider:

```go
import "net/http"

transport := &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     90 * time.Second,
}

client := &http.Client{
	Transport: transport,
	Timeout:   30 * time.Second,
}
```

### Context and Timeouts

All operations support context for timeout control:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.CreateVirtualAccount(ctx, req)
```

### Debug Mode Performance Impact

Debug mode adds logging overhead. Disable in production:

```go
client := gobriva.NewClient(gobriva.Config{
	// ... other config
	Debug: false, // Disable for production
})
```

## Development

### Project Structure

```
gobriva/
├── client.go          # Main client implementation
├── auth.go            # Authentication logic
├── va.go              # Virtual account operations
├── models.go          # Request/response types and helper functions
├── response_codes.go  # BRI response code definitions and error handling
├── client_test.go     # Comprehensive test suite (75+ tests)
├── go.mod             # Go module definition
├── README.md          # This documentation
└── .gitignore         # Git ignore patterns
```

### Building

```bash
# Build the library
go build ./...

# Run tests
go test ./...

# Generate documentation
go doc -all
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Code Quality

- Follow Go naming conventions
- Add documentation for public APIs
- Include comprehensive test coverage
- Use `gofmt` for code formatting
- Run `go vet` and `golint` for static analysis

## Changelog

### v1.0.1 (Latest)
- **Fixed**: Unit test for `NewStructuredBRIAPIResponse` function - now correctly extracts HTTP status code from response definition instead of using hardcoded value
- **Tests**: All 75 unit tests now pass successfully

### v1.0.0
- Initial release with full BRIVA WS SNAP BI API implementation
- OAuth2 authentication with RSA signing
- HMAC-SHA512 request signing
- Complete virtual account operations (create, update, inquiry, delete, status, reports)
- Comprehensive error handling with structured response codes
- Debug mode for HTTP request/response logging
- Interface-based design for testability

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Note**: This library is not officially affiliated with Bank Rakyat Indonesia (BRI). Use at your own risk and ensure
compliance with BRI's terms of service.