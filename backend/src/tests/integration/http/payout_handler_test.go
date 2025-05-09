package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/huydq/test/src/api/http/handlers"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	repoImpl "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/huydq/test/src/tests/fixtures"
	"github.com/huydq/test/src/usecase"
	"github.com/stretchr/testify/suite"
)

// PayoutHandlerTestSuite contains the setup for testing PayoutHandler
type PayoutHandlerTestSuite struct {
	tests.TestSuite
	payoutRepo       repositories.PayoutRepository
	payoutRecordRepo repositories.PayoutRecordRepository
	payoutUsecase    *usecase.PayoutUsecase
	payoutHandler    *handlers.PayoutHandler
	adminUser        *models.User
	adminRole        *models.Role
	businessUser     *models.User
	businessRole     *models.Role
	testPayouts      []*models.Payout
}

func (s *PayoutHandlerTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()

	// Initialize repositories
	s.payoutRepo = repoImpl.NewPayoutRepository(s.DB)
	s.payoutRecordRepo = repoImpl.NewPayoutRecordRepository(s.DB)

	// Create usecase with real repos
	s.payoutUsecase = usecase.NewPayoutUsecase(
		s.payoutRepo,
		s.payoutRecordRepo,
	)

	// Create handler with usecase
	s.payoutHandler = handlers.NewPayoutHandler(s.payoutUsecase)
}

func (s *PayoutHandlerTestSuite) SetupTest() {
	// Setup test admin user with role and permissions
	adminPermission := fixtures.GetMockPermission(1, "ユーザー管理", "USER_MANAGE")
	s.adminRole = fixtures.GetMockRoleWithPermissions(1, string(models.RoleCodeAdmin), []*models.Permission{adminPermission})
	s.adminUser = fixtures.GetMockUser("admin@example.com", "password123", s.adminRole)
	s.DB.Create(s.adminUser)

	// Setup business user for testing
	businessPermission := fixtures.GetMockPermission(8, "振込み承認（事業）", "TRANSFER_APPROVE_BUSINESS")
	s.businessRole = fixtures.GetMockRoleWithPermissions(3, string(models.RoleCodeBusinessUser), []*models.Permission{businessPermission})
	s.businessUser = fixtures.GetMockUser("business@example.com", "password123", s.businessRole)
	s.DB.Create(s.businessUser)

	// Create test payouts
	s.createTestPayouts()
}

func (s *PayoutHandlerTestSuite) TearDownTest() {
	// Clean up test data
	if len(s.testPayouts) > 0 {
		for _, payout := range s.testPayouts {
			s.DB.Exec("DELETE FROM payout_record WHERE payout_id = ?", payout.ID)
			s.DB.Exec("DELETE FROM payout WHERE id = ?", payout.ID)
		}
		s.testPayouts = nil
	}
	s.DB.Unscoped().Delete(&models.User{}, "email = ?", s.adminUser.Email)
	s.adminUser = nil
	s.DB.Unscoped().Delete(&models.User{}, "email = ?", s.businessUser.Email)
	s.businessUser = nil
}

func TestPayoutHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutHandlerTestSuite))
}

// Helper method to create test payouts
func (s *PayoutHandlerTestSuite) createTestPayouts() {
	// Check for and create test merchant records for foreign key references
	var merchant1Count int64
	s.DB.Model(&models.Merchant{}).Where("id = ?", 1).Count(&merchant1Count)

	if merchant1Count == 0 {
		merchant1 := &models.Merchant{
			ID:                1,
			PaymentProviderID: 1,
			PaymentMerchantID: "MERCH-001",
			MerchantName:      "Test Merchant 1",
			ShopID:            1001,
		}
		err := s.DB.Create(merchant1).Error
		s.NoError(err, "Should create test merchant 1")
	}

	var merchant2Count int64
	s.DB.Model(&models.Merchant{}).Where("id = ?", 2).Count(&merchant2Count)

	if merchant2Count == 0 {
		merchant2 := &models.Merchant{
			ID:                2,
			PaymentProviderID: 1,
			PaymentMerchantID: "MERCH-002",
			MerchantName:      "Test Merchant 2",
			ShopID:            1002,
		}
		err := s.DB.Create(merchant2).Error
		s.NoError(err, "Should create test merchant 2")
	}

	// Create a draft payout
	draftPayout, err := fixtures.CreateDraftPayout(s.DB, s.businessUser, fixtures.BaseDate)
	s.NoError(err)
	s.testPayouts = append(s.testPayouts, draftPayout)

	// Add some payout records to the draft payout
	for i := range 5 {
		record := &models.PayoutRecord{
			PayoutID:              draftPayout.ID,
			TransactionID:         1000 + i,
			ShopID:                1,
			Amount:                2000.00,
			BankName:              "Test Bank",
			BankCode:              "0001",
			BranchName:            "Test Branch",
			BranchCode:            "001",
			BankAccountType:       models.BankAccountTypeOrdinary,
			AccountName:           "Test Account",
			AccountNo:             "1234567890",
			TransferStatus:        models.TransferStatusInProgress,
			AozoraTransferApplyNo: "TEST-DRAFT-" + time.Now().Format("20060102-150405") + "-" + strconv.Itoa(i),
		}
		err := s.DB.Create(record).Error
		s.NoError(err) // Check for insertion error
	}

	// Create a processed payout
	processedPayout, err := fixtures.CreateProcessedPayout(s.DB, s.businessUser, fixtures.BaseDate, fixtures.NextDay)
	s.NoError(err)
	s.testPayouts = append(s.testPayouts, processedPayout)

	// Add some payout records to the processed payout
	for i := range 10 {
		record := &models.PayoutRecord{
			PayoutID:              processedPayout.ID,
			TransactionID:         2000 + i,
			ShopID:                2,
			Amount:                2000.00,
			BankName:              "Test Bank",
			BankCode:              "0002",
			BranchName:            "Test Branch",
			BranchCode:            "002",
			BankAccountType:       models.BankAccountTypeOrdinary,
			AccountName:           "Test Account",
			AccountNo:             "1234567890",
			TransferStatus:        models.TransferStatusProcessed,
			AozoraTransferApplyNo: "TEST-PROC-" + time.Now().Format("20060102-150405") + "-" + strconv.Itoa(i),
			SendingDate:           &fixtures.BaseDate,
			TransferRequestedAt:   &fixtures.BaseDate,
			TransferExecutedAt:    &fixtures.NextDay,
		}
		err := s.DB.Create(record).Error
		s.NoError(err) // Check for insertion error
	}
}

// createAdminAuthContext creates a context that mimics an authenticated admin user
func (s *PayoutHandlerTestSuite) createAdminAuthContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UserIDKey, s.adminUser.ID)
	ctx = context.WithValue(ctx, middleware.RoleCodeKey, string(models.RoleCodeAdmin))
	return ctx
}

func (s *PayoutHandlerTestSuite) TestListPayouts() {
	// Test with admin user authentication
	s.Run("AdminUser_Success", func() {
		// Create test request
		req := httptest.NewRequest("GET", "/admin/payouts?page=1&page_size=10", nil)
		req = req.WithContext(s.createAdminAuthContext())
		w := httptest.NewRecorder()

		// Execute handler directly
		s.payoutHandler.ListPayouts(w, req)

		// Verify response
		s.Equal(http.StatusOK, w.Code)

		// Parse response
		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)

		// Check response data
		s.Equal(resp.Success, true)

		// Extract data from response
		data, ok := resp.Data.(map[string]any)
		s.True(ok)

		// Assert payouts data exists
		payoutsData, ok := data["payouts"]
		s.True(ok)
		s.NotNil(payoutsData)

		// Verify our test payouts are included in the response
		payouts, ok := payoutsData.([]any)
		s.True(ok)
		s.GreaterOrEqual(len(payouts), 2)

		// Verify pagination data
		s.Equal(float64(1), data["page"])
		s.Equal(float64(10), data["page_size"])
		s.NotNil(data["total"])
	})
}

func (s *PayoutHandlerTestSuite) TestListPayouts_WithFilters() {
	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
		filterCheck    func(any) bool
	}{
		{
			name:           "Filter by Draft Status",
			queryParams:    "page=1&page_size=10&payout_status=1",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			filterCheck: func(payout any) bool {
				p, ok := payout.(map[string]any)
				return ok && p["payout_status"] == float64(models.PayoutStatusDraft)
			},
		},
		{
			name:           "Filter by Processed Status",
			queryParams:    "page=1&page_size=10&payout_status=3",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			filterCheck: func(payout any) bool {
				p, ok := payout.(map[string]any)
				return ok && p["payout_status"] == float64(models.PayoutStatusProcessed)
			},
		},
		{
			name:           "Filter by Invalid Date Format",
			queryParams:    "created_at=invalid-date",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			filterCheck:    nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Create test request with filter
			url := fmt.Sprintf("/admin/payouts?%s", tc.queryParams)
			req := httptest.NewRequest("GET", url, nil)
			req = req.WithContext(s.createAdminAuthContext())
			w := httptest.NewRecorder()

			// Execute handler directly
			s.payoutHandler.ListPayouts(w, req)

			// Verify response status
			s.Equal(tc.expectedStatus, w.Code)

			// For successful responses, check the content
			if tc.expectedStatus == http.StatusOK {
				var resp response.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				s.NoError(err)
				s.Equal(resp.Success, true)

				// Check response data
				data, ok := resp.Data.(map[string]any)
				s.True(ok)

				// Assert payouts data exists
				payoutsData, ok := data["payouts"]
				s.True(ok)

				payouts, ok := payoutsData.([]any)
				s.True(ok)

				// Verify count of payouts matches expected
				if tc.expectedCount > 0 {
					s.GreaterOrEqual(len(payouts), tc.expectedCount)
				}

				// Apply the filter check if provided
				if tc.filterCheck != nil {
					for _, payout := range payouts {
						s.True(tc.filterCheck(payout))
					}
				}
			}
		})
	}
}

func (s *PayoutHandlerTestSuite) TestPaginationAndSorting() {
	// Test with different page sizes
	s.Run("Pagination", func() {
		// Test page size 1
		req := httptest.NewRequest("GET", "/admin/payouts?page=1&page_size=1", nil)
		req = req.WithContext(s.createAdminAuthContext())
		w := httptest.NewRecorder()
		s.payoutHandler.ListPayouts(w, req)

		// Verify response
		s.Equal(http.StatusOK, w.Code)
		var resp1 response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp1)
		s.NoError(err)

		data1, ok := resp1.Data.(map[string]any)
		s.True(ok)
		payouts1, ok := data1["payouts"].([]any)
		s.True(ok)
		s.Equal(1, len(payouts1)) // Should return only 1 payout (page size 1)

		// Test page 2, size 1
		req = httptest.NewRequest("GET", "/admin/payouts?page=2&page_size=1", nil)
		req = req.WithContext(s.createAdminAuthContext())
		w = httptest.NewRecorder()
		s.payoutHandler.ListPayouts(w, req)

		// Verify response
		s.Equal(http.StatusOK, w.Code)
		var resp2 response.Response
		err = json.Unmarshal(w.Body.Bytes(), &resp2)
		s.NoError(err)

		data2, ok := resp2.Data.(map[string]any)
		s.True(ok)
		payouts2, ok := data2["payouts"].([]any)
		s.True(ok)
		s.Equal(1, len(payouts2)) // Should return 1 payout on second page

		// Verify pages contain different payouts
		payout1ID := payouts1[0].(map[string]any)["id"]
		payout2ID := payouts2[0].(map[string]any)["id"]
		s.NotEqual(payout1ID, payout2ID)
	})

	// Test sorting
	s.Run("Sorting", func() {
		// Test ascending order by ID
		req := httptest.NewRequest("GET", "/admin/payouts?sort_field=id&sort_order=ascend", nil)
		req = req.WithContext(s.createAdminAuthContext())
		w := httptest.NewRecorder()
		s.payoutHandler.ListPayouts(w, req)

		// Verify response
		s.Equal(http.StatusOK, w.Code)
		var resp response.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)

		data, ok := resp.Data.(map[string]any)
		s.True(ok)
		payouts, ok := data["payouts"].([]any)
		s.True(ok)
		s.GreaterOrEqual(len(payouts), 2)

		// Verify ascending sort
		firstID := payouts[0].(map[string]any)["id"].(float64)
		lastID := payouts[len(payouts)-1].(map[string]any)["id"].(float64)
		s.LessOrEqual(firstID, lastID)

		// Test descending order
		req = httptest.NewRequest("GET", "/admin/payouts?sort_field=id&sort_order=descend", nil)
		req = req.WithContext(s.createAdminAuthContext())
		w = httptest.NewRecorder()
		s.payoutHandler.ListPayouts(w, req)

		// Verify response
		s.Equal(http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		s.NoError(err)

		data, ok = resp.Data.(map[string]any)
		s.True(ok)
		payouts, ok = data["payouts"].([]any)
		s.True(ok)

		// Verify descending sort
		firstID = payouts[0].(map[string]any)["id"].(float64)
		lastID = payouts[len(payouts)-1].(map[string]any)["id"].(float64)
		s.GreaterOrEqual(firstID, lastID)
	})
}
