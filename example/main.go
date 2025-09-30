package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/nofendian17/gobriva"
)

func main() {
	// Configure slog to show info messages (debug logging uses Info level)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAxaAIcyRgIuB/8Ig8AmOKzM8/AoF2NaCDu9ddgESLDEUr28TD
n2zgY6dgzU4zsS3rCvW/W6kM0JgY9+DV7yZ9OmieN0FVf5tSi7G46dA6sgBLBrQ/
DygRlkx4V+yXrwj1mM/gCcvGrdnmHJL6iNayekzsGmAn/JzMzU9MTpnyng2htPdc
Sy4JUM3q0EpgPHrreMnsvwWQNLAmoH+m+lBsIj+f1stV5GHTghvGsi1K7K2TLgZx
SSISO50Syyq37b4hcCARAVSsMu/CCwSc5kEUwowVOZqwFI//vQ06AGkjPgebev3v
sLjn5fJQAXlrjF+l+Fcu3lf4gz2iu7D/dxiXXQIDAQABAoIBAAPt8WBEuWQ7Vx/m
8MNSMl86W02tjRzc3H5+j9168p7WvNwDHUNzjHSlNKYKASCFQendkz8uoDjnkhuG
DQrslrHaHnpvyKDsa2d1En2uEvSy6WSW1Rtaixr6pni54hXSlXuP3cVw5IbuN1Nn
7FmOZBpNb3MvNUXvFWo03COpKPgehhI2Zw8upzu3UxC5TtjHYfoFlyW80EP0pi9B
Vh/MolBXf4mL5DLErtovI+dibzn3OWpUp5bFJ6sIx+FqrvLFg9b9/uGZLvGd6h1H
y2G1KhZY+vJThO1bj2MXLcME9sZ3w998ddpi4/8lWRTQGpfKfhsq1+0LHwDP+90u
OoxC5ksCgYEA+tf/SjJa8hQ6f9YWBVUQPrKXn6O1IYd2SVDwuwSEiqixHxeeG0vu
lglPtkrutjzblppWTG6MN0pL2zhwNPo5o6So8r+esGh4AYt5BjQ/hVDME/STM5Af
23YsGr9LDkUFhxq9PcB1/pHmEIMxe0BNdLnU+VLx52WsajoNPKtOFdMCgYEAya/8
cAFYhQqHC32tW3TlwzBT5qZubiwD1OHLgyoKuVgKBJ40uYNaxwksTz6IEoQUWktr
FhhhoU4S0feiffNUCx9xU/fy+XP6HPNkj4EynbRfyFr/5ItqKSOVTMr7+nv/AcL7
97C9MCmk6FWFA+eCQCm2yB8zr84QjC70EOODcA8CgYATVyYS1XEXqyGbi6kk/hsD
ioeQQnILxMMFAh2dfcquWjVV3V9OYXtizBL+Tia7nFOd+AZhXECpXqwcmexk2Uoq
aN6x4L5egZ+HFvbc2JhxMfqaK0hSOHGMXT8nTMp/riiv8wrWQQmX+C3R5huhkiKm
tlFKa+/E1J0Hj7RHkjmyCwKBgDyPsq5zSQBSA/EIYOjIdkGhHmBw81Hzt4bR8klF
c4jqDcALPWvDLJv9fiehcDyXGoFuig5NbeuAxRf1Uv6c9UyNuXrsRjJvh9fvoe+R
bQB77BL+eD5JOqx1udwgS3+QgicmRIDAul5e8tys6U8d0jewDumSrPOKXd+qLbFw
j8QzAoGBAKyFfLxkBFkfRobIWqN4yO3228zJ066p0YnQq7fZ+S8SamLI75l6aMPG
ywoTfuSaliYUR3r17w8usFZ6+ssVw4qh08hfJ8Vu1QA9y1jGEzrgnUbpUWUpvxeO
SjcH/oB+CUW9hr0DRf84HA7YVu9VqPCKfPoDLWou0YcC9x3j6LcL
-----END RSA PRIVATE KEY-----
`

	brivaClient := gobriva.NewClient(gobriva.Config{
		PartnerID:     "ratnatravelbgrdev001",
		ClientID:      "nHF3UktU6DK8kjltoHHiAg8DsTPxNkiT",
		ClientSecret:  "gYe36iuWB8mpt5u8",
		PrivateKey:    privateKey,
		ChannelID:     "88888",
		IsSandbox:     true,
		Timeout:       10 * time.Second,
		Debug:         true,
		HTTPClient:    nil,
		Authenticator: nil,
	})

	// ============================================================================
	// COMPLETE BRI VIRTUAL ACCOUNT WORKFLOW EXAMPLE
	// ============================================================================

	fmt.Println("🚀 Starting Complete BRI Virtual Account Workflow Example")
	fmt.Println("==========================================================")

	// ============================================================================
	// 1. CREATE VIRTUAL ACCOUNT
	// ============================================================================
	fmt.Println("\n📝 Step 1: Creating Virtual Account")

	partnerServiceID := fmt.Sprintf("%-8s", "12345") // 8 character left-aligned: partnerServiceID + customerNo
	vaNumber := "1234567890123456"                   // 16 digit VA number for BRIVA
	createVaReq := gobriva.NewCreateVirtualAccountRequest(
		partnerServiceID,             // partnerServiceID (will be combined with customerNo)
		"CUST001",                    // customerNo
		vaNumber,                     // virtualAccountNo (16 digits for BRI VA)
		"John Doe - Invoice #INV001", // virtualAccountName
		"TRX001",                     // trxID
		100000.00,                    // amount
		"IDR",                        // currency
		"2025-12-31T23:59:59+07:00",  // expiredDate (ISO 8601 format)
	)
	// Note: partnerServiceID will be automatically formatted as "12345CUST" (8 chars, left-padded)

	createVaRes, err := brivaClient.CreateVirtualAccount(ctx, createVaReq)
	if err != nil {
		log.Fatalf("❌ Failed to create virtual account: %v", err)
	}
	fmt.Printf("✅ Virtual Account Created: %+v\n", createVaRes)

	// ============================================================================
	// 2. UPDATE VIRTUAL ACCOUNT
	// ============================================================================
	fmt.Println("\n📝 Step 2: Updating Virtual Account")

	updateVaReq := gobriva.NewUpdateVirtualAccountRequest(
		partnerServiceID,                     // partnerServiceID
		"CUST001",                            // customerNo
		"1234567890123456",                   // virtualAccountNo
		"John Doe Updated - Invoice #INV001", // virtualAccountName
		"TRX002",                             // trxID
		150000.00,                            // amount (increased)
		"IDR",                                // currency
		"2025-12-31T23:59:59+07:00",          // expiredDate
	)

	updateVaRes, err := brivaClient.UpdateVirtualAccount(ctx, updateVaReq)
	if err != nil {
		log.Fatalf("❌ Failed to update virtual account: %v", err)
	}
	fmt.Printf("✅ Virtual Account Updated: %+v\n", updateVaRes)

	// ============================================================================
	// 3. INQUIRY VIRTUAL ACCOUNT
	// ============================================================================
	fmt.Println("\n📝 Step 3: Inquiry Virtual Account")

	inquiryReq := gobriva.NewInquiryVirtualAccountRequest(
		partnerServiceID,   // partnerServiceID
		"CUST001",          // customerNo
		"1234567890123456", // virtualAccountNo
		"TRX003",           // trxID
	)

	inquiryRes, err := brivaClient.InquiryVirtualAccount(ctx, inquiryReq)
	if err != nil {
		log.Fatalf("❌ Failed to inquiry virtual account: %v", err)
	}
	fmt.Printf("✅ Virtual Account Inquiry: %+v\n", inquiryRes)

	// ============================================================================
	// 4. UPDATE VIRTUAL ACCOUNT STATUS TO PAID
	// ============================================================================
	fmt.Println("\n📝 Step 4: Updating Virtual Account Status to PAID")

	updateStatusReq := gobriva.NewUpdateVirtualAccountStatusRequest(
		partnerServiceID,   // partnerServiceID
		"CUST001",          // customerNo
		"1234567890123456", // virtualAccountNo
		"TRX004",           // trxID
		"Y",                // paidStatus (Y = paid, N = unpaid)
	)

	updateStatusRes, err := brivaClient.UpdateVirtualAccountStatus(ctx, updateStatusReq)
	if err != nil {
		log.Fatalf("❌ Failed to update virtual account status: %v", err)
	}
	fmt.Printf("✅ Virtual Account Status Updated to PAID: %+v\n", updateStatusRes)

	// ============================================================================
	// 5. INQUIRY VIRTUAL ACCOUNT STATUS
	// ============================================================================
	fmt.Println("\n📝 Step 5: Inquiry Virtual Account Status")

	inquiryStatusReq := &gobriva.InquiryVirtualAccountStatusRequest{
		PartnerServiceID: partnerServiceID, // 8 chars: partnerServiceID + customerNo
		CustomerNo:       "CUST001",
		VirtualAccountNo: "1234567890123456",
		InquiryRequestID: "INQ001",
	}

	inquiryStatusRes, err := brivaClient.InquiryVirtualAccountStatus(ctx, inquiryStatusReq)
	if err != nil {
		log.Fatalf("❌ Failed to inquiry virtual account status: %v", err)
	}
	fmt.Printf("✅ Virtual Account Status Inquiry: %+v\n", inquiryStatusRes)

	// ============================================================================
	// 6. GET VIRTUAL ACCOUNT REPORT
	// ============================================================================
	fmt.Println("\n📝 Step 6: Getting Virtual Account Report")

	reportReq := gobriva.NewVirtualAccountReportRequest(
		partnerServiceID, // partnerServiceID
		"2025-01-01",     // startDate
		"00:00:00",       // startTime
		"23:59:59",       // endTime
	)

	reportRes, err := brivaClient.GetVirtualAccountReport(ctx, reportReq)
	if err != nil {
		log.Fatalf("❌ Failed to get virtual account report: %v", err)
	}
	fmt.Printf("✅ Virtual Account Report: %+v\n", reportRes)

	// ============================================================================
	// 7. DELETE VIRTUAL ACCOUNT (Optional - cleanup)
	// ============================================================================
	fmt.Println("\n📝 Step 7: Deleting Virtual Account (Cleanup)")

	deleteReq := &gobriva.DeleteVirtualAccountRequest{
		PartnerServiceID: partnerServiceID, // 8 chars: partnerServiceID + customerNo
		CustomerNo:       "CUST001",
		VirtualAccountNo: "1234567890123456",
		TrxID:            "TRX005",
	}

	deleteRes, err := brivaClient.DeleteVirtualAccount(ctx, deleteReq)
	if err != nil {
		log.Printf("⚠️  Failed to delete virtual account (might already be deleted): %v", err)
	} else {
		fmt.Printf("✅ Virtual Account Deleted: %+v\n", deleteRes)
	}

	fmt.Println("\n🎉 Complete BRI Virtual Account Workflow Example Finished!")
	fmt.Println("==========================================================")
}
