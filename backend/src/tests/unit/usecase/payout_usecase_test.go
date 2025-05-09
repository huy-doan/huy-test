package usecase_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/domain/repositories/filter"
	repoImpl "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/huydq/test/src/tests/fixtures"
	"github.com/huydq/test/src/usecase"
	"github.com/stretchr/testify/suite"
)

type PayoutUsecaseTestSuite struct {
	tests.TestSuite
	payoutRepo       repositories.PayoutRepository
	payoutRecordRepo repositories.PayoutRecordRepository
	payoutUsecase    *usecase.PayoutUsecase
	testUser         *models.User
	testPayouts      []*models.Payout
}

func (s *PayoutUsecaseTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()

	// Initialize repositories with the test DB
	s.payoutRepo = repoImpl.NewPayoutRepository(s.DB)
	s.payoutRecordRepo = repoImpl.NewPayoutRecordRepository(s.DB)

	// Create the usecase with actual repositories
	s.payoutUsecase = usecase.NewPayoutUsecase(
		s.payoutRepo,
		s.payoutRecordRepo,
	)
}

func (s *PayoutUsecaseTestSuite) SetupTest() {
	// Create a test user for the tests
	businessRole := fixtures.GetMockRoleWithPermissions(3, string(models.RoleCodeBusinessUser), nil)
	s.testUser = fixtures.GetMockUser("payout_test@example.com", "password123", businessRole)
	s.DB.Create(s.testUser)

	// Create test payouts
	s.createTestPayouts()
}

func (s *PayoutUsecaseTestSuite) TearDownTest() {
	// Clean up test data
	if len(s.testPayouts) > 0 {
		for _, payout := range s.testPayouts {
			s.DB.Exec("DELETE FROM payout_record WHERE payout_id = ?", payout.ID)
			s.DB.Exec("DELETE FROM payout WHERE id = ?", payout.ID)
		}
		s.testPayouts = nil
	}
	s.DB.Unscoped().Delete(&models.User{}, "email = ?", s.testUser.Email)
	s.testUser = nil
}

func TestPayoutUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutUsecaseTestSuite))
}

// Helper method to create test payouts
func (s *PayoutUsecaseTestSuite) createTestPayouts() {
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
	draftPayout, err := fixtures.CreateDraftPayout(s.DB, s.testUser, fixtures.BaseDate)
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
			AccountNo:             "1234567890",
			AccountName:           "Test Account",
			TransferStatus:        models.TransferStatusInProgress,
			SendingDate:           &fixtures.BaseDate,
			AozoraTransferApplyNo: "TEST-DRAFT-" + time.Now().Format("20060102-150405") + "-" + strconv.Itoa(i),
		}
		err := s.DB.Create(record).Error
		s.NoError(err) // Check for insertion error
	}

	// Create a processed payout
	processedPayout, err := fixtures.CreateProcessedPayout(s.DB, s.testUser, fixtures.BaseDate, fixtures.NextDay)
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
			AccountNo:             "1234567890",
			AccountName:           "Test Account",
			TransferStatus:        models.TransferStatusProcessed,
			SendingDate:           &fixtures.BaseDate,
			TransferRequestedAt:   &fixtures.BaseDate,
			TransferExecutedAt:    &fixtures.NextDay,
			AozoraTransferApplyNo: "TEST-" + time.Now().Format("20060102-150405") + "-" + strconv.Itoa(i),
		}
		err := s.DB.Create(record).Error
		s.NoError(err) // Check for insertion error
	}
}

func (s *PayoutUsecaseTestSuite) TestGetAll_NoFilter() {
	ctx := context.Background()
	payoutFilter := filter.NewPayoutFilter()

	// Call the method
	payouts, totalPages, total, err := s.payoutUsecase.GetAll(ctx, payoutFilter)

	// Assert
	s.NoError(err)
	s.GreaterOrEqual(totalPages, 1)
	s.GreaterOrEqual(total, int64(2)) // We should have at least our 2 test payouts
	s.GreaterOrEqual(len(payouts), 2) // We should have at least our 2 test payouts

	// Check if record counts were properly set
	for _, payout := range payouts {
		if payout.ID == s.testPayouts[0].ID {
			s.Equal(5, payout.PayoutRecordCount)
			s.Equal(10000.0, payout.PayoutRecordSumAmount)
		} else if payout.ID == s.testPayouts[1].ID {
			s.Equal(10, payout.PayoutRecordCount)
			s.Equal(20000.0, payout.PayoutRecordSumAmount)
		}
	}
}

func (s *PayoutUsecaseTestSuite) TestGetAll_WithStatusFilter() {
	ctx := context.Background()

	// Test cases for different filter combinations
	testCases := []struct {
		name              string
		filter            *filter.PayoutFilter
		expectedStatus    int
		minimumCount      int64
		expectFoundPayout bool
		expectedPayoutID  int
	}{
		{
			name: "Filter by Draft Status",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				status := models.PayoutStatusDraft
				f.PayoutStatus = &status
				return f
			}(),
			expectedStatus:    models.PayoutStatusDraft,
			minimumCount:      1,
			expectFoundPayout: true,
			expectedPayoutID:  0, // Will be set to actual ID
		},
		{
			name: "Filter by Processed Status",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				status := models.PayoutStatusProcessed
				f.PayoutStatus = &status
				return f
			}(),
			expectedStatus:    models.PayoutStatusProcessed,
			minimumCount:      1,
			expectFoundPayout: true,
			expectedPayoutID:  0, // Will be set to actual ID
		},
	}

	// Set expected payout IDs
	for i, tc := range testCases {
		if tc.expectedStatus == models.PayoutStatusDraft {
			testCases[i].expectedPayoutID = s.testPayouts[0].ID
		} else if tc.expectedStatus == models.PayoutStatusProcessed {
			testCases[i].expectedPayoutID = s.testPayouts[1].ID
		}
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Call the method
			payouts, _, total, err := s.payoutUsecase.GetAll(ctx, tc.filter)

			// Assert
			s.NoError(err)
			s.GreaterOrEqual(total, tc.minimumCount)

			// All returned payouts should have the expected status
			for _, payout := range payouts {
				s.Equal(tc.expectedStatus, payout.PayoutStatus)

				// If we expect a specific payout, verify it's in the results
				if tc.expectFoundPayout && payout.ID == tc.expectedPayoutID {
					// Additional assertions for the specific payout
					s.Equal(s.testUser.ID, payout.UserID)
					// Check the record count and sum are set
					if tc.expectedStatus == models.PayoutStatusDraft {
						s.Equal(5, payout.PayoutRecordCount)
						s.Equal(10000.0, payout.PayoutRecordSumAmount)
					} else if tc.expectedStatus == models.PayoutStatusProcessed {
						s.Equal(10, payout.PayoutRecordCount)
						s.Equal(20000.0, payout.PayoutRecordSumAmount)
					}
				}
			}
		})
	}
}

func (s *PayoutUsecaseTestSuite) TestGetAll_WithDateFilter() {
	ctx := context.Background()

	// Test cases for date filters
	testCases := []struct {
		name         string
		filter       *filter.PayoutFilter
		minimumCount int64
	}{
		{
			name: "Filter by Sending Date - Base Date",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SendingDate = &fixtures.BaseDate
				return f
			}(),
			minimumCount: 2, // Both payouts have BaseDate as sending date
		},
		{
			name: "Filter by Sent Date - Next Day",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SentDate = &fixtures.NextDay
				return f
			}(),
			minimumCount: 1, // One payout has NextDay as sent date
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Call the method
			payouts, _, total, err := s.payoutUsecase.GetAll(ctx, tc.filter)

			// Assert
			s.NoError(err)
			s.GreaterOrEqual(total, tc.minimumCount)
			s.GreaterOrEqual(len(payouts), int(tc.minimumCount))

			// Verify each payout has the expected user relationship loaded
			for _, payout := range payouts {
				s.Equal(s.testUser.ID, payout.UserID)
				// Verify the record count and sum are set
				s.GreaterOrEqual(payout.PayoutRecordCount, 1)
				s.GreaterOrEqual(payout.PayoutRecordSumAmount, 2000.0)
			}
		})
	}
}

func (s *PayoutUsecaseTestSuite) TestGetAll_WithPagination() {
	ctx := context.Background()

	// Create a filter with pagination settings
	payoutFilter := filter.NewPayoutFilter()
	payoutFilter.SetPagination(1, 1) // First page, one item per page

	// Call the method for first page
	payouts1, totalPages, total, err := s.payoutUsecase.GetAll(ctx, payoutFilter)

	// Assert first page
	s.NoError(err)
	s.GreaterOrEqual(total, int64(2)) // We should have at least our 2 test payouts
	s.GreaterOrEqual(totalPages, 2)   // With page size 1, we should have 2 pages
	s.Len(payouts1, 1)                // Should return only 1 payout (page size)

	// Get second page
	payoutFilter.SetPagination(2, 1)
	payouts2, _, _, err := s.payoutUsecase.GetAll(ctx, payoutFilter)

	// Assert second page
	s.NoError(err)
	s.Len(payouts2, 1) // Should return 1 payout on second page

	// Verify pages contain different payouts
	s.NotEqual(payouts1[0].ID, payouts2[0].ID)
}
