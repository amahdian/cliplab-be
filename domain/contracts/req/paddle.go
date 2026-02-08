package req

import "time"

type PaddleWebhook struct {
	EventID    string      `json:"event_id"`
	EventType  string      `json:"event_type"`
	OccurredAt time.Time   `json:"occurred_at"`
	Data       interface{} `json:"data"`
}

type PaddleSubscriptionData struct {
	ID                   string                   `json:"id"`
	Status               string                   `json:"status"`
	CustomerID           string                   `json:"customer_id"`
	Items                []PaddleSubscriptionItem `json:"items"`
	CurrentBillingPeriod *PaddleBillingPeriod     `json:"current_billing_period"`
	CustomData           map[string]interface{}   `json:"custom_data"`
}

type PaddleSubscriptionItem struct {
	Price struct {
		ID string `json:"id"`
	} `json:"price"`
}

type PaddleBillingPeriod struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}
