package service

// and optionally a description or type for better documentation/formatting.
type CSVHeaderMappingEntry struct {
	JP   string // Japanese or original header
	EN   string // Normalized (json/English) header
	Desc string // Optional: description or note
}

// CSVHeaderMappings is a list of all header mappings
var CSVHeaderMappings = []CSVHeaderMappingEntry{
	{JP: "決済番号", EN: "payment_transaction_id", Desc: "Payment Transaction ID"},
	{JP: "加盟店ID", EN: "payment_merchant_id", Desc: "Merchant ID"},
	{JP: "屋号", EN: "merchant_business_name", Desc: "Merchant Business Name"},
	{JP: "店舗ID", EN: "shop_id", Desc: "Shop ID"},
	{JP: "店舗名", EN: "shop_name", Desc: "Shop Name"},
	{JP: "端末番号/PosID", EN: "terminal_code", Desc: "Terminal Code"},
	{JP: "取引ステータス", EN: "payment_transaction_status", Desc: "Transaction Status"},
	{JP: "取引日時", EN: "transaction_at", Desc: "Transaction At"},
	{JP: "取引金額", EN: "transaction_amount", Desc: "Transaction Amount"},
	{JP: "レシート番号", EN: "receipt_number", Desc: "Receipt Number"},
	{JP: "支払い方法", EN: "paypay_payment_method", Desc: "Payment Method"},
	{JP: "マーチャント決済ID", EN: "merchant_order_id", Desc: "Merchant Order ID"},
	{JP: "加盟店決済ID", EN: "merchant_order_id", Desc: "Merchant Order ID"},
	{JP: "支払い詳細", EN: "payment_detail", Desc: "Payment Detail"},
	{JP: "法人名", EN: "corporate_name", Desc: "Corporate Name"},
	{JP: "締め日", EN: "cutoff_date", Desc: "Cutoff Date"},
	{JP: "支払日", EN: "payment_date", Desc: "Payment Date"},
	{JP: "取引額", EN: "transaction_amount", Desc: "Transaction Amount (alt)"},
	{JP: "返金額", EN: "refund_amount", Desc: "Refund Amount"},
	{JP: "利用料", EN: "usage_fee", Desc: "Usage Fee"},
	{JP: "プラットフォーム使用料", EN: "platform_fee", Desc: "Platform Fee"},
	{JP: "初期費用", EN: "initial_fee", Desc: "Initial Fee"},
	{JP: "税", EN: "tax", Desc: "Tax"},
	{JP: "キャッシュバック", EN: "cashback", Desc: "Cashback"},
	{JP: "調整額", EN: "adjustment", Desc: "Adjustment"},
	{JP: "入金手数料", EN: "fee", Desc: "Fee"},
	{JP: "支払金額", EN: "amount", Desc: "Amount"},
}

// CSVHeaderMapping is a map for fast lookup (JP -> EN)
var CSVHeaderMapping = func() map[string]string {
	m := make(map[string]string, len(CSVHeaderMappings))
	for _, entry := range CSVHeaderMappings {
		m[entry.JP] = entry.EN
	}
	return m
}()

// CSVHeaderMappingService provides functionality for mapping CSV headers
type CSVHeaderMappingService struct{}

// NewCSVHeaderMappingService creates a new instance of CSVHeaderMappingService
func NewCSVHeaderMappingService() *CSVHeaderMappingService {
	return &CSVHeaderMappingService{}
}

// MapHeaders normalizes headers from Japanese to English and maps CSV records accordingly
func (s *CSVHeaderMappingService) MapHeaders(headers []string, records [][]string) ([]string, []map[string]string) {
	normHeaders := make([]string, len(headers))
	for i, h := range headers {
		normHeaders[i] = s.NormalizeHeader(h)
	}
	var mappedRecords []map[string]string
	for _, row := range records {
		record := make(map[string]string)
		for i, v := range row {
			if i < len(normHeaders) {
				record[normHeaders[i]] = v
			}
		}
		mappedRecords = append(mappedRecords, record)
	}
	return normHeaders, mappedRecords
}

// NormalizeHeader maps a header to its normalized form if available
func (s *CSVHeaderMappingService) NormalizeHeader(header string) string {
	if norm, ok := CSVHeaderMapping[header]; ok {
		return norm
	}
	return header // fallback: return as is if not found
}
