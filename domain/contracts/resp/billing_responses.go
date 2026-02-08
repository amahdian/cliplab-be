package resp

type CheckoutResponse struct {
	URL           string `json:"url"`
	TransactionID string `json:"transactionId"`
}
