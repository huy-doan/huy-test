package repository

import (
	"context"
	"testing"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"github.com/huydq/test/src/domain/repositories/filter"
	dbRepositories "github.com/huydq/test/src/infrastructure/persistence/repositories"
	"github.com/huydq/test/src/tests"
	"github.com/huydq/test/src/tests/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PayoutRepositoryTestSuite struct {
	tests.TestSuite
	payoutRepo repositories.PayoutRepository
	roleRepo   repositories.RoleRepository
	userRepo   repositories.UserRepository
	testUser   *models.User
	adminRole  *models.Role
}

func (s *PayoutRepositoryTestSuite) SetupSuite() {
	// Call the parent SetupSuite method
	s.TestSuite.SetupSuite()

	// Initialize repositories
	s.payoutRepo = dbRepositories.NewPayoutRepository(s.DB)
	s.roleRepo = dbRepositories.NewRoleRepository(s.DB)
	s.userRepo = dbRepositories.NewUserRepository(s.DB)

	// Setup test data
	s.setupTestData()
}

func (s *PayoutRepositoryTestSuite) setupTestData() {
	ctx := context.Background()

	// Get admin role
	s.adminRole = fixtures.GetAdminRole(ctx, s.roleRepo)
	require.NotNil(s.T(), s.adminRole, "Admin role should be available")

	// Create a test user with a unique ID
	testUser, err := fixtures.CreateUniqueTestUser(ctx, s.userRepo, "payout_test", "password123", s.adminRole)
	require.NoError(s.T(), err, "Should create test user")
	s.testUser = testUser
}

func (s *PayoutRepositoryTestSuite) TearDownSuite() {
	// Clean up test data
	if s.testUser != nil {
		s.DB.Delete(&models.User{}, s.testUser.ID)
	}
}

// TestList updates the test for the List method to account for the new return values
func (s *PayoutRepositoryTestSuite) TestList() {
	ctx := context.Background()

	// Use the fixed test dates from the fixtures
	baseDate := fixtures.BaseDate
	nextDay := fixtures.NextDay
	zeroDate := fixtures.ZeroDate

	// Define test payouts with different statuses and dates
	testPayouts := []*models.Payout{
		{
			UserID:                s.testUser.ID,
			Total:                 10000,
			TotalCount:            1,
			SendingDate:           baseDate,
			SentDate:              zeroDate, // Zero time for draft status
			PayoutStatus:          models.PayoutStatusDraft,
			AozoraTransferApplyNo: "TEST-001",
		},
		{
			UserID:                s.testUser.ID,
			Total:                 20000,
			TotalCount:            2,
			SendingDate:           nextDay,
			SentDate:              zeroDate, // Zero time for draft status
			PayoutStatus:          models.PayoutStatusDraft,
			AozoraTransferApplyNo: "TEST-002",
		},
		{
			UserID:                s.testUser.ID,
			Total:                 30000,
			TotalCount:            3,
			SendingDate:           baseDate,
			SentDate:              baseDate, // Same as sending date for processed status
			PayoutStatus:          models.PayoutStatusProcessed,
			AozoraTransferApplyNo: "TEST-003",
		},
	}

	// Insert test payouts
	for _, payout := range testPayouts {
		err := s.DB.Create(&payout).Error
		require.NoError(s.T(), err, "Should create test payout")
	}

	// Clean up after test
	defer func() {
		for _, payout := range testPayouts {
			s.DB.Delete(&models.Payout{}, payout.ID)
		}
	}()

	testCases := []struct {
		name                      string
		filter                    *filter.PayoutFilter
		expectedMinimumCount      int
		expectError               bool
		expectedMinimumTotalPages int
		expectedMinimumTotalCount int64
	}{
		{
			name: "List all with no filter",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SetPagination(1, 10)
				return f
			}(),
			expectedMinimumCount:      len(testPayouts),
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: int64(len(testPayouts)),
		},
		{
			name: "Filter by status - draft",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				status := models.PayoutStatusDraft
				f.PayoutStatus = &status
				f.SetPagination(1, 10) // Explicitly set pagination
				return f
			}(),
			expectedMinimumCount:      2,
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: 2,
		},
		{
			name: "Filter by status - processed",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				status := models.PayoutStatusProcessed
				f.PayoutStatus = &status
				f.SetPagination(1, 10) // Explicitly set pagination
				return f
			}(),
			expectedMinimumCount:      1,
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: 1,
		},
		{
			name: "Filter by sending date - base date",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SendingDate = &baseDate
				f.SetPagination(1, 10) // Explicitly set pagination
				return f
			}(),
			expectedMinimumCount:      2, // Two payouts scheduled for the base date
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: 2,
		},
		{
			name: "Filter by sending date - next day",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SendingDate = &nextDay
				f.SetPagination(1, 10) // Explicitly set pagination
				return f
			}(),
			expectedMinimumCount:      1, // One payout scheduled for the next day
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: 1,
		},
		{
			name: "Filter by sent date - base date",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SentDate = &baseDate
				f.SetPagination(1, 10) // Explicitly set pagination
				return f
			}(),
			expectedMinimumCount:      1, // One payout already sent on base date
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: 1,
		},
		{
			name: "Pagination: page 1, size 2",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SetPagination(1, 2)
				return f
			}(),
			expectedMinimumCount:      2, // Should return 2 payouts due to page size
			expectError:               false,
			expectedMinimumTotalPages: 2, // Total of 3 payouts with page size 2 = 2 pages
			expectedMinimumTotalCount: 3,
		},
		{
			name: "Pagination: page 2, size 2",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SetPagination(2, 2)
				return f
			}(),
			expectedMinimumCount:      1, // Should return 1 payout on second page
			expectError:               false,
			expectedMinimumTotalPages: 2,
			expectedMinimumTotalCount: 3,
		},
		{
			name: "Sort by sending_date desc",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SetSort("sending_date", "desc")
				f.SetPagination(1, 10) // Explicitly set pagination
				return f
			}(),
			expectedMinimumCount:      len(testPayouts),
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: int64(len(testPayouts)),
		},
		{
			name: "Sort by invalid field",
			filter: func() *filter.PayoutFilter {
				f := filter.NewPayoutFilter()
				f.SetSort("invalid_field", "desc") // Invalid field
				f.SetPagination(1, 10)             // Explicitly set pagination
				return f
			}(),
			expectedMinimumCount:      len(testPayouts), // Should still return results with default sort
			expectError:               false,
			expectedMinimumTotalPages: 1,
			expectedMinimumTotalCount: int64(len(testPayouts)),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			payouts, totalPages, totalCount, err := s.payoutRepo.List(ctx, tc.filter)

			if tc.expectError {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
				// Use GreaterOrEqual for counts to handle possible extra records in database
				assert.GreaterOrEqual(s.T(), len(payouts), tc.expectedMinimumCount)
				assert.GreaterOrEqual(s.T(), totalPages, tc.expectedMinimumTotalPages)
				assert.GreaterOrEqual(s.T(), totalCount, tc.expectedMinimumTotalCount)

				// Verify that User relationship is properly preloaded
				if len(payouts) > 0 {
					assert.NotNil(s.T(), payouts[0].User)
				}

				// Verify sort order if specified
				if len(tc.filter.Sort) > 0 &&
					tc.filter.Sort[0].Field == "sending_date" &&
					tc.filter.Sort[0].Direction == filter.Descending &&
					len(payouts) >= 2 {
					// In descending order, first payout's sending date should be later than or equal to second's
					assert.True(s.T(), !payouts[0].SendingDate.Before(payouts[1].SendingDate))
				}
			}
		})
	}
}

func TestPayoutRepository(t *testing.T) {
	suite.Run(t, new(PayoutRepositoryTestSuite))
}
