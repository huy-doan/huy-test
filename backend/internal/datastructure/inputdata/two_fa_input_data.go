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
	MFAType int `json:"mfa_type" binding:"required"`
}

// CanResendCodeInputData represents input data for checking if a code can be resent
type CanResendCodeInputData struct {
	UserID  int `json:"user_id" binding:"required"`
	MFAType int `json:"mfa_type" binding:"required"`
}
