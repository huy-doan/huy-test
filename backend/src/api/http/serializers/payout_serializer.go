package serializers

import "github.com/vnlab/makeshop-payment/src/domain/models"

type PayoutSerializer struct {
	Payout *models.Payout
}

// NewPayoutSerializer creates a new PayoutSerializer
func NewPayoutSerializer(payout *models.Payout) *PayoutSerializer {
	return &PayoutSerializer{
		Payout: payout,
	}
}

func (s *PayoutSerializer) Serialize() any {
	if s.Payout == nil {
		return nil
	}

	result := map[string]any{
		"id":                       s.Payout.ID,
		"payout_status":            s.Payout.PayoutStatus,
		"total":                    s.Payout.Total,
		"total_count":              s.Payout.TotalCount,
		"sending_date":             s.Payout.SendingDate.Format("2006-01-02"),
		"sent_date":                s.Payout.SentDate.Format("2006-01-02"),
		"aozora_transfer_apply_no": s.Payout.AozoraTransferApplyNo,
		"payout_record_count":      s.Payout.PayoutRecordCount,
		"payout_record_sum_amount": s.Payout.PayoutRecordSumAmount,
		"created_at":               s.Payout.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":               s.Payout.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if s.Payout.User != nil {
		result["payout_issuer"] = s.Payout.User.GetFullName()
	}

	// if s.Payout.Approval != nil {
	// 	result["approval"] = NewApprovalSerializer(s.Payout.Approval).Serialize()
	// }

	return result
}

func SerializePayoutCollection(payouts []*models.Payout) []any {
	if payouts == nil {
		return nil
	}

	var serializedPayouts []any
	for _, payout := range payouts {
		serializedPayouts = append(serializedPayouts, NewPayoutSerializer(payout).Serialize())
	}

	return serializedPayouts
}
