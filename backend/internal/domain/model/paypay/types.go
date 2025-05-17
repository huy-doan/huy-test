package model

import "github.com/huydq/test/internal/pkg/utils"

// PaymentDetail represents PayPay payment details stored as JSON in the database
// We're using the generic JSONField until we know the exact structure of the payment details
type PaymentDetail = utils.JSONField[any]

// Once we know the exact structure, we can replace it with a proper struct, for example:
/*
type PaymentDetail struct {
    PaymentMethod   string  `json:"payment_method"`
    CardType        string  `json:"card_type,omitempty"`
    LastFourDigits  string  `json:"last_four_digits,omitempty"`
    ExpirationDate  string  `json:"expiration_date,omitempty"`
    // Add other fields as needed when the structure becomes clear
}
*/
