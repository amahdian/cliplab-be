package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/amahdian/cliplab-be/domain/contracts/req"
	"github.com/amahdian/cliplab-be/domain/model"
	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/storage"
	"github.com/google/uuid"
)

type BillingSvc interface {
	HandleWebhook(payload []byte, signature string) error
	CreateCheckout(userId uuid.UUID, priceId string) (string, string, error)
}

type billingSvc struct {
	ctx       context.Context
	stg       storage.PgStorage
	envs      *env.Envs
	creditSvc CreditSvc
}

func newBillingSvc(ctx context.Context, stg storage.PgStorage, envs *env.Envs, creditSvc CreditSvc) BillingSvc {
	return &billingSvc{
		ctx:       ctx,
		stg:       stg,
		envs:      envs,
		creditSvc: creditSvc,
	}
}

func (s *billingSvc) HandleWebhook(payload []byte, signature string) error {
	// TODO: Implementation of signature verification
	// For now, we trust the payload if the secret key is set but skip it for development if needed
	// In a real app, you MUST verify the signature from Paddle

	var event req.PaddleWebhook
	if err := json.Unmarshal(payload, &event); err != nil {
		return errs.Wrapf(err, "failed to unmarshal paddle webhook")
	}

	switch event.EventType {
	case "subscription.created", "subscription.updated", "subscription.activated":
		return s.handleSubscriptionEvent(event)
	}

	return nil
}

func (s *billingSvc) handleSubscriptionEvent(event req.PaddleWebhook) error {
	dataBytes, _ := json.Marshal(event.Data)
	var subData req.PaddleSubscriptionData
	if err := json.Unmarshal(dataBytes, &subData); err != nil {
		return errs.Wrapf(err, "failed to unmarshal subscription data")
	}

	userIdStr, ok := subData.CustomData["userId"].(string)
	if !ok {
		return fmt.Errorf("userId not found in custom data")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return errs.Wrapf(err, "invalid userId in custom data")
	}

	priceId := ""
	if len(subData.Items) > 0 {
		priceId = subData.Items[0].Price.ID
	}

	// Update User fields
	user, err := s.stg.User(s.ctx).FindById(userId)
	if err != nil {
		return errs.Wrapf(err, "failed to find user with ID %s", userId.String())
	}
	if user == nil {
		return fmt.Errorf("user with ID %s not found", userId.String())
	}

	oldPriceId := ""
	if user.PriceID != nil {
		oldPriceId = *user.PriceID
	}

	user.PaddleCustomerID = &subData.CustomerID
	user.SubscriptionID = &subData.ID
	user.PriceID = &priceId
	user.SubscriptionStatus = &subData.Status

	if err := s.stg.User(s.ctx).UpdateOne(user, true); err != nil {
		return errs.Wrapf(err, "failed to update user subscription info")
	}

	// Update or Create Subscription record
	sub, err := s.stg.Subscription(s.ctx).FindByPaddleID(subData.ID)
	if err != nil {
		return errs.Wrapf(err, "failed to check existing subscription")
	}

	if sub == nil {
		sub = &model.Subscription{
			UserID:               userId,
			PaddleSubscriptionID: subData.ID,
		}
	}

	sub.PriceID = priceId
	sub.Status = subData.Status
	if subData.CurrentBillingPeriod != nil {
		sub.CurrentPeriodStart = subData.CurrentBillingPeriod.StartsAt
		sub.CurrentPeriodEnd = subData.CurrentBillingPeriod.EndsAt
	}

	if err := s.stg.Subscription(s.ctx).UpsertOne(sub, false); err != nil {
		return errs.Wrapf(err, "failed to save subscription record")
	}

	// Add credits if it's a new subscription or an upgrade
	if (event.EventType == "subscription.created" || event.EventType == "subscription.activated") && priceId != oldPriceId {
		if credits, ok := model.PlanCredits[priceId]; ok {
			if err := s.creditSvc.AddCredits(userId, credits); err != nil {
				return errs.Wrapf(err, "failed to add credits for subscription")
			}
		}
	}

	return nil
}

func (s *billingSvc) CreateCheckout(userId uuid.UUID, priceId string) (string, string, error) {
	// 0. Fetch user to get email
	user, err := s.stg.User(s.ctx).FindById(userId)
	if err != nil || user == nil {
		return "", "", errs.Newf(errs.NotFound, err, "user not found")
	}

	// 1. Prepare Paddle transaction request
	payload := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"price_id": priceId,
				"quantity": 1,
			},
		},
		"customer": map[string]interface{}{
			"email": user.Email,
		},
		"custom_data": map[string]interface{}{
			"userId": userId.String(),
		},
		"collection_mode": "automatic",
	}

	payloadBytes, _ := json.Marshal(payload)

	// 2. Call Paddle API to create transaction
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.envs.Paddle.APIBaseURL+"/transactions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", "", errs.Wrapf(err, "failed to create paddle request")
	}

	req.Header.Set("Authorization", "Bearer "+s.envs.Paddle.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", errs.Wrapf(err, "paddle api call failed")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("paddle api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			ID                string `json:"id"`
			HostedCheckoutURL string `json:"hosted_checkout_url"`
			Checkout          struct {
				URL string `json:"url"`
			} `json:"checkout"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", errs.Wrapf(err, "failed to decode paddle response")
	}

	checkoutURL := result.Data.HostedCheckoutURL
	if checkoutURL == "" {
		checkoutURL = result.Data.Checkout.URL
	}

	// If it's still empty but we have an ID, or if it points to localhost (common in sandbox)
	// we can fallback to the standard Paddle checkout URL format
	if result.Data.ID != "" && (checkoutURL == "" || strings.Contains(checkoutURL, "localhost")) {
		baseUrl := "https://checkout.paddle.com/checkout/buy"
		if strings.Contains(s.envs.Paddle.APIBaseURL, "sandbox") {
			baseUrl = "https://sandbox-checkout.paddle.com/checkout/buy"
		}
		checkoutURL = fmt.Sprintf("%s?_ptxn=%s", baseUrl, result.Data.ID)
	}

	if checkoutURL == "" {
		return "", "", fmt.Errorf("paddle returned empty checkout url. full response: %s", string(body))
	}

	return checkoutURL, result.Data.ID, nil
}
