package outputdata

import "github.com/huydq/test/internal/domain/model/user"

// LoginOutputData represents the result of a login operation
type LoginOutputData struct {
	Token       string     `json:"token,omitempty"`
	User        *user.User `json:"user"`
	RequiresMFA bool       `json:"requires_mfa"`
	MFAInfo     *MFAInfo   `json:"mfa_info,omitempty"`
}

// MFAInfo contains information about MFA requirements
type MFAInfo struct {
	Type      string `json:"type"`
	ExpiresIn int    `json:"expires_in"`
}
