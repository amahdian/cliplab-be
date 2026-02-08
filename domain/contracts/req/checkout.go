package req

type CreateCheckout struct {
	PriceID string `json:"priceId" binding:"required"`
}
