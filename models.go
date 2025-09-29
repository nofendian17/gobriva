package gobriva

import "fmt"

// Amount represents monetary amount with currency
type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// AdditionalInfo represents additional // NewCreateVirtualAccountRequest creates a new CreateVirtualAccountRequest with default values
func NewCreateVirtualAccountRequest(partnerServiceID, customerNo, vaNo, vaName, trxID string, amount float64, currency, expiredDate string) *CreateVirtualAccountRequest {
	return &CreateVirtualAccountRequest{
		PartnerServiceID:   partnerServiceID,
		CustomerNo:         customerNo,
		VirtualAccountNo:   vaNo,
		VirtualAccountName: vaName,
		TotalAmount: Amount{
			Value:    fmt.Sprintf("%.2f", amount),
			Currency: currency,
		},
		ExpiredDate: expiredDate,
		TrxID:       trxID,
	}
}

// AdditionalInfo represents additional information for VA
type AdditionalInfo struct {
	Description string `json:"description,omitempty"`
}

// CreateVirtualAccountRequest represents the request to create a virtual account
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

// CreateVirtualAccountResponse represents the response from creating a virtual account
type CreateVirtualAccountResponse struct {
	ResponseCode       string              `json:"responseCode"`
	ResponseMessage    string              `json:"responseMessage"`
	VirtualAccountData *VirtualAccountData `json:"virtualAccountData,omitempty"`
}

// UpdateVirtualAccountRequest represents the request to update a virtual account
type UpdateVirtualAccountRequest struct {
	PartnerServiceID   string         `json:"partnerServiceId"`
	CustomerNo         string         `json:"customerNo"`
	VirtualAccountNo   string         `json:"virtualAccountNo"`
	VirtualAccountName string         `json:"virtualAccountName"`
	TotalAmount        Amount         `json:"totalAmount"`
	ExpiredDate        string         `json:"expiredDate"`
	TrxID              string         `json:"trxId"`
	AdditionalInfo     AdditionalInfo `json:"additionalInfo,omitempty"`
}

// UpdateVirtualAccountResponse represents the response from updating a virtual account
type UpdateVirtualAccountResponse struct {
	ResponseCode       string              `json:"responseCode"`
	ResponseMessage    string              `json:"responseMessage"`
	VirtualAccountData *VirtualAccountData `json:"virtualAccountData,omitempty"`
}

// UpdateVirtualAccountStatusRequest represents the request to update VA status
type UpdateVirtualAccountStatusRequest struct {
	PartnerServiceID string `json:"partnerServiceId"`
	CustomerNo       string `json:"customerNo"`
	VirtualAccountNo string `json:"virtualAccountNo"`
	TrxID            string `json:"trxId"`
	PaidStatus       string `json:"paidStatus"`
}

// UpdateVirtualAccountStatusResponse represents the response from updating VA status
type UpdateVirtualAccountStatusResponse struct {
	ResponseCode       string              `json:"responseCode"`
	ResponseMessage    string              `json:"responseMessage"`
	VirtualAccountData *VirtualAccountData `json:"virtualAccountData,omitempty"`
}

// InquiryVirtualAccountStatusRequest represents the request for VA status inquiry
type InquiryVirtualAccountStatusRequest struct {
	PartnerServiceID string `json:"partnerServiceId"`
	CustomerNo       string `json:"customerNo"`
	VirtualAccountNo string `json:"virtualAccountNo"`
	InquiryRequestID string `json:"inquiryRequestId"`
}

// InquiryVirtualAccountStatusResponse represents the response from VA status inquiry
type InquiryVirtualAccountStatusResponse struct {
	ResponseCode       string              `json:"responseCode"`
	ResponseMessage    string              `json:"responseMessage"`
	VirtualAccountData *VirtualAccountData `json:"virtualAccountData,omitempty"`
	AdditionalInfo     AdditionalInfo      `json:"additionalInfo,omitempty"`
}

// InquiryVirtualAccountRequest represents the request for VA inquiry
type InquiryVirtualAccountRequest struct {
	PartnerServiceID string `json:"partnerServiceId"`
	CustomerNo       string `json:"customerNo"`
	VirtualAccountNo string `json:"virtualAccountNo"`
	TrxID            string `json:"trxId"`
}

// InquiryVirtualAccountResponse represents the response from VA inquiry
type InquiryVirtualAccountResponse struct {
	ResponseCode       string              `json:"responseCode"`
	ResponseMessage    string              `json:"responseMessage"`
	VirtualAccountData *VirtualAccountData `json:"virtualAccountData,omitempty"`
}

// DeleteVirtualAccountRequest represents the request to delete a virtual account
type DeleteVirtualAccountRequest struct {
	PartnerServiceID string `json:"partnerServiceId"`
	CustomerNo       string `json:"customerNo"`
	VirtualAccountNo string `json:"virtualAccountNo"`
	TrxID            string `json:"trxId"`
}

// DeleteVirtualAccountResponse represents the response from deleting a virtual account
type DeleteVirtualAccountResponse struct {
	ResponseCode       string              `json:"responseCode"`
	ResponseMessage    string              `json:"responseMessage"`
	VirtualAccountData *VirtualAccountData `json:"virtualAccountData,omitempty"`
}

// VirtualAccountReportRequest represents the request for VA report
type VirtualAccountReportRequest struct {
	PartnerServiceID string `json:"partnerServiceId"`
	StartDate        string `json:"startDate"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	EndDate          string `json:"endDate,omitempty"`
}

// VirtualAccountReportResponse represents the response from VA report
type VirtualAccountReportResponse struct {
	ResponseCode       string                      `json:"responseCode"`
	ResponseMessage    string                      `json:"responseMessage"`
	VirtualAccountData []VirtualAccountTransaction `json:"virtualAccountData,omitempty"`
}

// VirtualAccountData represents virtual account information
type VirtualAccountData struct {
	InstitutionCode    string         `json:"institutionCode,omitempty"`
	PartnerServiceID   string         `json:"partnerServiceId"`
	CustomerNo         string         `json:"customerNo"`
	VirtualAccountNo   string         `json:"virtualAccountNo"`
	VirtualAccountName string         `json:"virtualAccountName"`
	TrxID              string         `json:"trxId"`
	TotalAmount        Amount         `json:"totalAmount,omitempty"`
	ExpiredDate        string         `json:"expiredDate,omitempty"`
	AdditionalInfo     AdditionalInfo `json:"additionalInfo,omitempty"`
	PaidStatus         string         `json:"paidStatus,omitempty"`
}

// VirtualAccountTransaction represents a transaction in VA report
type VirtualAccountTransaction struct {
	PartnerServiceID   string     `json:"partnerServiceId"`
	CustomerNo         string     `json:"customerNo"`
	VirtualAccountNo   string     `json:"virtualAccountNo"`
	VirtualAccountName string     `json:"virtualAccountName"`
	SourceAccountNo    string     `json:"sourceAccountNo"`
	PaidAmount         Amount     `json:"paidAmount"`
	TrxDateTime        string     `json:"trxDateTime"`
	TrxID              string     `json:"trxId"`
	InquiryRequestID   string     `json:"inquiryRequestId"`
	PaymentRequestID   string     `json:"paymentRequestId"`
	TotalAmount        Amount     `json:"totalAmount"`
	FreeTexts          []FreeText `json:"freeTexts,omitempty"`
}

// FreeText represents free text information in multiple languages
type FreeText struct {
	English   string `json:"english"`
	Indonesia string `json:"indonesia"`
}

// TokenRequest represents the OAuth2 token request
type TokenRequest struct {
	GrantType string `json:"grantType"`
}

// APIError represents an error from the BRI API
type APIError struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("BRI API Error [%s]: %s", e.ResponseCode, e.ResponseMessage)
}

// Helper functions for creating requests

// NewUpdateVirtualAccountRequest creates a new UpdateVirtualAccountRequest with default values
func NewUpdateVirtualAccountRequest(partnerServiceID, customerNo, vaNo, vaName, trxID string, amount float64, currency, expiredDate string) *UpdateVirtualAccountRequest {
	return &UpdateVirtualAccountRequest{
		PartnerServiceID:   partnerServiceID,
		CustomerNo:         customerNo,
		VirtualAccountNo:   vaNo,
		VirtualAccountName: vaName,
		TotalAmount: Amount{
			Value:    fmt.Sprintf("%.2f", amount),
			Currency: currency,
		},
		ExpiredDate: expiredDate,
		TrxID:       trxID,
	}
}

// NewUpdateVirtualAccountStatusRequest creates a new UpdateVirtualAccountStatusRequest
func NewUpdateVirtualAccountStatusRequest(partnerServiceID, customerNo, vaNo, trxID, paidStatus string) *UpdateVirtualAccountStatusRequest {
	return &UpdateVirtualAccountStatusRequest{
		PartnerServiceID: partnerServiceID,
		CustomerNo:       customerNo,
		VirtualAccountNo: vaNo,
		TrxID:            trxID,
		PaidStatus:       paidStatus,
	}
}

// NewInquiryVirtualAccountRequest creates a new InquiryVirtualAccountRequest
func NewInquiryVirtualAccountRequest(partnerServiceID, customerNo, vaNo, trxID string) *InquiryVirtualAccountRequest {
	return &InquiryVirtualAccountRequest{
		PartnerServiceID: partnerServiceID,
		CustomerNo:       customerNo,
		VirtualAccountNo: vaNo,
		TrxID:            trxID,
	}
}

// NewVirtualAccountReportRequest creates a new VirtualAccountReportRequest
func NewVirtualAccountReportRequest(partnerServiceID, startDate, startTime, endTime string) *VirtualAccountReportRequest {
	return &VirtualAccountReportRequest{
		PartnerServiceID: partnerServiceID,
		StartDate:        startDate,
		StartTime:        startTime,
		EndTime:          endTime,
	}
}
