package gobriva

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CreateVirtualAccount creates a new virtual account
func (c *Client) CreateVirtualAccount(ctx context.Context, req *CreateVirtualAccountRequest) (*CreateVirtualAccountResponse, error) {
	// Ensure authentication
	if err := c.auth.EnsureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Make request
	resp, err := c.makeRequest(ctx, "POST", "/snap/v1.0/transfer-va/create-va", req)
	if err != nil {
		return nil, fmt.Errorf("failed to make create virtual account request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read create virtual account response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return nil, &StructuredBRIAPIResponse{
			ResponseCode:       errorResp.ResponseCode,
			ResponseMessage:    errorResp.ResponseMessage,
			ResponseDefinition: GetBRIVAResponseDefinition(errorResp.ResponseCode),
			HTTPStatusCode:     resp.StatusCode,
			Timestamp:          time.Now(),
		}
	}

	var createResp CreateVirtualAccountResponse
	if err := json.Unmarshal(respBody, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create virtual account response: %w", err)
	}

	return &createResp, nil
}

// UpdateVirtualAccount updates an existing virtual account
func (c *Client) UpdateVirtualAccount(ctx context.Context, req *UpdateVirtualAccountRequest) (*UpdateVirtualAccountResponse, error) {
	// Ensure authentication
	if err := c.auth.EnsureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Make request
	resp, err := c.makeRequest(ctx, "PUT", "/snap/v1.0/transfer-va/update-va", req)
	if err != nil {
		return nil, fmt.Errorf("failed to make update virtual account request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read update virtual account response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return nil, &StructuredBRIAPIResponse{
			ResponseCode:       errorResp.ResponseCode,
			ResponseMessage:    errorResp.ResponseMessage,
			ResponseDefinition: GetBRIVAResponseDefinition(errorResp.ResponseCode),
			HTTPStatusCode:     resp.StatusCode,
			Timestamp:          time.Now(),
		}
	}

	var updateResp UpdateVirtualAccountResponse
	if err := json.Unmarshal(respBody, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update virtual account response: %w", err)
	}

	return &updateResp, nil
}

// UpdateVirtualAccountStatus updates the status of a virtual account
func (c *Client) UpdateVirtualAccountStatus(ctx context.Context, req *UpdateVirtualAccountStatusRequest) (*UpdateVirtualAccountStatusResponse, error) {
	// Ensure authentication
	if err := c.auth.EnsureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Make request
	resp, err := c.makeRequest(ctx, "PUT", "/snap/v1.0/transfer-va/update-status", req)
	if err != nil {
		return nil, fmt.Errorf("failed to make update virtual account status request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read update virtual account status response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return nil, &StructuredBRIAPIResponse{
			ResponseCode:       errorResp.ResponseCode,
			ResponseMessage:    errorResp.ResponseMessage,
			ResponseDefinition: GetBRIVAResponseDefinition(errorResp.ResponseCode),
			HTTPStatusCode:     resp.StatusCode,
			Timestamp:          time.Now(),
		}
	}

	var statusResp UpdateVirtualAccountStatusResponse
	if err := json.Unmarshal(respBody, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update virtual account status response: %w", err)
	}

	return &statusResp, nil
}

// InquiryVirtualAccount gets information about a virtual account
func (c *Client) InquiryVirtualAccount(ctx context.Context, req *InquiryVirtualAccountRequest) (*InquiryVirtualAccountResponse, error) {
	// Ensure authentication
	if err := c.auth.EnsureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Make request
	resp, err := c.makeRequest(ctx, "POST", "/snap/v1.0/transfer-va/inquiry-va", req)
	if err != nil {
		return nil, fmt.Errorf("failed to make inquiry virtual account request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read inquiry virtual account response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return nil, &StructuredBRIAPIResponse{
			ResponseCode:       errorResp.ResponseCode,
			ResponseMessage:    errorResp.ResponseMessage,
			ResponseDefinition: GetBRIVAResponseDefinition(errorResp.ResponseCode),
			HTTPStatusCode:     resp.StatusCode,
			Timestamp:          time.Now(),
		}
	}

	var inquiryResp InquiryVirtualAccountResponse
	if err := json.Unmarshal(respBody, &inquiryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inquiry virtual account response: %w", err)
	}

	return &inquiryResp, nil
}

// DeleteVirtualAccount deletes a virtual account
func (c *Client) DeleteVirtualAccount(ctx context.Context, req *DeleteVirtualAccountRequest) (*DeleteVirtualAccountResponse, error) {
	// Ensure authentication
	if err := c.auth.EnsureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Make request
	resp, err := c.makeRequest(ctx, "DELETE", "/snap/v1.0/transfer-va/delete-va", req)
	if err != nil {
		return nil, fmt.Errorf("failed to make delete virtual account request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read delete virtual account response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return nil, &StructuredBRIAPIResponse{
			ResponseCode:       errorResp.ResponseCode,
			ResponseMessage:    errorResp.ResponseMessage,
			ResponseDefinition: GetBRIVAResponseDefinition(errorResp.ResponseCode),
			HTTPStatusCode:     resp.StatusCode,
			Timestamp:          time.Now(),
		}
	}

	var deleteResp DeleteVirtualAccountResponse
	if err := json.Unmarshal(respBody, &deleteResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal delete virtual account response: %w", err)
	}

	return &deleteResp, nil
}

// GetVirtualAccountReport gets a report of virtual account transactions
func (c *Client) GetVirtualAccountReport(ctx context.Context, req *VirtualAccountReportRequest) (*VirtualAccountReportResponse, error) {
	// Ensure authentication
	if err := c.auth.EnsureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Make request
	resp, err := c.makeRequest(ctx, "POST", "/snap/v1.0/transfer-va/report", req)
	if err != nil {
		return nil, fmt.Errorf("failed to make virtual account report request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read virtual account report response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return nil, &StructuredBRIAPIResponse{
			ResponseCode:       errorResp.ResponseCode,
			ResponseMessage:    errorResp.ResponseMessage,
			ResponseDefinition: GetBRIVAResponseDefinition(errorResp.ResponseCode),
			HTTPStatusCode:     resp.StatusCode,
			Timestamp:          time.Now(),
		}
	}

	var reportResp VirtualAccountReportResponse
	if err := json.Unmarshal(respBody, &reportResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal virtual account report response: %w", err)
	}

	return &reportResp, nil
}

// InquiryVirtualAccountStatus inquires the status of a virtual account
func (c *Client) InquiryVirtualAccountStatus(ctx context.Context, req *InquiryVirtualAccountStatusRequest) (*InquiryVirtualAccountStatusResponse, error) {
	// Ensure authentication
	if err := c.auth.EnsureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Make request
	resp, err := c.makeRequest(ctx, "POST", "/snap/v1.0/transfer-va/status", req)
	if err != nil {
		return nil, fmt.Errorf("failed to make inquiry virtual account status request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read inquiry virtual account status response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return nil, &StructuredBRIAPIResponse{
			ResponseCode:       errorResp.ResponseCode,
			ResponseMessage:    errorResp.ResponseMessage,
			ResponseDefinition: GetBRIVAResponseDefinition(errorResp.ResponseCode),
			HTTPStatusCode:     resp.StatusCode,
			Timestamp:          time.Now(),
		}
	}

	var inquiryResp InquiryVirtualAccountStatusResponse
	if err := json.Unmarshal(respBody, &inquiryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inquiry virtual account status response: %w", err)
	}

	return &inquiryResp, nil
}
