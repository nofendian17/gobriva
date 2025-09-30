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
- **Comprehensive test suite** - 76+ unit tests with high code coverage
- **Debug mode** for HTTP request/response logging with timing measurements
- **Simplified error handling** - Direct parsing from API responses with automatic field extraction
- **HTTP status code-based categorization** - Standard error categorization without complex mappings
- **Context support** - All operations support context for cancellation and timeouts

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Security](#security)
- [Performance](#performance)
- [Development](#development)
- [License](#license)

## Installation

```bash
go get github.com/nofendian17/gobriva@latest
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
	"log/slog"

	"github.com/nofendian17/gobriva"
)

func main() {
	// Optional: create a custom slog.Logger and pass it to the client so the
	// library logs use your application's logger (no global changes necessary).
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	client := gobriva.NewClient(gobriva.Config{
		PartnerID:    "your-partner-id",
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		PrivateKey:   "your-private-key-pem",
		ChannelID:    "your-channel-id",
		IsSandbox:    true,
		Debug:        true,
		Logger:       logger, // pass your logger here
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

## Usage Examples

### Complete Virtual Account Workflow

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nofendian17/gobriva"
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

	ctx := context.Background()

	// 1. Create Virtual Account
	createReq := gobriva.NewCreateVirtualAccountRequest(
		"12345",                    // partnerServiceID
		"CUST001",                  // customerNo
		"12345678901234567890",     // virtualAccountNo
		"John Doe - Invoice #001",  // virtualAccountName
		"TRX001",                   // trxID
		100000.00,                  // amount
		"IDR",                      // currency
		"2024-12-31T23:59:59+07:00", // expiredDate
	)

	va, err := client.CreateVirtualAccount(ctx, createReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created VA: %s\n", va.VirtualAccountData.VirtualAccountNo)

	// 2. Update Virtual Account
	updateReq := gobriva.NewUpdateVirtualAccountRequest(
		"12345",                    // partnerServiceID
		"CUST001",                  // customerNo
		"12345678901234567890",     // virtualAccountNo
		"John Doe - Invoice #001",  // virtualAccountName
		"TRX001",                   // trxID
		150000.00,                  // updated amount
		"IDR",                      // currency
		"2024-12-31T23:59:59+07:00", // expiredDate
	)

	updatedVA, err := client.UpdateVirtualAccount(ctx, updateReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated VA: %s\n", updatedVA.VirtualAccountData.VirtualAccountNo)

	// 3. Inquiry Virtual Account
	inquiryReq := gobriva.NewInquiryVirtualAccountRequest(
		"12345",                // partnerServiceID
		"CUST001",              // customerNo
		"12345678901234567890", // virtualAccountNo
		"TRX002",               // trxID
	)

	vaInfo, err := client.InquiryVirtualAccount(ctx, inquiryReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("VA Info: %+v\n", vaInfo)

	// 4. Update Status to Paid
	statusReq := gobriva.NewUpdateVirtualAccountStatusRequest(
		"12345",                // partnerServiceID
		"CUST001",              // customerNo
		"12345678901234567890", // virtualAccountNo
		"TRX003",               // trxID
		"Y",                    // paidStatus (Y=paid, N=unpaid)
	)

	statusResp, err := client.UpdateVirtualAccountStatus(ctx, statusReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status updated: %s\n", statusResp.ResponseMessage)

	// 5. Get Transaction Report
	reportReq := gobriva.NewVirtualAccountReportRequest(
		"12345",     // partnerServiceID
		"2024-01-01", // startDate
		"00:00:00",   // startTime
		"23:59:59",   // endTime
	)

	report, err := client.GetVirtualAccountReport(ctx, reportReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d transactions\n", len(report.VirtualAccountData))
}
```

### Error Handling with Field Extraction

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/nofendian17/gobriva"
)

func handleVirtualAccountError(err error) {
	var briErr *gobriva.StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		fmt.Printf("BRI API Error [%s]: %s\n", briErr.ResponseCode, briErr.ResponseMessage)
		fmt.Printf("HTTP Status: %d\n", briErr.HTTPStatusCode)
		fmt.Printf("Category: %s\n", briErr.GetCategory())

		// Check error type
		switch briErr.GetCategory() {
		case gobriva.CategoryBadRequest:
			fmt.Println("Bad request - check your input data")
			if briErr.IsClientError() {
				fmt.Println("Client error (4xx) - validate request parameters")
			}
		case gobriva.CategoryUnauthorized:
			fmt.Println("Authentication failed - check credentials")
		case gobriva.CategoryInternalServerError:
			fmt.Println("Server error - retry later or contact support")
		case gobriva.CategoryPending:
			fmt.Println("Unknown response code - manual verification required")
		}

		// The error message automatically includes field extraction
		fmt.Printf("Full error: %s\n", briErr.Error())
	}
}

func main() {
	client := gobriva.NewClient(gobriva.Config{
		PartnerID:    "your-partner-id",
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		PrivateKey:   "your-private-key-pem",
		ChannelID:    "your-channel-id",
		IsSandbox:    true,
	})

	// Example with invalid data (will trigger field extraction)
	req := &gobriva.CreateVirtualAccountRequest{
		PartnerServiceID:   "", // Empty - will cause error
		CustomerNo:         "CUST001",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: gobriva.Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "TRX001",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err != nil {
		handleVirtualAccountError(err)
	}
}
```

### Custom HTTP Client and Timeouts

```go
package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/nofendian17/gobriva"
)

func main() {
	// Custom HTTP client with timeout and TLS config
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Set to true for testing only
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	client := gobriva.NewClient(gobriva.Config{
		PartnerID:    "your-partner-id",
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		PrivateKey:   "your-private-key-pem",
		ChannelID:    "your-channel-id",
		IsSandbox:    true,
		Debug:        false,
		HTTPClient:   httpClient, // Use custom HTTP client
	})

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	req := gobriva.NewCreateVirtualAccountRequest(
		"12345", "CUST001", "12345678901234567890",
		"Test Account", "TRX001", 100000.00, "IDR",
		"2024-12-31T23:59:59+07:00",
	)

	resp, err := client.CreateVirtualAccount(ctx, req)
	if err != nil {
		// Handle error
		return
	}

	fmt.Printf("Created VA: %s\n", resp.VirtualAccountData.VirtualAccountNo)
}
```

### Testing with Mock Dependencies

```go
package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/nofendian17/gobriva"
)

// MockHTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// MockAuthenticator for testing
type MockAuthenticator struct {
	EnsureAuthenticatedFunc func(context.Context) error
}

func (m *MockAuthenticator) EnsureAuthenticated(ctx context.Context) error {
	return m.EnsureAuthenticatedFunc(ctx)
}

func TestCreateVirtualAccount(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Success",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "CUST001",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Test Account",
						"totalAmount": {"value": "100000.00", "currency": "IDR"},
						"expiredDate": "2024-12-31T23:59:59+07:00",
						"trxId": "TRX001"
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := gobriva.NewClient(gobriva.Config{
		PartnerID:     "test-partner",
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		PrivateKey:    "test-key",
		ChannelID:     "test-channel",
		IsSandbox:     true,
		HTTPClient:    mockHTTP,
		Authenticator: mockAuth,
	})

	req := gobriva.NewCreateVirtualAccountRequest(
		"12345", "CUST001", "12345678901234567890",
		"Test Account", "TRX001", 100000.00, "IDR",
		"2024-12-31T23:59:59+07:00",
	)

	resp, err := client.CreateVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if resp.VirtualAccountData.VirtualAccountNo != "12345678901234567890" {
		t.Errorf("Expected VA number '12345678901234567890', got '%s'", resp.VirtualAccountData.VirtualAccountNo)
	}
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
	Logger        *slog.Logger  // Optional: pass a custom slog.Logger; client will use it locally
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

The private key should be in PEM format with 2048 bits:

```pem
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC...
-----END PRIVATE KEY-----
```

## API Reference

### Client Creation

```go
func NewClient(config Config) *Client
```

Creates a new BRI Virtual Account client with the provided configuration.

### Virtual Account Operations

#### CreateVirtualAccount

Creates a new virtual account.

```go
func (c *Client) CreateVirtualAccount(ctx context.Context, req *CreateVirtualAccountRequest) (*CreateVirtualAccountResponse, error)
```

**Request Construction:**

```go
// Using struct
req := &CreateVirtualAccountRequest{
    PartnerServiceID:   "12345",
    CustomerNo:         "CUST001",
    VirtualAccountNo:   "12345678901234567890",
    VirtualAccountName: "John Doe",
    TotalAmount: Amount{
        Value:    "100000.00",
        Currency: "IDR",
    },
    ExpiredDate: "2024-12-31T23:59:59+07:00",
    TrxID:       "TRX001",
}

// Using helper function
req := NewCreateVirtualAccountRequest(
    "12345",                // partnerServiceID
    "CUST001",              // customerNo
    "12345678901234567890", // virtualAccountNo
    "John Doe",             // virtualAccountName
    "TRX001",               // trxID
    100000.00,              // amount
    "IDR",                  // currency
    "2024-12-31T23:59:59+07:00", // expiredDate
)
```

#### UpdateVirtualAccount

Updates an existing virtual account.

```go
func (c *Client) UpdateVirtualAccount(ctx context.Context, req *UpdateVirtualAccountRequest) (*UpdateVirtualAccountResponse, error)
```

**Request Construction:**

```go
req := NewUpdateVirtualAccountRequest(
    "12345",                // partnerServiceID
    "CUST001",              // customerNo
    "12345678901234567890", // virtualAccountNo
    "John Doe Updated",     // virtualAccountName
    "TRX002",               // trxID
    150000.00,              // amount
    "IDR",                  // currency
    "2024-12-31T23:59:59+07:00", // expiredDate
)
```

#### UpdateVirtualAccountStatus

Updates the payment status of a virtual account.

```go
func (c *Client) UpdateVirtualAccountStatus(ctx context.Context, req *UpdateVirtualAccountStatusRequest) (*UpdateVirtualAccountStatusResponse, error)
```

**Request Construction:**

```go
req := NewUpdateVirtualAccountStatusRequest(
    "12345",                // partnerServiceID
    "CUST001",              // customerNo
    "12345678901234567890", // virtualAccountNo
    "TRX003",               // trxID
    "Y",                    // paidStatus ("Y" = paid, "N" = unpaid)
)
```

#### InquiryVirtualAccount

Retrieves information about a virtual account.

```go
func (c *Client) InquiryVirtualAccount(ctx context.Context, req *InquiryVirtualAccountRequest) (*InquiryVirtualAccountResponse, error)
```

**Request Construction:**

```go
req := NewInquiryVirtualAccountRequest(
    "12345",                // partnerServiceID
    "CUST001",              // customerNo
    "12345678901234567890", // virtualAccountNo
    "TRX004",               // trxID
)
```

#### InquiryVirtualAccountStatus

Retrieves the payment status of a virtual account.

```go
func (c *Client) InquiryVirtualAccountStatus(ctx context.Context, req *InquiryVirtualAccountStatusRequest) (*InquiryVirtualAccountStatusResponse, error)
```

#### DeleteVirtualAccount

Deletes a virtual account.

```go
func (c *Client) DeleteVirtualAccount(ctx context.Context, req *DeleteVirtualAccountRequest) (*DeleteVirtualAccountResponse, error)
```

#### GetVirtualAccountReport

Retrieves transaction reports for virtual accounts within a date range.

```go
func (c *Client) GetVirtualAccountReport(ctx context.Context, req *VirtualAccountReportRequest) (*VirtualAccountReportResponse, error)
```

**Request Construction:**

```go
req := NewVirtualAccountReportRequest(
    "12345",     // partnerServiceID
    "2024-01-01", // startDate (YYYY-MM-DD)
    "00:00:00",   // startTime (HH:MM:SS)
    "23:59:59",   // endTime (HH:MM:SS)
)
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

The library provides comprehensive error handling with structured error types that parse directly from API responses.

### Error Types

#### StructuredBRIAPIResponse

```go
type StructuredBRIAPIResponse struct {
	ResponseCode    string    // The actual response code from API
	ResponseMessage string    // The actual response message from API
	HTTPStatusCode  int       // HTTP status code
	Timestamp       time.Time // When the error occurred
}
```

### Error Categories

Errors are categorized based on HTTP status codes:

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
func (e *StructuredBRIAPIResponse) GetCategory() HttpCategory
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

		// Check if it's a client error (4xx)
		if briErr.IsClientError() {
			log.Printf("Client error - check your request data")
		}

		// Check if response code is unknown/pending
		if briErr.IsPending() {
			log.Printf("Unknown response code - manual verification required")
		}
	}
	return err
}
```

### Field Extraction from Errors

The library automatically extracts field names from error messages for better debugging:

```go
// Example API error response:
// {"responseCode": "4002701", "responseMessage": "Invalid Mandatory Field institutionCode"}

_, err := client.CreateVirtualAccount(context.Background(), req)
if briErr, ok := err.(*gobriva.StructuredBRIAPIResponse); ok {
	// Error message will include field extraction:
	// "BRI API Error [4002701]: Invalid Mandatory Field institutionCode (field: institutionCode)"
	log.Printf("Error: %s", briErr.Error())
}
```

## Testing

### Test Coverage

Current test coverage: **76+ unit tests** with comprehensive coverage of all operations and error scenarios.

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

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Note**: This library is not officially affiliated with Bank Rakyat Indonesia (BRI). Use at your own risk and ensure
compliance with BRI's terms of service.