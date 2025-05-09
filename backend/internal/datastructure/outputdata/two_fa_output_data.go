package outputdata

import "github.com/huydq/test/internal/domain/model/user"

// GenerateTwoFAOutputData represents the response for generating a 2FA token
type GenerateTwoFAOutputData struct {
	MFAType   int   `json:"mfa_type"`
	ExpiresIn int64 `json:"expires_in"`
}

// VerifyTwoFAOutputData represents the response for verifying a 2FA token
type VerifyTwoFAOutputData struct {
	Token string     `json:"token"`
	User  *user.User `json:"user"`
}

// CanResendCodeOutputData represents the response for checking if a code can be resent
type CanResendCodeOutputData struct {
	CanResend     bool `json:"can_resend"`
	RemainingTime int  `json:"remaining_time"` // in seconds
}
