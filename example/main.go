package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"time"

	"github.com/nofendian17/gobriva"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	privateKey := `--PUT YOUR PRIVATE KEY HERE--`

	brivaClient := gobriva.NewClient(gobriva.Config{
		PartnerID:     "YOUR_PARTNER_ID",
		ClientID:      "YOUR_CLIENT_ID",
		ClientSecret:  "YOUR_SECRET_KEY",
		PrivateKey:    privateKey,
		ChannelID:     "YOUR_CHANNEL_ID",
		IsSandbox:     true,
		Timeout:       10 * time.Second,
		Debug:         true,
		HTTPClient:    nil,
		Authenticator: nil,
		Logger:        logger,
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

	// Configurable variables for all use cases
	var (
		clientServiceID = "22416"    // BRI assigned serviceID - change per client
		clientName      = "John Doe" // Customer name - change per customer
		invoiceAmount   = 10000.00   // Transaction amount in IDR - change as needed
		expiryDays      = 1          // Days until expiry - change as needed
	)

	// Standard variable assignments for all use cases
	serviceID := clientServiceID
	uniqueNum := uniqueNumber()
	partnerServiceID := padLeftSpace(serviceID, 8)
	customerNumber := padLeftSpace(uniqueNum, 20)
	vaNumber := padLeftSpace(serviceID+customerNumber, 28)
	vaName := fmt.Sprintf("%s - VA#%s", clientName, serviceID+customerNumber)
	trxID := uniqueNum
	amount := invoiceAmount
	currency := "IDR"
	expireDate := time.Now().AddDate(0, 0, expiryDays).Format("2006-01-02T15:04:05-07:00")

	createVaReq := gobriva.NewCreateVirtualAccountRequest(
		partnerServiceID,
		customerNumber,
		vaNumber,
		vaName,
		trxID,
		amount,
		currency,
		expireDate,
	)

	createVaRes, err := brivaClient.CreateVirtualAccount(ctx, createVaReq)
	if err != nil {
		log.Fatalf("❌ Failed to create virtual account: %v", err)
	}
	fmt.Printf("✅ Virtual Account Created: %+v\n", createVaRes)

	// ============================================================================
	// 2. UPDATE VIRTUAL ACCOUNT
	// ============================================================================
	fmt.Println("\n📝 Step 2: Updating Virtual Account")

	// Update specific variables
	var (
		updateCustomerNo       = "CUST001"                            // Customer number for update
		updateVirtualAccountNo = "1234567890123456"                   // Specific VA number to update
		updateCustomerName     = "John Doe Updated - Invoice #INV001" // Updated customer name
		updateTrxID            = "TRX002"                             // Update transaction ID
		updateAmount           = 150000.00                            // Updated amount
		updateExpiryDate       = "2025-12-31T23:59:59+07:00"          // Updated expiry date
	)

	// Use existing variables for update operation (already declared above)
	// serviceID, uniqueNum, partnerServiceID, customerNumber, vaNumber, vaName, trxID, amount, currency, expireDate are already set

	updateVaReq := gobriva.NewUpdateVirtualAccountRequest(
		partnerServiceID,       // partnerServiceID
		updateCustomerNo,       // customerNo
		updateVirtualAccountNo, // virtualAccountNo
		updateCustomerName,     // virtualAccountName
		updateTrxID,            // trxID
		updateAmount,           // amount (increased)
		currency,               // currency
		updateExpiryDate,       // expiredDate
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

	// Using existing variables for consistency across all use cases
	inquiryReq := gobriva.NewInquiryVirtualAccountRequest(
		partnerServiceID,       // partnerServiceID
		updateCustomerNo,       // customerNo (using same as update)
		updateVirtualAccountNo, // virtualAccountNo (using same as update)
		trxID,                  // trxID
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

	// Using existing variables for consistency across all use cases
	updateStatusReq := gobriva.NewUpdateVirtualAccountStatusRequest(
		partnerServiceID,       // partnerServiceID
		updateCustomerNo,       // customerNo
		updateVirtualAccountNo, // virtualAccountNo
		trxID,                  // trxID
		"Y",                    // paidStatus (Y = paid, N = unpaid)
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

	// Using existing variables for consistency across all use cases
	inquiryStatusReq := &gobriva.InquiryVirtualAccountStatusRequest{
		PartnerServiceID: partnerServiceID, // 8 chars: partnerServiceID + customerNo
		CustomerNo:       updateCustomerNo,
		VirtualAccountNo: updateVirtualAccountNo,
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

	// Using existing variables for consistency across all use cases
	deleteReq := &gobriva.DeleteVirtualAccountRequest{
		PartnerServiceID: partnerServiceID, // 8 chars: partnerServiceID + customerNo
		CustomerNo:       updateCustomerNo,
		VirtualAccountNo: updateVirtualAccountNo,
		TrxID:            trxID,
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

func uniqueNumber() string {
	// m = 10^20
	m := new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)
	n, _ := rand.Int(rand.Reader, m)
	return n.Text(10)
}

func padLeftSpace(s string, length int) string {
	return fmt.Sprintf("%*s", length, s)
}
