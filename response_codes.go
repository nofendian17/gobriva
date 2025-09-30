package gobriva

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// HttpCategory represents the category of an error
type HttpCategory string

const (
	CategorySuccess             HttpCategory = "Success"
	CategoryBadRequest          HttpCategory = "BadRequest"
	CategoryUnauthorized        HttpCategory = "Unauthorized"
	CategoryForbidden           HttpCategory = "Forbidden"
	CategoryNotFound            HttpCategory = "NotFound"
	CategoryMethodNotAllowed    HttpCategory = "MethodNotAllowed"
	CategoryConflict            HttpCategory = "Conflict"
	CategoryInternalServerError HttpCategory = "InternalServerError"
	CategoryBadGateway          HttpCategory = "BadGateway"
	CategoryServiceUnavailable  HttpCategory = "ServiceUnavailable"
	CategoryPending             HttpCategory = "Pending"
)

// BRIResponseCode represents a structured BRI API response code
// Format: HTTPSTATUS(3) + SERVICECODE(2) + CASECODE(2)
// Example: "2002700" = HTTP 200 + Service 27 + Case 00
type BRIResponseCode struct {
	HTTPStatus  int    // HTTP status code (200, 400, 401, etc.)
	ServiceCode int    // Service identifier (27 for BRIVA)
	CaseCode    int    // Specific case within the service
	FullCode    string // Complete 7-digit response code
}

// String returns the full 7-digit response code
func (rc *BRIResponseCode) String() string {
	return rc.FullCode
}

// IsSuccess checks if this is a success response code
func (rc *BRIResponseCode) IsSuccess() bool {
	return rc.HTTPStatus >= 200 && rc.HTTPStatus < 300
}

// IsClientError checks if this is a client error (4xx)
func (rc *BRIResponseCode) IsClientError() bool {
	return rc.HTTPStatus >= 400 && rc.HTTPStatus < 500
}

// IsServerError checks if this is a server error (5xx)
func (rc *BRIResponseCode) IsServerError() bool {
	return rc.HTTPStatus >= 500
}

// GetHTTPStatus returns the HTTP status code
func (rc *BRIResponseCode) GetHTTPStatus() int {
	return rc.HTTPStatus
}

// GetServiceCode returns the service code
func (rc *BRIResponseCode) GetServiceCode() int {
	return rc.ServiceCode
}

// GetCaseCode returns the case code
func (rc *BRIResponseCode) GetCaseCode() int {
	return rc.CaseCode
}

// BRIVAResponseDefinition contains detailed information about a BRIVA response code
type BRIVAResponseDefinition struct {
	ResponseCode *BRIResponseCode
	Category     HttpCategory
	Description  string
	Field        string // Specific field that caused the error (if applicable)
}

// BRIVA Response Code Definitions
var brivaResponseDefinitions = map[string]*BRIVAResponseDefinition{
	// Success Codes
	"2002600": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 26, CaseCode: 0, FullCode: "2002600"},
		Category:     CategorySuccess,
		Description:  "Inquiry status successful",
	},
	"2002700": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 27, CaseCode: 0, FullCode: "2002700"},
		Category:     CategorySuccess,
		Description:  "Request processed successfully",
	},
	"2002701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 27, CaseCode: 1, FullCode: "2002701"},
		Category:     CategorySuccess,
		Description:  "Virtual Account created successfully",
	},
	"2002800": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 28, CaseCode: 0, FullCode: "2002800"},
		Category:     CategorySuccess,
		Description:  "Virtual Account updated successfully",
	},
	"2002900": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 29, CaseCode: 0, FullCode: "2002900"},
		Category:     CategorySuccess,
		Description:  "Virtual Account status updated successfully",
	},
	"2003000": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 30, CaseCode: 0, FullCode: "2003000"},
		Category:     CategorySuccess,
		Description:  "Virtual Account inquiry successful",
	},
	"2003100": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 31, CaseCode: 0, FullCode: "2003100"},
		Category:     CategorySuccess,
		Description:  "Virtual Account deleted successfully",
	},
	"2003500": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 200, ServiceCode: 35, CaseCode: 0, FullCode: "2003500"},
		Category:     CategorySuccess,
		Description:  "Report generated successfully",
	},

	// Bad Request Codes (400xxxx)
	"4002701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 1, FullCode: "4002701"},
		Category:     CategoryBadRequest,
		Description:  "Invalid field format",
		Field:        "virtualAccountNo",
	},
	"4002702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 2, FullCode: "4002702"},
		Category:     CategoryBadRequest,
		Description:  "Invalid mandatory field",
		Field:        "partnerServiceId",
	},
	"4002703": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 3, FullCode: "4002703"},
		Category:     CategoryBadRequest,
		Description:  "Invalid field value",
	},
	"4002704": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 4, FullCode: "4002704"},
		Category:     CategoryBadRequest,
		Description:  "Invalid amount format or value",
		Field:        "totalAmount",
	},
	"4002705": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 5, FullCode: "4002705"},
		Category:     CategoryBadRequest,
		Description:  "Invalid account information",
	},
	"4002706": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 6, FullCode: "4002706"},
		Category:     CategoryBadRequest,
		Description:  "Invalid date format",
		Field:        "expiredDate",
	},
	"4002707": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 7, FullCode: "4002707"},
		Category:     CategoryBadRequest,
		Description:  "Invalid time format",
	},
	"4002708": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 8, FullCode: "4002708"},
		Category:     CategoryBadRequest,
		Description:  "Invalid currency code",

		Field: "currency",
	},
	"4002709": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 9, FullCode: "4002709"},
		Category:     CategoryBadRequest,
		Description:  "Invalid partner service ID",

		Field: "partnerServiceId",
	},
	"4002710": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 10, FullCode: "4002710"},
		Category:     CategoryBadRequest,
		Description:  "Invalid customer number",

		Field: "customerNo",
	},
	"4002711": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 11, FullCode: "4002711"},
		Category:     CategoryBadRequest,
		Description:  "Invalid virtual account number",

		Field: "virtualAccountNo",
	},
	"4002712": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 12, FullCode: "4002712"},
		Category:     CategoryBadRequest,
		Description:  "Invalid virtual account name",

		Field: "virtualAccountName",
	},
	"4002713": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 13, FullCode: "4002713"},
		Category:     CategoryBadRequest,
		Description:  "Invalid transaction ID",

		Field: "trxId",
	},
	"4002714": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 14, FullCode: "4002714"},
		Category:     CategoryBadRequest,
		Description:  "Invalid paid status",

		Field: "paidStatus",
	},
	"4002715": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 15, FullCode: "4002715"},
		Category:     CategoryBadRequest,
		Description:  "Invalid inquiry request ID",

		Field: "inquiryRequestId",
	},
	"4002716": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 16, FullCode: "4002716"},
		Category:     CategoryBadRequest,
		Description:  "Invalid report date range",

		Field: "startDate",
	},
	"4002717": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 27, CaseCode: 17, FullCode: "4002717"},
		Category:     CategoryBadRequest,
		Description:  "Invalid report time range",

		Field: "startTime",
	},
	"4002600": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 26, CaseCode: 0, FullCode: "4002600"},
		Category:     CategoryBadRequest,
		Description:  "Bad Request",
	},
	"4002601": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 26, CaseCode: 1, FullCode: "4002601"},
		Category:     CategoryBadRequest,
		Description:  "Invalid Field Format",
	},
	"4002602": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 400, ServiceCode: 26, CaseCode: 2, FullCode: "4002602"},
		Category:     CategoryBadRequest,
		Description:  "Invalid Mandatory Field",
	},

	// Unauthorized Codes (401xxxx)
	"4012701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 27, CaseCode: 1, FullCode: "4012701"},
		Category:     CategoryUnauthorized,
		Description:  "Invalid signature",
	},
	"4012702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 27, CaseCode: 2, FullCode: "4012702"},
		Category:     CategoryUnauthorized,
		Description:  "Invalid timestamp",
	},
	"4012703": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 27, CaseCode: 3, FullCode: "4012703"},
		Category:     CategoryUnauthorized,
		Description:  "Invalid access token",
	},
	"4012704": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 27, CaseCode: 4, FullCode: "4012704"},
		Category:     CategoryUnauthorized,
		Description:  "Access token expired",
	},
	"4012705": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 27, CaseCode: 5, FullCode: "4012705"},
		Category:     CategoryUnauthorized,
		Description:  "Invalid credentials",
	},
	"4012706": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 27, CaseCode: 6, FullCode: "4012706"},
		Category:     CategoryUnauthorized,
		Description:  "Invalid client key",
	},
	"4012707": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 27, CaseCode: 7, FullCode: "4012707"},
		Category:     CategoryUnauthorized,
		Description:  "Invalid private key",
	},
	"4012600": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 401, ServiceCode: 26, CaseCode: 0, FullCode: "4012600"},
		Category:     CategoryUnauthorized,
		Description:  "Unauthorized. Client Forbidden Access API",
	},

	// Forbidden Codes (403xxxx)
	"4032701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 403, ServiceCode: 27, CaseCode: 1, FullCode: "4032701"},
		Category:     CategoryForbidden,
		Description:  "Insufficient permission",
	},
	"4032702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 403, ServiceCode: 27, CaseCode: 2, FullCode: "4032702"},
		Category:     CategoryForbidden,
		Description:  "Access denied",
	},
	"4032703": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 403, ServiceCode: 27, CaseCode: 3, FullCode: "4032703"},
		Category:     CategoryForbidden,
		Description:  "Partner not active",
	},
	"4032704": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 403, ServiceCode: 27, CaseCode: 4, FullCode: "4032704"},
		Category:     CategoryForbidden,
		Description:  "Channel not allowed",
	},
	"4032705": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 403, ServiceCode: 27, CaseCode: 5, FullCode: "4032705"},
		Category:     CategoryForbidden,
		Description:  "IP not whitelisted",
	},

	// Not Found Codes (404xxxx)
	"4042701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 404, ServiceCode: 27, CaseCode: 1, FullCode: "4042701"},
		Category:     CategoryNotFound,
		Description:  "Virtual Account not found",
	},
	"4042702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 404, ServiceCode: 27, CaseCode: 2, FullCode: "4042702"},
		Category:     CategoryNotFound,
		Description:  "Customer not found",
	},
	"4042703": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 404, ServiceCode: 27, CaseCode: 3, FullCode: "4042703"},
		Category:     CategoryNotFound,
		Description:  "Partner service not found",
	},
	"4042704": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 404, ServiceCode: 27, CaseCode: 4, FullCode: "4042704"},
		Category:     CategoryNotFound,
		Description:  "Transaction not found",
	},
	"4042612": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 404, ServiceCode: 26, CaseCode: 12, FullCode: "4042612"},
		Category:     CategoryNotFound,
		Description:  "Invalid Bill/Virtual Account",
	},
	"4042613": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 404, ServiceCode: 26, CaseCode: 13, FullCode: "4042613"},
		Category:     CategoryNotFound,
		Description:  "Invalid Amount",
	},

	// Method Not Allowed Codes (405xxxx)
	"4052701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 405, ServiceCode: 27, CaseCode: 1, FullCode: "4052701"},
		Category:     CategoryMethodNotAllowed,
		Description:  "HTTP method not allowed",
	},
	"4052702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 405, ServiceCode: 27, CaseCode: 2, FullCode: "4052702"},
		Category:     CategoryMethodNotAllowed,
		Description:  "HTTP method not allowed for this endpoint",
	},

	// Conflict Codes (409xxxx)
	"4092701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 409, ServiceCode: 27, CaseCode: 1, FullCode: "4092701"},
		Category:     CategoryConflict,
		Description:  "Virtual Account already exists",
	},
	"4092702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 409, ServiceCode: 27, CaseCode: 2, FullCode: "4092702"},
		Category:     CategoryConflict,
		Description:  "Virtual Account number already exists",
	},
	"4092703": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 409, ServiceCode: 27, CaseCode: 3, FullCode: "4092703"},
		Category:     CategoryConflict,
		Description:  "Transaction ID already exists",
	},
	"4092704": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 409, ServiceCode: 27, CaseCode: 4, FullCode: "4092704"},
		Category:     CategoryConflict,
		Description:  "Customer number already exists",
	},
	"4092601": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 409, ServiceCode: 26, CaseCode: 1, FullCode: "4092601"},
		Category:     CategoryConflict,
		Description:  "Conflict",
	},

	// Internal Server Error Codes (500xxxx)
	"5002701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 500, ServiceCode: 27, CaseCode: 1, FullCode: "5002701"},
		Category:     CategoryInternalServerError,
		Description:  "Internal server error",
	},
	"5002702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 500, ServiceCode: 27, CaseCode: 2, FullCode: "5002702"},
		Category:     CategoryInternalServerError,
		Description:  "Database error",
	},
	"5002703": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 500, ServiceCode: 27, CaseCode: 3, FullCode: "5002703"},
		Category:     CategoryInternalServerError,
		Description:  "External service error",
	},
	"5002704": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 500, ServiceCode: 27, CaseCode: 4, FullCode: "5002704"},
		Category:     CategoryInternalServerError,
		Description:  "System under maintenance",
	},
	"5002705": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 500, ServiceCode: 27, CaseCode: 5, FullCode: "5002705"},
		Category:     CategoryInternalServerError,
		Description:  "System unavailable",
	},
	"5002600": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 500, ServiceCode: 26, CaseCode: 0, FullCode: "5002600"},
		Category:     CategoryInternalServerError,
		Description:  "General Error",
	},

	// Bad Gateway Codes (502xxxx)
	"5022701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 502, ServiceCode: 27, CaseCode: 1, FullCode: "5022701"},
		Category:     CategoryBadGateway,
		Description:  "Bad gateway",
	},
	"5022702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 502, ServiceCode: 27, CaseCode: 2, FullCode: "5022702"},
		Category:     CategoryBadGateway,
		Description:  "External service timeout",
	},

	// Timeout Codes (504xxxx)
	"5042700": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 504, ServiceCode: 27, CaseCode: 0, FullCode: "5042700"},
		Category:     CategoryServiceUnavailable, // Assuming timeout is service unavailable
		Description:  "Timeout",
	},
	"5042600": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 504, ServiceCode: 26, CaseCode: 0, FullCode: "5042600"},
		Category:     CategoryServiceUnavailable,
		Description:  "Timeout",
	},

	// Service Unavailable Codes (503xxxx)
	"5032701": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 503, ServiceCode: 27, CaseCode: 1, FullCode: "5032701"},
		Category:     CategoryServiceUnavailable,
		Description:  "Service unavailable",
	},
	"5032702": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 503, ServiceCode: 27, CaseCode: 2, FullCode: "5032702"},
		Category:     CategoryServiceUnavailable,
		Description:  "Rate limit exceeded",
	},
	"5032703": {
		ResponseCode: &BRIResponseCode{HTTPStatus: 503, ServiceCode: 27, CaseCode: 3, FullCode: "5032703"},
		Category:     CategoryServiceUnavailable,
		Description:  "Circuit breaker open",
	},
}

// GetBRIVAResponseDefinition returns detailed information about a BRIVA response code
func GetBRIVAResponseDefinition(code string) *BRIVAResponseDefinition {
	definition := brivaResponseDefinitions[code]
	if definition != nil {
		return definition
	}

	// Handle unknown response codes with pending status
	return getPendingResponseDefinition(code)
}

// getPendingResponseDefinition creates a default definition for unknown response codes
func getPendingResponseDefinition(code string) *BRIVAResponseDefinition {
	// Try to parse the response code to determine HTTP status
	var httpStatus = 500 // default
	if len(code) == 7 && isDigitString(code) {
		if status, err := strconv.Atoi(code[0:3]); err == nil {
			httpStatus = status
		}
	}

	// Determine category based on HTTP status
	var category HttpCategory
	switch {
	case httpStatus >= 200 && httpStatus < 300:
		category = CategorySuccess
	case httpStatus >= 400 && httpStatus < 500:
		category = CategoryBadRequest
	case httpStatus >= 500 && httpStatus < 600:
		category = CategoryInternalServerError
	default:
		// For non-standard HTTP statuses, treat as pending
		category = CategoryPending
	}

	return &BRIVAResponseDefinition{
		ResponseCode: &BRIResponseCode{
			HTTPStatus:  httpStatus,
			ServiceCode: 0,
			CaseCode:    0,
			FullCode:    code,
		},
		Category:    category,
		Description: fmt.Sprintf("Unknown response code: %s - Status pending, requires manual verification", code),
	}
}

// StructuredBRIAPIResponse provides response information from the API
type StructuredBRIAPIResponse struct {
	ResponseCode    string    // The actual response code from API
	ResponseMessage string    // The actual response message from API
	HTTPStatusCode  int       // HTTP status code
	Timestamp       time.Time // When the error occurred
}

// Error implements the error interface
func (e *StructuredBRIAPIResponse) Error() string {
	msg := fmt.Sprintf("BRI API Error [%s]: %s", e.ResponseCode, e.ResponseMessage)
	// Try to extract field name from response message for certain error types
	if field := e.extractFieldFromMessage(); field != "" {
		msg += fmt.Sprintf(" (field: %s)", field)
	}
	return msg
}

// extractFieldFromMessage attempts to extract field name from response message
func (e *StructuredBRIAPIResponse) extractFieldFromMessage() string {
	// Handle "Invalid Mandatory Field <fieldName>" pattern
	if strings.HasPrefix(e.ResponseMessage, "Invalid Mandatory Field ") {
		field := strings.TrimPrefix(e.ResponseMessage, "Invalid Mandatory Field ")
		return field
	}
	// Handle "Invalid Field Format <fieldName>" pattern
	if strings.HasPrefix(e.ResponseMessage, "Invalid Field Format ") {
		field := strings.TrimPrefix(e.ResponseMessage, "Invalid Field Format ")
		return field
	}
	// Handle "Invalid field format <fieldName>" pattern (lowercase)
	if strings.HasPrefix(e.ResponseMessage, "Invalid field format ") {
		field := strings.TrimPrefix(e.ResponseMessage, "Invalid field format ")
		return field
	}
	// Handle "Invalid field value <fieldName>" pattern
	if strings.HasPrefix(e.ResponseMessage, "Invalid field value ") {
		field := strings.TrimPrefix(e.ResponseMessage, "Invalid field value ")
		return field
	}
	return ""
}

// GetTimestamp returns when the error occurred
func (e *StructuredBRIAPIResponse) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetCategory returns the response category based on HTTP status code
func (e *StructuredBRIAPIResponse) GetCategory() HttpCategory {
	switch {
	case e.HTTPStatusCode >= 200 && e.HTTPStatusCode < 300:
		return CategorySuccess
	case e.HTTPStatusCode >= 400 && e.HTTPStatusCode < 500:
		return CategoryBadRequest
	case e.HTTPStatusCode >= 500 && e.HTTPStatusCode < 600:
		return CategoryInternalServerError
	default:
		return CategoryPending
	}
}

// IsSuccess checks if this is a success response
func (e *StructuredBRIAPIResponse) IsSuccess() bool {
	return e.HTTPStatusCode >= 200 && e.HTTPStatusCode < 300
}

// IsClientError checks if this is a client error (4xx)
func (e *StructuredBRIAPIResponse) IsClientError() bool {
	return e.HTTPStatusCode >= 400 && e.HTTPStatusCode < 500
}

// IsPending checks if this response has a pending status that requires manual verification
func (e *StructuredBRIAPIResponse) IsPending() bool {
	return e.GetCategory() == CategoryPending
}

// NewStructuredBRIAPIResponse creates a new structured BRI API response
func NewStructuredBRIAPIResponse(responseCode, responseMessage string) *StructuredBRIAPIResponse {
	// Extract HTTP status code from response code (first 3 digits)
	var httpStatusCode int
	if len(responseCode) >= 3 {
		if status, err := strconv.Atoi(responseCode[0:3]); err == nil {
			httpStatusCode = status
		} else {
			httpStatusCode = 500 // default fallback
		}
	} else {
		httpStatusCode = 500 // default fallback
	}

	return &StructuredBRIAPIResponse{
		ResponseCode:    responseCode,
		ResponseMessage: responseMessage,
		HTTPStatusCode:  httpStatusCode,
		Timestamp:       time.Now(),
	}
}

// ParseResponseCodeFromMessage extracts response code from API response message

// isDigitString checks if string contains only digits
func isDigitString(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
