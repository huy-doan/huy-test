package usecase

import (
	"context"
	"errors"
	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
	"time"
)

var (
	ErrAccountLocked     = errors.New("account.locked")
	ErrAccountTempLocked = errors.New("account.temp_locked")
	ErrAccountPermLocked = errors.New("account.perm_locked")
	ErrUnauthorized      = errors.New("account.unauthorized")
)

type LockedAccountUsecase struct {
	lockedAccountRepo repositories.LockedAccountRepository
	userRepo          repositories.UserRepository
}

func NewLockedAccountUsecase(
	lockedAccountRepo repositories.LockedAccountRepository,
	userRepo repositories.UserRepository,
) *LockedAccountUsecase {
	return &LockedAccountUsecase{
		lockedAccountRepo: lockedAccountRepo,
		userRepo:          userRepo,
	}
}

// getOrCreateLockedAccount gets an existing locked account or creates a new one
func (u *LockedAccountUsecase) GetOrCreateLockedAccount(ctx context.Context, email string) (*models.LockedAccount, error) {
	lockedAccount, err := u.lockedAccountRepo.GetByEmail(ctx, email)
	if err != nil {
		lockedAccount = &models.LockedAccount{
			Email:    email,
			IsLocked: false,
			Count:    0,
		}
	}

	user, _ := u.userRepo.FindByEmail(ctx, email)
	if user != nil {
		lockedAccount.UserID = &user.ID
	}

	return lockedAccount, nil
}

// handlePermanentLock handles permanent account locking
func (u *LockedAccountUsecase) handlePermanentLock(ctx context.Context, lockedAccount *models.LockedAccount) error {
	now := time.Now()
	lockedAccount.IsLocked = true
	lockedAccount.LockedAt = &now
	lockedAccount.ExpiredAt = nil

	err := u.lockedAccountRepo.Update(ctx, lockedAccount)
	if err != nil {
		return err
	}

	return ErrAccountPermLocked
}

// handleTemporaryLock handles temporary account locking
func (u *LockedAccountUsecase) handleTemporaryLock(ctx context.Context, lockedAccount *models.LockedAccount) error {
	if lockedAccount.ExpiredAt == nil {
		now := time.Now()
		expiredAt := now.Add(time.Duration(models.TempLockDuration) * time.Minute)
		lockedAccount.LockedAt = &now
		lockedAccount.ExpiredAt = &expiredAt
		lockedAccount.IsLocked = true
	}

	err := u.lockedAccountRepo.Update(ctx, lockedAccount)
	if err != nil {
		return err
	}

	return ErrAccountTempLocked
}

// updateLockedAccount updates or creates a locked account record
func (u *LockedAccountUsecase) updateLockedAccount(ctx context.Context, lockedAccount *models.LockedAccount) error {
	if lockedAccount.ID > 0 {
		return u.lockedAccountRepo.Update(ctx, lockedAccount)
	}

	return u.lockedAccountRepo.Create(ctx, lockedAccount)
}

// HandleFailedLogin processes a failed login attempt
func (u *LockedAccountUsecase) HandleFailedLogin(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil && user == nil {
		return u.handlePermanentLock(ctx, &models.LockedAccount{
			Email:    email,
			IsLocked: true,
		})
	}

	lockedAccount, err := u.GetOrCreateLockedAccount(ctx, email)
	if err != nil {
		return err
	}

	lockedAccount.Count++
	if lockedAccount.ShouldPermanentlyLock() {
		return u.handlePermanentLock(ctx, lockedAccount)
	} else if lockedAccount.ShouldTemporarilyLock() {
		return u.handleTemporaryLock(ctx, lockedAccount)
	}

	return u.updateLockedAccount(ctx, lockedAccount)
}

// CheckAccountStatus checks if an account is locked
func (u *LockedAccountUsecase) CheckAccountStatus(ctx context.Context, email string) error {
	lockedAccount, err := u.lockedAccountRepo.GetByEmail(ctx, email)
	if err == nil && lockedAccount.IsLocked {
		if lockedAccount.IsPermanentlyLocked() {
			return ErrAccountPermLocked
		}

		if lockedAccount.IsTemporarilyLocked() {
			return ErrAccountTempLocked
		}
	}

	return nil
}

// UnlockAccount unlocks an account (admin only)
func (u *LockedAccountUsecase) UnlockAccount(ctx context.Context, adminID int, targetUserID int) error {
	admin, err := u.userRepo.FindByID(ctx, adminID)
	if err != nil || admin.IsAdmin() {
		return ErrUnauthorized
	}

	lockedAccount, err := u.lockedAccountRepo.GetByUserID(ctx, targetUserID)
	if err != nil {
		return err
	}

	return u.ResetLockAccount(ctx, lockedAccount)
}

func (u *LockedAccountUsecase) UnlockAccountByEmail(ctx context.Context, email string) error {
	lockedAccount, err := u.lockedAccountRepo.GetByEmail(ctx, email)

	if err == nil {
		return u.ResetLockAccount(ctx, lockedAccount)
	}

	return nil
}

// ResetLockAccount resets all lock-related fields of a LockedAccount to their default values
func (u *LockedAccountUsecase) ResetLockAccount(ctx context.Context, lockedAccount *models.LockedAccount) error {
	lockedAccount.IsLocked = false
	lockedAccount.Count = 0
	lockedAccount.ExpiredAt = nil
	lockedAccount.LockedAt = nil

	return u.lockedAccountRepo.Update(ctx, lockedAccount)
}

// GetRemainingAttempts returns the remaining attempts before temporary and permanent locks
func (u *LockedAccountUsecase) GetRemainingAttempts(ctx context.Context, email string) (tempLockRemaining, permLockRemaining int, err error) {
	lockedAccount, err := u.GetOrCreateLockedAccount(ctx, email)
	if err != nil {
		return 0, 0, err
	}

	if lockedAccount == nil {
		return models.TempLockThreshold, models.PermLockThreshold, nil
	}

	tempLockRemaining = models.TempLockThreshold - lockedAccount.Count
	permLockRemaining = models.PermLockThreshold - lockedAccount.Count

	return tempLockRemaining, permLockRemaining, nil
}
