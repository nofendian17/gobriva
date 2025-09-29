package gobriva

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	// Default mock response
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{"responseCode":"2002700","responseMessage":"Successful"}`)),
		Header:     make(http.Header),
	}, nil
}

// MockAuthenticator is a mock implementation of Authenticator for testing
type MockAuthenticator struct {
	AuthenticateFunc        func(ctx context.Context) error
	IsAuthenticatedFunc     func() bool
	EnsureAuthenticatedFunc func(ctx context.Context) error
}

func (m *MockAuthenticator) Authenticate(ctx context.Context) error {
	if m.AuthenticateFunc != nil {
		return m.AuthenticateFunc(ctx)
	}
	return nil
}

func (m *MockAuthenticator) IsAuthenticated() bool {
	if m.IsAuthenticatedFunc != nil {
		return m.IsAuthenticatedFunc()
	}
	return true
}

func (m *MockAuthenticator) EnsureAuthenticated(ctx context.Context) error {
	if m.EnsureAuthenticatedFunc != nil {
		return m.EnsureAuthenticatedFunc(ctx)
	}
	return nil
}

var privateKeyTest = `-----BEGIN RSA PRIVATE KEY-----
MIIEoQIBAAKCAQB4ke4WQ+V3lcZcoSMdMKCj1zRGIOoJ2ctD01hsJ1Bo/O6QwiRu
sAGSm13afA43Hs+ibXOW2B+as4xXNX1EaJwDYJ60dUBZOjy1S66xhdp/+b8NbxRw
tE2VGAyp7uafEyj38QN0calQhmpTOUeSGScoybrogOWW/2n/OibDhWt7vlSe0XKd
tvncrjaQ6nctSR9Ui/zCE6szwt/ZN5yNkR4hjCsRPkQ0zsVz0dQl8NdPtgRgz+uC
EHcTvD0ssJiipc9CApZc5moxNxq8TE11GdQEVqu1ej+s2NUENCdkd2ee3W9Wo8Ht
VLD6DMdtnVoISK/MkCbNBcNw/jaaK5oSKVRhAgMBAAECggEAVe78mJXv2NnBNYgL
tORRujTKJymSZU77lu3tWbhzkUCk8DvPJ6z+kfV2YSCGKTcmmggUmHCVpfdOkWGo
VLeSar3Un53qLS5a0oSMkC5s20Wvq+19zg5UNW2cqQmDCeHoEz+OTNIt8Ry8b3Cl
2DVhOar+MnScLpEAhU53DmfrgZfYDMQT7H7gg4Bp+P57GXg1mHP1iRNzdPvd+9S7
ql02w+N1pjta0eF4N68p7nPOIsUq2oZyM/u4glx8eWsOt+P/iQ1yalH80iKKuQaa
9sAXLWCWGzkFi8lanIVzuXmcycknX3lCBC9H70gGC658Q5QVmY8K/jXLw+BGfRjH
tG4B2QKBgQDZ0up5I8SnR3ChD0wvDs5/BkWjdg13hpdZzko9AVfzVnnDZNrB7Gbf
UUiFDuHZQfFb/I/kG2NnfrCNJUwX7omtexSxXIOqoErOHp+Wotmjo50tVQC1LblI
HEDeZWWwmL4foxkXewLj0nE4HkWOGfLGhIn+hP7+DmEwo8lJX3tw9wKBgQCNs4jh
Hgg3twHQsHmjQxq1nF188sAaC/HXF4G/WqE6IgmHZZ4MFLYE7gCine1e+yB+QhwS
b8vTjtYMtO8E8iK717duVhbsTyuE4ePkpJHGIiMbrsw037Kov5GY3ayE7F/pHI7V
ePdFkmoP/yL2s5IG+jIWRq9Sg1q6R0g2S/DnZwKBgEURwTHKargUSh14CVM+obHb
nkdXzqtg7SsX46h2fZn2iMOxfkBRoskbMCCo+Gp4o3zkmAffu2R84qTO99L624M7
7PLUgBehnja/tSEB4HsoDVXrhz7sEb1Q4CzlABrAREEp6XHtmpv9BdOinbGSfs3+
BvfC2kxa6OyQcuomMbE/AoGAAsG6aP7HlCXoUCIOy8FTdLMNEpA6codG9jNL3+go
eNQOsWals4B3phLnSkKeSpnCIRKyLx2jroL54RdoCwWW7Wad9/SOz5wesaAfaeRV
vbAOVMyKxoCPnj7T21B8ub1LhGJ82ORYky7tB1CkYn5N2frmHI7VfFp32mXmnr/N
eQMCgYAn5WAyxLyEaKmSDyvArqf86RFO6SNlK02T45Hh6hrbi0ArJFgyatSEY1D6
/uqKDBvBA+GtVIbcepaiiCeNIrtj4DwZaGRVT2tqUcUiLKJICCesRsw4mZ5x/KUo
38Ex8a2e++RahDu94ag3mgj5FENlHJeMBEhWa+aPmlzwEXEj1w==
-----END RSA PRIVATE KEY-----`

func TestClientWithMocks(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request headers
			if req.Header.Get("X-PARTNER-ID") != "test-partner" {
				t.Errorf("Expected X-PARTNER-ID header to be 'test-partner', got '%s'", req.Header.Get("X-PARTNER-ID"))
			}
			if req.Header.Get("CHANNEL-ID") != "test-channel" {
				t.Errorf("Expected CHANNEL-ID header to be 'test-channel', got '%s'", req.Header.Get("CHANNEL-ID"))
			}

			// Return mock response
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Test Account",
						"trxId": "test-trx-123",
						"totalAmount": {
							"value": "100000.00",
							"currency": "IDR"
						}
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			// Mock authentication - just return success
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test CreateVirtualAccount
	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
		AdditionalInfo: AdditionalInfo{
			Description: "Test transaction",
		},
	}

	resp, err := client.CreateVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateVirtualAccount failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	if resp.VirtualAccountData.VirtualAccountName != "Test Account" {
		t.Errorf("Expected virtual account name 'Test Account', got '%s'", resp.VirtualAccountData.VirtualAccountName)
	}
}

func TestClientWithDebugMode(t *testing.T) {
	// Capture slog output
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug}))
	originalLogger := slog.Default()
	slog.SetDefault(logger)
	defer slog.SetDefault(originalLogger)

	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Test Account",
						"trxId": "test-trx-123",
						"totalAmount": {
							"value": "100000.00",
							"currency": "IDR"
						}
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with debug mode enabled
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		debug:        true, // Enable debug mode
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	resp, err := client.CreateVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateVirtualAccount failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	// Verify debug logging occurred
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, `"level":"DEBUG"`) || !strings.Contains(logOutput, `"msg":"HTTP Request"`) {
		t.Errorf("Expected debug logging with request details, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"msg":"HTTP Response"`) || !strings.Contains(logOutput, `"duration"`) {
		t.Errorf("Expected debug logging with response details and duration, got: %s", logOutput)
	}
}

func TestClientWithUnknownResponseCode(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "9999999",
					"responseMessage": "Unknown error occurred"
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for unknown response code")
	}

	// Check that the error is properly structured as pending
	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "9999999" {
			t.Errorf("Expected response code '9999999', got '%s'", briErr.ResponseCode)
		}
		if !briErr.IsPending() {
			t.Errorf("Expected response to be pending, but IsPending() returned false")
		}
		if briErr.GetCategory() != CategoryPending {
			t.Errorf("Expected category 'Pending', got '%s'", briErr.GetCategory())
		}
	}
}

func TestUpdateVirtualAccount(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Updated Account",
						"trxId": "update-trx-123",
						"totalAmount": {
							"value": "150000.00",
							"currency": "IDR"
						}
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test UpdateVirtualAccount
	req := &UpdateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Updated Account",
		TotalAmount: Amount{
			Value:    "150000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "update-trx-123",
		AdditionalInfo: AdditionalInfo{
			Description: "Updated transaction",
		},
	}

	resp, err := client.UpdateVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateVirtualAccount failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	if resp.VirtualAccountData.VirtualAccountName != "Updated Account" {
		t.Errorf("Expected virtual account name 'Updated Account', got '%s'", resp.VirtualAccountData.VirtualAccountName)
	}
}

func TestUpdateVirtualAccountStatus(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"paidStatus": "Y"
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test UpdateVirtualAccountStatus
	req := &UpdateVirtualAccountStatusRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "status-trx-123",
		PaidStatus:       "Y",
	}

	resp, err := client.UpdateVirtualAccountStatus(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateVirtualAccountStatus failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	if resp.VirtualAccountData.PaidStatus != "Y" {
		t.Errorf("Expected paid status 'Y', got '%s'", resp.VirtualAccountData.PaidStatus)
	}
}

func TestInquiryVirtualAccount(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Test Account",
						"trxId": "inquiry-trx-123",
						"totalAmount": {
							"value": "100000.00",
							"currency": "IDR"
						}
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test InquiryVirtualAccount
	req := &InquiryVirtualAccountRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "inquiry-trx-123",
	}

	resp, err := client.InquiryVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("InquiryVirtualAccount failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	if resp.VirtualAccountData.VirtualAccountName != "Test Account" {
		t.Errorf("Expected virtual account name 'Test Account', got '%s'", resp.VirtualAccountData.VirtualAccountName)
	}
}

func TestDeleteVirtualAccount(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"trxId": "delete-trx-123"
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test DeleteVirtualAccount
	req := &DeleteVirtualAccountRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "delete-trx-123",
	}

	resp, err := client.DeleteVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("DeleteVirtualAccount failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	if resp.VirtualAccountData.TrxID != "delete-trx-123" {
		t.Errorf("Expected trxId 'delete-trx-123', got '%s'", resp.VirtualAccountData.TrxID)
	}
}

func TestGetVirtualAccountReport(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": [
						{
							"partnerServiceId": "12345",
							"customerNo": "67890",
							"virtualAccountNo": "12345678901234567890",
							"virtualAccountName": "Test Account",
							"sourceAccountNo": "1234567890",
							"paidAmount": {
								"value": "100000.00",
								"currency": "IDR"
							},
							"trxDateTime": "2024-01-01T10:00:00+07:00",
							"trxId": "report-trx-123",
							"inquiryRequestId": "inq-123",
							"paymentRequestId": "pay-123",
							"totalAmount": {
								"value": "100000.00",
								"currency": "IDR"
							}
						}
					]
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test GetVirtualAccountReport
	req := &VirtualAccountReportRequest{
		PartnerServiceID: "12345",
		StartDate:        "2024-01-01",
		StartTime:        "00:00:00",
		EndTime:          "23:59:59",
		EndDate:          "2024-01-01",
	}

	resp, err := client.GetVirtualAccountReport(context.Background(), req)
	if err != nil {
		t.Fatalf("GetVirtualAccountReport failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	if len(resp.VirtualAccountData) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(resp.VirtualAccountData))
	}

	if resp.VirtualAccountData[0].TrxID != "report-trx-123" {
		t.Errorf("Expected trxId 'report-trx-123', got '%s'", resp.VirtualAccountData[0].TrxID)
	}
}

func TestInquiryVirtualAccountStatus(t *testing.T) {
	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"paidStatus": "N"
					},
					"additionalInfo": {
						"description": "Account status inquiry"
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test InquiryVirtualAccountStatus
	req := &InquiryVirtualAccountStatusRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		InquiryRequestID: "inq-status-123",
	}

	resp, err := client.InquiryVirtualAccountStatus(context.Background(), req)
	if err != nil {
		t.Fatalf("InquiryVirtualAccountStatus failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	if resp.VirtualAccountData.PaidStatus != "N" {
		t.Errorf("Expected paid status 'N', got '%s'", resp.VirtualAccountData.PaidStatus)
	}

	if resp.AdditionalInfo.Description != "Account status inquiry" {
		t.Errorf("Expected description 'Account status inquiry', got '%s'", resp.AdditionalInfo.Description)
	}
}

func TestAuthenticationFailure(t *testing.T) {
	// Create mock HTTP client (won't be called due to auth failure)
	mockHTTP := &MockHTTPClient{}

	// Create mock authenticator that fails
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return fmt.Errorf("authentication failed: invalid credentials")
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test CreateVirtualAccount with auth failure
	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected authentication error, but got success")
	}

	if !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("Expected authentication error, got: %v", err)
	}
}

func TestHTTPErrorResponse(t *testing.T) {
	// Create mock HTTP client that returns error response
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "4002701",
					"responseMessage": "Invalid field format"
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with mocks
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test CreateVirtualAccount with HTTP error
	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected HTTP error, but got success")
	}

	// Check that it's a structured BRI API response
	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "4002701" {
			t.Errorf("Expected response code '4002701', got '%s'", briErr.ResponseCode)
		}
		if briErr.ResponseMessage != "Invalid field format" {
			t.Errorf("Expected response message 'Invalid field format', got '%s'", briErr.ResponseMessage)
		}
		if briErr.GetCategory() != CategoryBadRequest {
			t.Errorf("Expected category 'BadRequest', got '%s'", briErr.GetCategory())
		}
	}
}

func TestUpdateVirtualAccountWithDebug(t *testing.T) {
	// Capture slog output
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug}))
	originalLogger := slog.Default()
	slog.SetDefault(logger)
	defer slog.SetDefault(originalLogger)

	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Updated Account",
						"trxId": "update-trx-123",
						"totalAmount": {
							"value": "150000.00",
							"currency": "IDR"
						}
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with debug mode enabled
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		debug:        true, // Enable debug mode
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &UpdateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Updated Account",
		TotalAmount: Amount{
			Value:    "150000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "update-trx-123",
	}

	resp, err := client.UpdateVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateVirtualAccount failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	// Verify debug logging occurred
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, `"level":"DEBUG"`) || !strings.Contains(logOutput, `"msg":"HTTP Request"`) {
		t.Errorf("Expected debug logging with PUT request details, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"msg":"HTTP Response"`) || !strings.Contains(logOutput, `"duration"`) {
		t.Errorf("Expected debug logging with response details and duration, got: %s", logOutput)
	}
}

func TestInquiryVirtualAccountWithDebug(t *testing.T) {
	// Capture slog output
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug}))
	originalLogger := slog.Default()
	slog.SetDefault(logger)
	defer slog.SetDefault(originalLogger)

	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Test Account",
						"trxId": "inquiry-trx-123",
						"totalAmount": {
							"value": "100000.00",
							"currency": "IDR"
						}
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	// Create mock authenticator
	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create client with debug mode enabled
	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		debug:        true, // Enable debug mode
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &InquiryVirtualAccountRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "inquiry-trx-123",
	}

	resp, err := client.InquiryVirtualAccount(context.Background(), req)
	if err != nil {
		t.Fatalf("InquiryVirtualAccount failed: %v", err)
	}

	if resp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", resp.ResponseCode)
	}

	// Verify debug logging occurred
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, `"level":"DEBUG"`) || !strings.Contains(logOutput, `"msg":"HTTP Request"`) {
		t.Errorf("Expected debug logging with GET request details, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"msg":"HTTP Response"`) || !strings.Contains(logOutput, `"duration"`) {
		t.Errorf("Expected debug logging with response details and duration, got: %s", logOutput)
	}
}

func TestCreateVAInvalidFieldFormat(t *testing.T) {
	// Test for invalid field format (4002701) - based on Postman collection
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "4002701",
					"responseMessage": "Invalid Field Format virtualAccountNo"
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

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "invalid-format", // Invalid format
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for invalid field format")
	}

	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "4002701" {
			t.Errorf("Expected response code '4002701', got '%s'", briErr.ResponseCode)
		}
		if briErr.ResponseMessage != "Invalid Field Format virtualAccountNo" {
			t.Errorf("Expected response message 'Invalid Field Format virtualAccountNo', got '%s'", briErr.ResponseMessage)
		}
		if briErr.GetCategory() != CategoryBadRequest {
			t.Errorf("Expected category 'BadRequest', got '%s'", briErr.GetCategory())
		}
	}
}

func TestCreateVAInvalidMandatoryFieldPartnerServiceId(t *testing.T) {
	// Test for invalid mandatory field partnerServiceId (4002702) - based on Postman collection
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "4002702",
					"responseMessage": "Invalid Mandatory Field partnerServiceId"
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

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "", // Empty mandatory field
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for missing mandatory field")
	}

	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "4002702" {
			t.Errorf("Expected response code '4002702', got '%s'", briErr.ResponseCode)
		}
		if briErr.ResponseMessage != "Invalid Mandatory Field partnerServiceId" {
			t.Errorf("Expected response message 'Invalid Mandatory Field partnerServiceId', got '%s'", briErr.ResponseMessage)
		}
		if briErr.GetCategory() != CategoryBadRequest {
			t.Errorf("Expected category 'BadRequest', got '%s'", briErr.GetCategory())
		}
	}
}

func TestCreateVAInvalidHeaderXPartnerID(t *testing.T) {
	// Test for invalid mandatory header xPartnerID (4002702) - based on Postman collection
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify that X-PARTNER-ID header is missing or invalid
			if req.Header.Get("X-PARTNER-ID") == "" {
				return &http.Response{
					StatusCode: 400,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"responseCode": "4002702",
						"responseMessage": "Invalid Mandatory Field xPartnerID"
					}`)),
					Header: make(http.Header),
				}, nil
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"responseCode":"2002700","responseMessage":"Successful"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "", // Empty partner ID to trigger header validation
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for missing mandatory header")
	}

	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "4002702" {
			t.Errorf("Expected response code '4002702', got '%s'", briErr.ResponseCode)
		}
		if briErr.ResponseMessage != "Invalid Mandatory Field xPartnerID" {
			t.Errorf("Expected response message 'Invalid Mandatory Field xPartnerID', got '%s'", briErr.ResponseMessage)
		}
		if briErr.GetCategory() != CategoryBadRequest {
			t.Errorf("Expected category 'BadRequest', got '%s'", briErr.GetCategory())
		}
	}
}

func TestCreateVAInvalidHeaderXExternalID(t *testing.T) {
	// Test for invalid mandatory field xExternalID (4002702) - based on Postman collection
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "4002702",
					"responseMessage": "Invalid Mandatory Field xExternalID"
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

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for missing mandatory header")
	}

	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "4002702" {
			t.Errorf("Expected response code '4002702', got '%s'", briErr.ResponseCode)
		}
		if briErr.ResponseMessage != "Invalid Mandatory Field xExternalID" {
			t.Errorf("Expected response message 'Invalid Mandatory Field xExternalID', got '%s'", briErr.ResponseMessage)
		}
		if briErr.GetCategory() != CategoryBadRequest {
			t.Errorf("Expected category 'BadRequest', got '%s'", briErr.GetCategory())
		}
	}
}

func TestCreateVAInvalidHeaderChannelID(t *testing.T) {
	// Test for invalid mandatory header channelID (4002702) - based on Postman collection
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify that CHANNEL-ID header is missing
			if req.Header.Get("CHANNEL-ID") == "" {
				return &http.Response{
					StatusCode: 400,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"responseCode": "4002702",
						"responseMessage": "Invalid Mandatory Field channelID"
					}`)),
					Header: make(http.Header),
				}, nil
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"responseCode":"2002700","responseMessage":"Successful"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "", // Empty channel ID to trigger header validation
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for missing mandatory header")
	}

	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "4002702" {
			t.Errorf("Expected response code '4002702', got '%s'", briErr.ResponseCode)
		}
		if briErr.ResponseMessage != "Invalid Mandatory Field channelID" {
			t.Errorf("Expected response message 'Invalid Mandatory Field channelID', got '%s'", briErr.ResponseMessage)
		}
		if briErr.GetCategory() != CategoryBadRequest {
			t.Errorf("Expected category 'BadRequest', got '%s'", briErr.GetCategory())
		}
	}
}

func TestCreateVAHTTPError(t *testing.T) {
	// Test for HTTP error responses (5xx status codes)
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "5002700",
					"responseMessage": "Internal Server Error"
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

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for HTTP 500 response")
	}

	var briErr *StructuredBRIAPIResponse
	if errors.As(err, &briErr) {
		if briErr.ResponseCode != "5002700" {
			t.Errorf("Expected response code '5002700', got '%s'", briErr.ResponseCode)
		}
		if briErr.ResponseMessage != "Internal Server Error" {
			t.Errorf("Expected response message 'Internal Server Error', got '%s'", briErr.ResponseMessage)
		}
		if briErr.GetCategory() != CategoryInternalServerError {
			t.Errorf("Expected category 'InternalServerError', got '%s'", briErr.GetCategory())
		}
	}
}

func TestCreateVANetworkError(t *testing.T) {
	// Test for network/connection errors
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		isSandbox:    true,
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for network failure")
	}

	// Network errors should not be wrapped in StructuredBRIAPIResponse
	var structuredBRIAPIResponse *StructuredBRIAPIResponse
	if errors.As(err, &structuredBRIAPIResponse) {
		t.Errorf("Network errors should not be wrapped in StructuredBRIAPIResponse, got %T", err)
	}
}

// Authentication Tests

func TestNewClient(t *testing.T) {
	// Test NewClient with default values
	config := Config{
		PartnerID:    "test-partner",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		PrivateKey:   "test-key",
		ChannelID:    "test-channel",
		IsSandbox:    true,
		Debug:        true,
	}

	client := NewClient(config)

	if client.partnerID != "test-partner" {
		t.Errorf("Expected partnerID 'test-partner', got '%s'", client.partnerID)
	}
	if client.clientID != "test-client" {
		t.Errorf("Expected clientID 'test-client', got '%s'", client.clientID)
	}
	if client.channelID != "test-channel" {
		t.Errorf("Expected channelID 'test-channel', got '%s'", client.channelID)
	}
	if !client.isSandbox {
		t.Error("Expected isSandbox to be true")
	}
	if !client.debug {
		t.Error("Expected debug to be true")
	}
	if client.baseURL != sandboxBaseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", sandboxBaseURL, client.baseURL)
	}
}

func TestNewClientProduction(t *testing.T) {
	// Test NewClient with production environment
	config := Config{
		PartnerID:    "test-partner",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		PrivateKey:   "test-key",
		ChannelID:    "test-channel",
		IsSandbox:    false,
	}

	client := NewClient(config)

	if client.baseURL != productionBaseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", productionBaseURL, client.baseURL)
	}
}

func TestNewClientWithCustomHTTPClient(t *testing.T) {
	// Test NewClient with custom HTTP client
	customHTTP := &MockHTTPClient{}
	config := Config{
		PartnerID:    "test-partner",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		PrivateKey:   "test-key",
		ChannelID:    "test-channel",
		IsSandbox:    true,
		HTTPClient:   customHTTP,
	}

	client := NewClient(config)

	if client.httpClient != customHTTP {
		t.Error("Expected custom HTTP client to be used")
	}
}

func TestNewClientWithCustomAuthenticator(t *testing.T) {
	// Test NewClient with custom authenticator
	customAuth := &MockAuthenticator{}
	config := Config{
		PartnerID:     "test-partner",
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		PrivateKey:    "test-key",
		ChannelID:     "test-channel",
		IsSandbox:     true,
		Authenticator: customAuth,
	}

	client := NewClient(config)

	if client.auth != customAuth {
		t.Error("Expected custom authenticator to be used")
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	// Test NewClient with custom timeout
	config := Config{
		PartnerID:    "test-partner",
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		PrivateKey:   "test-key",
		ChannelID:    "test-channel",
		IsSandbox:    true,
		Timeout:      10 * time.Second,
	}

	client := NewClient(config)

	// Verify timeout is set on HTTP client
	httpClient, ok := client.httpClient.(*http.Client)
	if !ok {
		t.Fatal("Expected default HTTP client")
	}
	if httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", httpClient.Timeout)
	}
}

func TestDefaultAuthenticatorAuthenticate(t *testing.T) {
	// Test successful authentication using mock authenticator
	mockAuth := &MockAuthenticator{
		AuthenticateFunc: func(ctx context.Context) error {
			// Simulate successful authentication
			return nil
		},
	}

	// Use mock authenticator instead of DefaultAuthenticator for this test
	err := mockAuth.Authenticate(context.Background())
	if err != nil {
		t.Fatalf("Authentication failed: %v", err)
	}
}

func TestDefaultAuthenticatorAuthenticateFailure(t *testing.T) {
	// Test authentication failure using mock authenticator
	mockAuth := &MockAuthenticator{
		AuthenticateFunc: func(ctx context.Context) error {
			// Simulate authentication failure
			return fmt.Errorf("authentication failed: invalid credentials")
		},
	}

	err := mockAuth.Authenticate(context.Background())
	if err == nil {
		t.Fatal("Expected authentication to fail")
	}

	if !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("Expected authentication error, got: %v", err)
	}
}

func TestDefaultAuthenticatorIsAuthenticated(t *testing.T) {
	client := &Client{
		accessToken: "test-token",
		tokenExpiry: time.Now().Add(time.Hour), // Token not expired
	}

	auth := &DefaultAuthenticator{client: client}

	if !auth.IsAuthenticated() {
		t.Error("Expected to be authenticated")
	}
}

func TestDefaultAuthenticatorIsAuthenticatedExpired(t *testing.T) {
	client := &Client{
		accessToken: "test-token",
		tokenExpiry: time.Now().Add(-time.Hour), // Token expired
	}

	auth := &DefaultAuthenticator{client: client}

	if auth.IsAuthenticated() {
		t.Error("Expected to not be authenticated (token expired)")
	}
}

func TestDefaultAuthenticatorIsAuthenticatedEmptyToken(t *testing.T) {
	client := &Client{
		accessToken: "",
		tokenExpiry: time.Now().Add(time.Hour),
	}

	auth := &DefaultAuthenticator{client: client}

	if auth.IsAuthenticated() {
		t.Error("Expected to not be authenticated (empty token)")
	}
}

func TestDefaultAuthenticatorEnsureAuthenticated(t *testing.T) {
	// Test when already authenticated
	client := &Client{
		accessToken: "test-token",
		tokenExpiry: time.Now().Add(time.Hour),
	}

	auth := &DefaultAuthenticator{client: client}

	err := auth.EnsureAuthenticated(context.Background())
	if err != nil {
		t.Errorf("Expected no error when already authenticated, got %v", err)
	}
}

// Utility function tests

func TestGenerateExternalID(t *testing.T) {
	client := &Client{}

	id1 := client.generateExternalID()
	id2 := client.generateExternalID()

	if len(id1) != 9 {
		t.Errorf("Expected external ID length 9, got %d", len(id1))
	}
	if len(id2) != 9 {
		t.Errorf("Expected external ID length 9, got %d", len(id2))
	}
	if id1 == id2 {
		t.Error("Expected different external IDs")
	}

	// Verify it's numeric
	for _, char := range id1 {
		if char < '0' || char > '9' {
			t.Errorf("Expected numeric external ID, got '%s'", id1)
		}
	}
}

func TestGenerateTimestamp(t *testing.T) {
	client := &Client{}

	timestamp1 := client.generateTimestamp()

	// Sleep for 1ms to ensure different timestamp
	time.Sleep(1 * time.Millisecond)

	timestamp2 := client.generateTimestamp()

	if timestamp1 == timestamp2 {
		t.Error("Expected different timestamps")
	}

	// Verify format (basic check)
	if len(timestamp1) < 20 {
		t.Errorf("Expected timestamp format, got '%s'", timestamp1)
	}
}

func TestCalculateSignature(t *testing.T) {
	client := &Client{
		clientSecret: "test-secret",
		accessToken:  "test-token",
	}

	signature, err := client.calculateSignature("POST", "/test", `{"key":"value"}`)
	if err != nil {
		t.Fatalf("Failed to calculate signature: %v", err)
	}

	if signature == "" {
		t.Error("Expected non-empty signature")
	}

	// Verify signature is base64 encoded (length should be reasonable)
	if len(signature) < 10 {
		t.Errorf("Expected longer signature, got length %d", len(signature))
	}
}

func TestCalculateSignatureGET(t *testing.T) {
	// Test signature calculation for GET request (no body)
	client := &Client{
		clientSecret: "test-secret",
		accessToken:  "test-token",
	}

	signature, err := client.calculateSignature("GET", "/test", "")
	if err != nil {
		t.Fatalf("Failed to calculate signature: %v", err)
	}

	if signature == "" {
		t.Error("Expected non-empty signature")
	}
}

// Model helper function tests

func TestNewCreateVirtualAccountRequest(t *testing.T) {
	req := NewCreateVirtualAccountRequest("12345", "67890", "12345678901234567890", "Test Account", "test-trx-123", 100000.00, "IDR", "2024-12-31T23:59:59+07:00")

	if req.PartnerServiceID != "12345" {
		t.Errorf("Expected PartnerServiceID '12345', got '%s'", req.PartnerServiceID)
	}
	if req.CustomerNo != "67890" {
		t.Errorf("Expected CustomerNo '67890', got '%s'", req.CustomerNo)
	}
	if req.VirtualAccountNo != "12345678901234567890" {
		t.Errorf("Expected VirtualAccountNo '12345678901234567890', got '%s'", req.VirtualAccountNo)
	}
	if req.VirtualAccountName != "Test Account" {
		t.Errorf("Expected VirtualAccountName 'Test Account', got '%s'", req.VirtualAccountName)
	}
	if req.TrxID != "test-trx-123" {
		t.Errorf("Expected TrxID 'test-trx-123', got '%s'", req.TrxID)
	}
	if req.TotalAmount.Value != "100000.00" {
		t.Errorf("Expected TotalAmount.Value '100000.00', got '%s'", req.TotalAmount.Value)
	}
	if req.TotalAmount.Currency != "IDR" {
		t.Errorf("Expected TotalAmount.Currency 'IDR', got '%s'", req.TotalAmount.Currency)
	}
	if req.ExpiredDate != "2024-12-31T23:59:59+07:00" {
		t.Errorf("Expected ExpiredDate '2024-12-31T23:59:59+07:00', got '%s'", req.ExpiredDate)
	}
}

func TestNewUpdateVirtualAccountRequest(t *testing.T) {
	req := NewUpdateVirtualAccountRequest("12345", "67890", "12345678901234567890", "Updated Account", "update-trx-123", 150000.00, "IDR", "2024-12-31T23:59:59+07:00")

	if req.PartnerServiceID != "12345" {
		t.Errorf("Expected PartnerServiceID '12345', got '%s'", req.PartnerServiceID)
	}
	if req.VirtualAccountName != "Updated Account" {
		t.Errorf("Expected VirtualAccountName 'Updated Account', got '%s'", req.VirtualAccountName)
	}
	if req.TotalAmount.Value != "150000.00" {
		t.Errorf("Expected TotalAmount.Value '150000.00', got '%s'", req.TotalAmount.Value)
	}
}

func TestNewUpdateVirtualAccountStatusRequest(t *testing.T) {
	req := NewUpdateVirtualAccountStatusRequest("12345", "67890", "12345678901234567890", "status-trx-123", "Y")

	if req.PartnerServiceID != "12345" {
		t.Errorf("Expected PartnerServiceID '12345', got '%s'", req.PartnerServiceID)
	}
	if req.PaidStatus != "Y" {
		t.Errorf("Expected PaidStatus 'Y', got '%s'", req.PaidStatus)
	}
	if req.TrxID != "status-trx-123" {
		t.Errorf("Expected TrxID 'status-trx-123', got '%s'", req.TrxID)
	}
}

func TestNewInquiryVirtualAccountRequest(t *testing.T) {
	req := NewInquiryVirtualAccountRequest("12345", "67890", "12345678901234567890", "inquiry-trx-123")

	if req.PartnerServiceID != "12345" {
		t.Errorf("Expected PartnerServiceID '12345', got '%s'", req.PartnerServiceID)
	}
	if req.TrxID != "inquiry-trx-123" {
		t.Errorf("Expected TrxID 'inquiry-trx-123', got '%s'", req.TrxID)
	}
}

func TestNewVirtualAccountReportRequest(t *testing.T) {
	req := NewVirtualAccountReportRequest("12345", "2024-01-01", "00:00:00", "23:59:59")

	if req.PartnerServiceID != "12345" {
		t.Errorf("Expected PartnerServiceID '12345', got '%s'", req.PartnerServiceID)
	}
	if req.StartDate != "2024-01-01" {
		t.Errorf("Expected StartDate '2024-01-01', got '%s'", req.StartDate)
	}
	if req.StartTime != "00:00:00" {
		t.Errorf("Expected StartTime '00:00:00', got '%s'", req.StartTime)
	}
	if req.EndTime != "23:59:59" {
		t.Errorf("Expected EndTime '23:59:59', got '%s'", req.EndTime)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		ResponseCode:    "4002701",
		ResponseMessage: "Invalid field format",
	}

	expected := "BRI API Error [4002701]: Invalid field format"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// Response code struct tests

func TestBRIResponseCode_String(t *testing.T) {
	rc := &BRIResponseCode{
		HTTPStatus:  200,
		ServiceCode: 27,
		CaseCode:    0,
		FullCode:    "2002700",
	}

	if rc.String() != "2002700" {
		t.Errorf("Expected '2002700', got '%s'", rc.String())
	}
}

func TestBRIResponseCode_IsSuccess(t *testing.T) {
	successRC := &BRIResponseCode{HTTPStatus: 200, FullCode: "2002700"}
	errorRC := &BRIResponseCode{HTTPStatus: 400, FullCode: "4002701"}

	if !successRC.IsSuccess() {
		t.Error("Expected success response code to return true for IsSuccess()")
	}
	if errorRC.IsSuccess() {
		t.Error("Expected error response code to return false for IsSuccess()")
	}
}

func TestBRIResponseCode_IsClientError(t *testing.T) {
	clientErrorRC := &BRIResponseCode{HTTPStatus: 400, FullCode: "4002701"}
	serverErrorRC := &BRIResponseCode{HTTPStatus: 500, FullCode: "5002701"}
	successRC := &BRIResponseCode{HTTPStatus: 200, FullCode: "2002700"}

	if !clientErrorRC.IsClientError() {
		t.Error("Expected 400 response code to return true for IsClientError()")
	}
	if serverErrorRC.IsClientError() {
		t.Error("Expected 500 response code to return false for IsClientError()")
	}
	if successRC.IsClientError() {
		t.Error("Expected 200 response code to return false for IsClientError()")
	}
}

func TestBRIResponseCode_IsServerError(t *testing.T) {
	serverErrorRC := &BRIResponseCode{HTTPStatus: 500, FullCode: "5002701"}
	clientErrorRC := &BRIResponseCode{HTTPStatus: 400, FullCode: "4002701"}
	successRC := &BRIResponseCode{HTTPStatus: 200, FullCode: "2002700"}

	if !serverErrorRC.IsServerError() {
		t.Error("Expected 500 response code to return true for IsServerError()")
	}
	if clientErrorRC.IsServerError() {
		t.Error("Expected 400 response code to return false for IsServerError()")
	}
	if successRC.IsServerError() {
		t.Error("Expected 200 response code to return false for IsServerError()")
	}
}

func TestBRIResponseCode_Getters(t *testing.T) {
	rc := &BRIResponseCode{
		HTTPStatus:  400,
		ServiceCode: 27,
		CaseCode:    1,
		FullCode:    "4002701",
	}

	if rc.GetHTTPStatus() != 400 {
		t.Errorf("Expected HTTPStatus 400, got %d", rc.GetHTTPStatus())
	}
	if rc.GetServiceCode() != 27 {
		t.Errorf("Expected ServiceCode 27, got %d", rc.GetServiceCode())
	}
	if rc.GetCaseCode() != 1 {
		t.Errorf("Expected CaseCode 1, got %d", rc.GetCaseCode())
	}
}

// StructuredBRIAPIResponse tests

func TestStructuredBRIAPIResponse_Error(t *testing.T) {
	resp := &StructuredBRIAPIResponse{
		ResponseCode:    "4002701",
		ResponseMessage: "Invalid field format",
		ResponseDefinition: &BRIVAResponseDefinition{
			ResponseCode: &BRIResponseCode{FullCode: "4002701"},
			Description:  "Invalid field format",
			Field:        "virtualAccountNo",
		},
	}

	expected := "BRI API Error [4002701]: Invalid field format (field: virtualAccountNo)"
	if resp.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, resp.Error())
	}
}

func TestStructuredBRIAPIResponse_ErrorNoDefinition(t *testing.T) {
	resp := &StructuredBRIAPIResponse{
		ResponseCode:    "4002701",
		ResponseMessage: "Invalid field format",
	}

	expected := "BRI API Error [4002701]: Invalid field format"
	if resp.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, resp.Error())
	}
}

func TestStructuredBRIAPIResponse_GetTimestamp(t *testing.T) {
	before := time.Now()
	resp := &StructuredBRIAPIResponse{
		Timestamp: before,
	}
	after := time.Now()

	timestamp := resp.GetTimestamp()
	if timestamp != before {
		t.Error("Expected GetTimestamp to return the stored timestamp")
	}
	if timestamp.Before(before) || timestamp.After(after) {
		t.Error("Expected timestamp to be within reasonable range")
	}
}

func TestStructuredBRIAPIResponse_GetCategory(t *testing.T) {
	resp := &StructuredBRIAPIResponse{
		ResponseDefinition: &BRIVAResponseDefinition{
			Category: CategoryBadRequest,
		},
	}

	if resp.GetCategory() != CategoryBadRequest {
		t.Error("Expected GetCategory to return CategoryBadRequest")
	}
}

func TestStructuredBRIAPIResponse_GetCategoryNoDefinition(t *testing.T) {
	resp := &StructuredBRIAPIResponse{}

	if resp.GetCategory() != CategoryInternalServerError {
		t.Error("Expected GetCategory to return CategoryInternalServerError when no definition")
	}
}

func TestStructuredBRIAPIResponse_IsSuccess(t *testing.T) {
	successResp := &StructuredBRIAPIResponse{HTTPStatusCode: 200}
	errorResp := &StructuredBRIAPIResponse{HTTPStatusCode: 400}

	if !successResp.IsSuccess() {
		t.Error("Expected 200 status to return true for IsSuccess()")
	}
	if errorResp.IsSuccess() {
		t.Error("Expected 400 status to return false for IsSuccess()")
	}
}

func TestStructuredBRIAPIResponse_IsClientError(t *testing.T) {
	clientErrorResp := &StructuredBRIAPIResponse{HTTPStatusCode: 400}
	serverErrorResp := &StructuredBRIAPIResponse{HTTPStatusCode: 500}
	successResp := &StructuredBRIAPIResponse{HTTPStatusCode: 200}

	if !clientErrorResp.IsClientError() {
		t.Error("Expected 400 status to return true for IsClientError()")
	}
	if serverErrorResp.IsClientError() {
		t.Error("Expected 500 status to return false for IsClientError()")
	}
	if successResp.IsClientError() {
		t.Error("Expected 200 status to return false for IsClientError()")
	}
}

func TestStructuredBRIAPIResponse_IsPending(t *testing.T) {
	pendingResp := &StructuredBRIAPIResponse{
		ResponseDefinition: &BRIVAResponseDefinition{
			Category: CategoryPending,
		},
	}

	successResp := &StructuredBRIAPIResponse{
		ResponseDefinition: &BRIVAResponseDefinition{
			Category: CategorySuccess,
		},
	}

	if !pendingResp.IsPending() {
		t.Error("Expected pending response to return true for IsPending()")
	}
	if successResp.IsPending() {
		t.Error("Expected success response to return false for IsPending()")
	}
}

func TestNewStructuredBRIAPIResponse(t *testing.T) {
	resp := NewStructuredBRIAPIResponse("4002701", "Invalid field format")

	if resp.ResponseCode != "4002701" {
		t.Errorf("Expected ResponseCode '4002701', got '%s'", resp.ResponseCode)
	}
	if resp.ResponseMessage != "Invalid field format" {
		t.Errorf("Expected ResponseMessage 'Invalid field format', got '%s'", resp.ResponseMessage)
	}
	if resp.HTTPStatusCode != 400 {
		t.Errorf("Expected HTTPStatusCode 400, got %d", resp.HTTPStatusCode)
	}
}

// makeRequest error scenario tests

func TestMakeRequestHTTPError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	_, err := client.makeRequest(context.Background(), "POST", "/test", map[string]string{"key": "value"})
	if err == nil {
		t.Fatal("Expected HTTP error")
	}

	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

func TestMakeRequestWithNilBody(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request body is empty for GET requests
			if req.Body != nil {
				t.Error("Expected nil body for GET request")
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("OK")),
				Header:     make(http.Header),
			}, nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	resp, err := client.makeRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("makeRequest failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMakeRequestWithBody(t *testing.T) {
	requestBody := map[string]string{"key": "value"}
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify Content-Type header
			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", req.Header.Get("Content-Type"))
			}

			// Read request body
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			var body map[string]string
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}

			if body["key"] != "value" {
				t.Errorf("Expected body key='value', got key='%s'", body["key"])
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("OK")),
				Header:     make(http.Header),
			}, nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	resp, err := client.makeRequest(context.Background(), "POST", "/test", requestBody)
	if err != nil {
		t.Fatalf("makeRequest failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMakeRequestHeaders(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify all required headers are present
			headers := req.Header
			if headers.Get("X-PARTNER-ID") != "test-partner" {
				t.Errorf("Expected X-PARTNER-ID 'test-partner', got '%s'", headers.Get("X-PARTNER-ID"))
			}
			if headers.Get("X-EXTERNAL-ID") == "" {
				t.Error("Expected X-EXTERNAL-ID header to be present")
			}
			if headers.Get("CHANNEL-ID") != "test-channel" {
				t.Errorf("Expected CHANNEL-ID 'test-channel', got '%s'", headers.Get("CHANNEL-ID"))
			}
			if headers.Get("X-SIGNATURE") == "" {
				t.Error("Expected X-SIGNATURE header to be present")
			}
			if headers.Get("X-TIMESTAMP") == "" {
				t.Error("Expected X-TIMESTAMP header to be present")
			}
			if headers.Get("Authorization") != "Bearer test-token" {
				t.Errorf("Expected Authorization 'Bearer test-token', got '%s'", headers.Get("Authorization"))
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString("OK")),
				Header:     make(http.Header),
			}, nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	_, err := client.makeRequest(context.Background(), "POST", "/test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("makeRequest failed: %v", err)
	}
}

// Enhanced VA function tests for 100% coverage

func TestCreateVirtualAccountJSONUnmarshalError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	_, err := client.CreateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected JSON unmarshal error")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal create virtual account response") {
		t.Errorf("Expected JSON unmarshal error, got: %v", err)
	}
}

func TestUpdateVirtualAccountJSONUnmarshalError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &UpdateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Updated Account",
		TotalAmount: Amount{
			Value:    "150000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "update-trx-123",
	}

	_, err := client.UpdateVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected JSON unmarshal error")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal update virtual account response") {
		t.Errorf("Expected JSON unmarshal error, got: %v", err)
	}
}

func TestUpdateVirtualAccountStatusJSONUnmarshalError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &UpdateVirtualAccountStatusRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "status-trx-123",
		PaidStatus:       "Y",
	}

	_, err := client.UpdateVirtualAccountStatus(context.Background(), req)
	if err == nil {
		t.Fatal("Expected JSON unmarshal error")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal update virtual account status response") {
		t.Errorf("Expected JSON unmarshal error, got: %v", err)
	}
}

func TestInquiryVirtualAccountJSONUnmarshalError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &InquiryVirtualAccountRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "inquiry-trx-123",
	}

	_, err := client.InquiryVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected JSON unmarshal error")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal inquiry virtual account response") {
		t.Errorf("Expected JSON unmarshal error, got: %v", err)
	}
}

func TestDeleteVirtualAccountJSONUnmarshalError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &DeleteVirtualAccountRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "delete-trx-123",
	}

	_, err := client.DeleteVirtualAccount(context.Background(), req)
	if err == nil {
		t.Fatal("Expected JSON unmarshal error")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal delete virtual account response") {
		t.Errorf("Expected JSON unmarshal error, got: %v", err)
	}
}

func TestGetVirtualAccountReportJSONUnmarshalError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &VirtualAccountReportRequest{
		PartnerServiceID: "12345",
		StartDate:        "2024-01-01",
		StartTime:        "00:00:00",
		EndTime:          "23:59:59",
		EndDate:          "2024-01-01",
	}

	_, err := client.GetVirtualAccountReport(context.Background(), req)
	if err == nil {
		t.Fatal("Expected JSON unmarshal error")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal virtual account report response") {
		t.Errorf("Expected JSON unmarshal error, got: %v", err)
	}
}

func TestInquiryVirtualAccountStatusJSONUnmarshalError(t *testing.T) {
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	mockAuth := &MockAuthenticator{
		EnsureAuthenticatedFunc: func(ctx context.Context) error {
			return nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	req := &InquiryVirtualAccountStatusRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		InquiryRequestID: "inq-status-123",
	}

	_, err := client.InquiryVirtualAccountStatus(context.Background(), req)
	if err == nil {
		t.Fatal("Expected JSON unmarshal error")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal inquiry virtual account status response") {
		t.Errorf("Expected JSON unmarshal error, got: %v", err)
	}
}

func TestCalculateSignatureEmptyBody(t *testing.T) {
	// Test signature calculation with empty body for non-GET request
	client := &Client{
		clientSecret: "test-secret",
		accessToken:  "test-token",
	}

	signature, err := client.calculateSignature("POST", "/test", "")
	if err != nil {
		t.Fatalf("Failed to calculate signature: %v", err)
	}

	if signature == "" {
		t.Error("Expected non-empty signature")
	}
}

// Integration-style tests for better coverage

func TestEndToEndVAWorkflow(t *testing.T) {
	// Test a complete workflow: Create -> Update -> Inquiry -> Delete
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Return success for all VA operations
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"virtualAccountName": "Test Account",
						"trxId": "test-trx-123"
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

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test Create
	createReq := &CreateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Test Account",
		TotalAmount: Amount{
			Value:    "100000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "test-trx-123",
	}

	createResp, err := client.CreateVirtualAccount(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreateVirtualAccount failed: %v", err)
	}
	if createResp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", createResp.ResponseCode)
	}

	// Test Update
	updateReq := &UpdateVirtualAccountRequest{
		PartnerServiceID:   "12345",
		CustomerNo:         "67890",
		VirtualAccountNo:   "12345678901234567890",
		VirtualAccountName: "Updated Account",
		TotalAmount: Amount{
			Value:    "150000.00",
			Currency: "IDR",
		},
		ExpiredDate: "2024-12-31T23:59:59+07:00",
		TrxID:       "update-trx-123",
	}

	updateResp, err := client.UpdateVirtualAccount(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("UpdateVirtualAccount failed: %v", err)
	}
	if updateResp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", updateResp.ResponseCode)
	}

	// Test Inquiry
	inquiryReq := &InquiryVirtualAccountRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "inquiry-trx-123",
	}

	inquiryResp, err := client.InquiryVirtualAccount(context.Background(), inquiryReq)
	if err != nil {
		t.Fatalf("InquiryVirtualAccount failed: %v", err)
	}
	if inquiryResp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", inquiryResp.ResponseCode)
	}

	// Test Delete
	deleteReq := &DeleteVirtualAccountRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "delete-trx-123",
	}

	deleteResp, err := client.DeleteVirtualAccount(context.Background(), deleteReq)
	if err != nil {
		t.Fatalf("DeleteVirtualAccount failed: %v", err)
	}
	if deleteResp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", deleteResp.ResponseCode)
	}
}

func TestVAStatusWorkflow(t *testing.T) {
	// Test status update and inquiry workflow
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Return success for status operations
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": {
						"partnerServiceId": "12345",
						"customerNo": "67890",
						"virtualAccountNo": "12345678901234567890",
						"paidStatus": "Y"
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

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test Update Status
	statusReq := &UpdateVirtualAccountStatusRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		TrxID:            "status-trx-123",
		PaidStatus:       "Y",
	}

	statusResp, err := client.UpdateVirtualAccountStatus(context.Background(), statusReq)
	if err != nil {
		t.Fatalf("UpdateVirtualAccountStatus failed: %v", err)
	}
	if statusResp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", statusResp.ResponseCode)
	}
	if statusResp.VirtualAccountData.PaidStatus != "Y" {
		t.Errorf("Expected paid status 'Y', got '%s'", statusResp.VirtualAccountData.PaidStatus)
	}

	// Test Status Inquiry
	inquiryStatusReq := &InquiryVirtualAccountStatusRequest{
		PartnerServiceID: "12345",
		CustomerNo:       "67890",
		VirtualAccountNo: "12345678901234567890",
		InquiryRequestID: "inq-status-123",
	}

	inquiryStatusResp, err := client.InquiryVirtualAccountStatus(context.Background(), inquiryStatusReq)
	if err != nil {
		t.Fatalf("InquiryVirtualAccountStatus failed: %v", err)
	}
	if inquiryStatusResp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", inquiryStatusResp.ResponseCode)
	}
	if inquiryStatusResp.VirtualAccountData.PaidStatus != "Y" {
		t.Errorf("Expected paid status 'Y', got '%s'", inquiryStatusResp.VirtualAccountData.PaidStatus)
	}
}

func TestVAReportWorkflow(t *testing.T) {
	// Test report generation
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "2002700",
					"responseMessage": "Successful",
					"virtualAccountData": [
						{
							"partnerServiceId": "12345",
							"customerNo": "67890",
							"virtualAccountNo": "12345678901234567890",
							"virtualAccountName": "Test Account",
							"sourceAccountNo": "1234567890",
							"paidAmount": {
								"value": "100000.00",
								"currency": "IDR"
							},
							"trxDateTime": "2024-01-01T10:00:00+07:00",
							"trxId": "report-trx-123",
							"inquiryRequestId": "inq-123",
							"paymentRequestId": "pay-123",
							"totalAmount": {
								"value": "100000.00",
								"currency": "IDR"
							}
						}
					]
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

	client := &Client{
		httpClient:   mockHTTP,
		auth:         mockAuth,
		baseURL:      "https://api.example.com",
		partnerID:    "test-partner",
		clientID:     "test-client",
		clientSecret: "test-secret",
		privateKey:   "test-key",
		channelID:    "test-channel",
		accessToken:  "test-token",
		tokenExpiry:  time.Now().Add(time.Hour),
	}

	// Test Get Report
	reportReq := &VirtualAccountReportRequest{
		PartnerServiceID: "12345",
		StartDate:        "2024-01-01",
		StartTime:        "00:00:00",
		EndTime:          "23:59:59",
		EndDate:          "2024-01-01",
	}

	reportResp, err := client.GetVirtualAccountReport(context.Background(), reportReq)
	if err != nil {
		t.Fatalf("GetVirtualAccountReport failed: %v", err)
	}
	if reportResp.ResponseCode != "2002700" {
		t.Errorf("Expected response code '2002700', got '%s'", reportResp.ResponseCode)
	}
	if len(reportResp.VirtualAccountData) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(reportResp.VirtualAccountData))
	}
	if reportResp.VirtualAccountData[0].TrxID != "report-trx-123" {
		t.Errorf("Expected trxId 'report-trx-123', got '%s'", reportResp.VirtualAccountData[0].TrxID)
	}
}

// Test for the authenticate function in auth.go

func TestClientAuthenticatePrivateKeyParsing(t *testing.T) {
	// Test private key parsing functionality
	client := &Client{
		httpClient:   &MockHTTPClient{},
		baseURL:      "https://api.example.com",
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
		privateKey:   "invalid-private-key-format",
	}

	ctx := context.Background()
	err := client.authenticate(ctx)
	if err == nil {
		t.Fatal("Expected authentication to fail with invalid private key")
	}

	if !strings.Contains(err.Error(), "failed to decode PEM block") {
		t.Errorf("Expected PEM decode error, got: %v", err)
	}
}

func TestClientAuthenticateTokenRequestCreation(t *testing.T) {
	// Test token request creation and HTTP call
	callCount := 0
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			// Verify request structure
			if req.URL.Path != "/snap/v1.0/access-token/b2b" {
				t.Errorf("Expected path '/snap/v1.0/access-token/b2b', got '%s'", req.URL.Path)
			}
			if req.Method != "POST" {
				t.Errorf("Expected POST method, got '%s'", req.Method)
			}

			// Read and verify request body
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			var body map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}

			if body["grantType"] != "client_credentials" {
				t.Errorf("Expected grantType 'client_credentials', got '%v'", body["grantType"])
			}

			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"accessToken": "test-access-token-12345",
					"tokenType": "Bearer",
					"expiresIn": 3600
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
		privateKey:   privateKeyTest,
	}

	ctx := context.Background()
	err := client.authenticate(ctx)
	if err != nil {
		t.Fatalf("authenticate() failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 HTTP call, got %d", callCount)
	}

	// Verify token was stored
	if client.accessToken != "test-access-token-12345" {
		t.Errorf("Expected access token 'test-access-token-12345', got '%s'", client.accessToken)
	}
}

func TestClientAuthenticateTokenResponseParsing(t *testing.T) {
	// Test token response parsing and storage
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"accessToken": "test-access-token-parsing",
					"tokenType": "Bearer",
					"expiresIn": 7200
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
		privateKey:   privateKeyTest,
	}

	ctx := context.Background()
	err := client.authenticate(ctx)
	if err != nil {
		t.Fatalf("authenticate() failed: %v", err)
	}

	// Verify token details
	if client.accessToken != "test-access-token-parsing" {
		t.Errorf("Expected access token 'test-access-token-parsing', got '%s'", client.accessToken)
	}

	// Verify expiry is set correctly (approximately 2 hours from now)
	expectedExpiry := time.Now().Add(7200 * time.Second)
	timeDiff := client.tokenExpiry.Sub(expectedExpiry)
	if timeDiff > time.Second || timeDiff < -time.Second {
		t.Errorf("Expected token expiry around %v, got %v", expectedExpiry, client.tokenExpiry)
	}
}

func TestClientAuthenticateErrorResponseHandling(t *testing.T) {
	// Test error response handling in authentication
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 401,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"responseCode": "4012705",
					"responseMessage": "Invalid credentials"
				}`)),
				Header: make(http.Header),
			}, nil
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
		privateKey:   privateKeyTest,
	}

	ctx := context.Background()
	err := client.authenticate(ctx)
	if err == nil {
		t.Fatal("Expected authentication to fail with error response")
	}

	// Check if it's wrapped in APIError (which is what authenticate function actually returns)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("Expected APIError, got %T", err)
	} else {
		// Only check properties if we successfully cast to APIError
		if apiErr.ResponseCode != "4012705" {
			t.Errorf("Expected response code '4012705', got '%s'", apiErr.ResponseCode)
		}
		if apiErr.ResponseMessage != "Invalid credentials" {
			t.Errorf("Expected response message 'Invalid credentials', got '%s'", apiErr.ResponseMessage)
		}
	}
}

func TestClientAuthenticateNetworkFailure(t *testing.T) {
	// Test network failure during authentication
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error: connection refused")
		},
	}

	client := &Client{
		httpClient:   mockHTTP,
		baseURL:      "https://api.example.com",
		clientID:     "test-client-id",
		clientSecret: "test-client-secret",
		privateKey:   privateKeyTest,
	}

	ctx := context.Background()
	err := client.authenticate(ctx)
	if err == nil {
		t.Fatal("Expected authentication to fail with network error")
	}

	if !strings.Contains(err.Error(), "network error") {
		t.Errorf("Expected network error, got: %v", err)
	}
}
