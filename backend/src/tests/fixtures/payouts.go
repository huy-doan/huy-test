package fixtures

import (
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"gorm.io/gorm"
)

// Fixed test dates for payouts to ensure consistent test results
var (
	// BaseDate represents a fixed date for testing (April 1, 2025)
	BaseDate = time.Date(2025, 4, 1, 10, 0, 0, 0, time.UTC)

	// NextDay represents a day after the BaseDate (April 2, 2025)
	NextDay = time.Date(2025, 4, 2, 10, 0, 0, 0, time.UTC)

	// PastDay represents a day before the BaseDate (March 31, 2025)
	PastDay = time.Date(2025, 3, 31, 10, 0, 0, 0, time.UTC)

	// FutureDay represents a week after the BaseDate (April 8, 2025)
	FutureDay = time.Date(2025, 4, 8, 10, 0, 0, 0, time.UTC)

	// ZeroDate represents a fixed date to use for draft payouts' sent_date
	ZeroDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
)

// CreateTestPayout creates a test payout with the specified parameters
func CreateTestPayout(
	db *gorm.DB,
	user *models.User,
	total float64,
	totalCount int,
	sendingDate time.Time,
	sentDate time.Time,
	status int,
	transferApplyNo string,
) (*models.Payout, error) {
	payout := &models.Payout{
		UserID:                user.ID,
		Total:                 total,
		TotalCount:            totalCount,
		SendingDate:           sendingDate,
		SentDate:              sentDate,
		PayoutStatus:          status,
		AozoraTransferApplyNo: transferApplyNo,
	}

	err := db.Create(payout).Error
	if err != nil {
		return nil, err
	}

	return payout, nil
}

// CreateDraftPayout creates a test payout with draft status
func CreateDraftPayout(db *gorm.DB, user *models.User, sendingDate time.Time) (*models.Payout, error) {
	// Draft payouts use ZeroDate for sent_date since it can't be NULL
	return CreateTestPayout(
		db,
		user,
		10000.00,
		1,
		sendingDate,
		ZeroDate, // Use zero date for draft status
		models.PayoutStatusDraft,
		"TEST-DRAFT-"+time.Now().Format("20060102-150405"),
	)
}

// CreateProcessedPayout creates a test payout with processed status
func CreateProcessedPayout(db *gorm.DB, user *models.User, sendingDate time.Time, sentDate time.Time) (*models.Payout, error) {
	return CreateTestPayout(
		db,
		user,
		20000.00,
		2,
		sendingDate,
		sentDate, // Processed payouts have a real sent date
		models.PayoutStatusProcessed,
		"TEST-PROC-"+time.Now().Format("20060102-150405"),
	)
}
