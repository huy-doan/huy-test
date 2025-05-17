package mapper

import (
	"strconv"
	"time"

	"github.com/huydq/test/internal/controller/base"
	payoutModel "github.com/huydq/test/internal/domain/model/payout"
)

type PayoutResponse struct {
	ID                    string    `json:"id"`
	PayoutStatus          int       `json:"payout_status"`
	Total                 float64   `json:"total"`
	TotalCount            int       `json:"total_count"`
	SendingDate           time.Time `json:"sending_date"`
	SentDate              time.Time `json:"sent_date"`
	AozoraTransferApplyNo string    `json:"aozora_transfer_apply_no"`
	PayoutRecordCount     int       `json:"payout_record_count"`
	PayoutRecordSumAmount float64   `json:"payout_record_sum_amount"`
	CreatedAt             string    `json:"created_at"`
	UpdatedAt             string    `json:"updated_at"`
	PayoutIssuer          string    `json:"payout_issuer"`
}

type PayoutListSuccessResponse struct {
	Payouts []PayoutResponse `json:"payouts"`
	base.PaginationResponse
}

func toPayoutResponse(p *payoutModel.Payout) PayoutResponse {
	response := PayoutResponse{
		ID:                    strconv.Itoa(p.ID),
		PayoutStatus:          int(p.PayoutStatus),
		Total:                 p.Total,
		TotalCount:            p.TotalCount,
		SendingDate:           p.SendingDate,
		SentDate:              p.SentDate,
		AozoraTransferApplyNo: p.AozoraTransferApplyNo,
		PayoutRecordCount:     p.PayoutRecordCount,
		PayoutRecordSumAmount: p.PayoutRecordSumAmount,
		CreatedAt:             p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             p.UpdatedAt.Format(time.RFC3339),
	}

	if p.User != nil {
		response.PayoutIssuer = p.User.FullName
	}

	return response
}

func ToPayoutListSuccessResponse(payouts []*payoutModel.Payout, page, pageSize, totalPages int, total int64) PayoutListSuccessResponse {
	payoutResponses := make([]PayoutResponse, len(payouts))
	for i, p := range payouts {
		payoutResponses[i] = toPayoutResponse(p)
	}

	return PayoutListSuccessResponse{
		Payouts: payoutResponses,
		PaginationResponse: base.PaginationResponse{
			Page:       int(page),
			PageSize:   int(pageSize),
			TotalPages: int(totalPages),
			Total:      total,
		},
	}
}
