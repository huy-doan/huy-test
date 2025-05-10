package mapper

import (
	"github.com/labstack/echo/v4"
)

type MFARequiredResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *MFARequiredData `json:"data"`
}

type MFARequiredData struct {
	RequiresMFA bool             `json:"requires_mfa"`
	User        *MFARequiredUser `json:"user"`
	ExpiresIn   int64            `json:"expires_in"`
	MFATypeID   int              `json:"mfa_type_id"`
}

type MFARequiredUser struct {
	Email   string `json:"email"`
	MFAType string `json:"mfa_type"`
}

type MFARequiredMapper struct {
	ctx echo.Context
}

func NewMFARequiredMapper(ctx echo.Context) *MFARequiredMapper {
	return &MFARequiredMapper{
		ctx: ctx,
	}
}

func (m *MFARequiredMapper) ToMFARequiredResponse(
	email string,
	mfaType string,
	expiresIn int64,
	mfaTypeID int,
) *MFARequiredResponse {
	return &MFARequiredResponse{
		Success: true,
		Message: "2FA認証が必要です",
		Data: &MFARequiredData{
			RequiresMFA: true,
			User: &MFARequiredUser{
				Email:   email,
				MFAType: mfaType,
			},
			ExpiresIn: expiresIn,
			MFATypeID: mfaTypeID,
		},
	}
}
