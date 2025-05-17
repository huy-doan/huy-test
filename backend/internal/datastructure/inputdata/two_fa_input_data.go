package inputdata

// VerifyTwoFAInputData represents input data for verifying a 2FA token
type VerifyTwoFAInputData struct {
	Email     string `json:"email" binding:"required,email"`
	Token     string `json:"token" binding:"required"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// GenerateTwoFAInputData represents input data for generating a 2FA token
type GenerateTwoFAInputData struct {
	UserID  int `json:"user_id" binding:"required"`
	MfaType int `json:"mfa_type" binding:"required"`
}

// CanResendCodeInputData represents input data for checking if a code can be resent
type CanResendCodeInputData struct {
	UserID  int `json:"user_id" binding:"required"`
	MfaType int `json:"mfa_type" binding:"required"`
}

// ResendCodeInputData represents input data for resending a 2FA code
type ResendCodeInputData struct {
	Email   string `json:"email" binding:"required,email"`
	MfaType int    `json:"mfa_type,omitempty"`
}
