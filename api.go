package pakasir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type transactionRequest struct {
	Project string `json:"project"`
	APIKey  string `json:"api_key"`
	OrderID string `json:"order_id"`
	Amount  int    `json:"amount"`
}

type transactionResponse struct {
	Data        interface{} `json:"data"`
	Message     string      `json:"message"`
	Payment     *struct {
		PaymentNumber string `json:"payment_number"`
		ExpiredAt     string `json:"expired_at"`
	} `json:"payment"`
	Transaction *struct {
		PaymentMethod string `json:"payment_method"`
		Status        string `json:"status"`
		CompletedAt   string `json:"completed_at"`
	} `json:"transaction"`
	Success bool `json:"success"`
}

func (c *Client) doPost(ctx context.Context, url string, body interface{}) (*transactionResponse, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var res transactionResponse
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if res.Data == nil && res.Payment == nil && !res.Success && res.Transaction == nil {
		if res.Message != "" {
			return nil, fmt.Errorf("%s", res.Message)
		}
		return nil, fmt.Errorf("failed to process request")
	}

	return &res, nil
}

func (c *Client) doGet(ctx context.Context, url string) (*transactionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var res transactionResponse
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if res.Data == nil && res.Transaction == nil {
		if res.Message != "" {
			return nil, fmt.Errorf("%s", res.Message)
		}
		return nil, fmt.Errorf("failed to get payment detail")
	}

	return &res, nil
}

func (c *Client) CreatePayment(ctx context.Context, method PaymentMethod, orderID string, amount int) (PaymentPayload, error) {
	orderID = sanitizeUrlSafe(orderID)
	payload := c.GetPaymentURL(method, orderID, amount)

	body := transactionRequest{
		Project: payload.Project,
		APIKey:  c.cfg.APIKey,
		OrderID: payload.OrderID,
		Amount:  payload.Amount,
	}

	url := fmt.Sprintf("%s/api/transactioncreate/%s", baseAPIURL, method)
	res, err := c.doPost(ctx, url, body)
	if err != nil {
		return PaymentPayload{}, err
	}

	if res.Payment != nil {
		if res.Payment.PaymentNumber != "" {
			payload.PaymentNumber = &res.Payment.PaymentNumber
		}
		if res.Payment.ExpiredAt != "" {
			t, err := time.Parse(time.RFC3339Nano, res.Payment.ExpiredAt)
			if err == nil {
				payload.ExpiredAt = &t
			}
		}
	}

	return payload, nil
}

func (c *Client) DetailPayment(ctx context.Context, orderID string, amount int) (PaymentPayload, error) {
	orderID = sanitizeUrlSafe(orderID)

	u := fmt.Sprintf("%s/api/transactiondetail?project=%s&amount=%d&order_id=%s&api_key=%s",
		baseAPIURL, c.cfg.Slug, amount, orderID, c.cfg.APIKey)

	res, err := c.doGet(ctx, u)
	if err != nil {
		return PaymentPayload{}, err
	}

	payload := c.GetPaymentURL(PaymentMethod(res.Transaction.PaymentMethod), orderID, amount)

	payload.Status = res.Transaction.Status
	payload.PaymentNumber = nil

	if res.Transaction.CompletedAt != "" {
		t, err := time.Parse(time.RFC3339Nano, res.Transaction.CompletedAt)
		if err == nil {
			payload.CompletedAt = &t
		}
	}
	payload.ExpiredAt = nil

	return payload, nil
}

func (c *Client) CancelPayment(ctx context.Context, orderID string, amount int) (PaymentPayload, error) {
	orderID = sanitizeUrlSafe(orderID)

	body := transactionRequest{
		Project: c.cfg.Slug,
		APIKey:  c.cfg.APIKey,
		OrderID: orderID,
		Amount:  amount,
	}

	_, err := c.doPost(ctx, baseAPIURL+"/api/transactioncancel", body)
	if err != nil {
		return PaymentPayload{}, err
	}

	payload, err := c.DetailPayment(ctx, orderID, amount)
	if err != nil {
		return PaymentPayload{}, err
	}

	payload.Status = "canceled"
	return payload, nil
}

func (c *Client) SimulationPayment(ctx context.Context, orderID string, amount int) (PaymentPayload, error) {
	orderID = sanitizeUrlSafe(orderID)

	body := transactionRequest{
		Project: c.cfg.Slug,
		APIKey:  c.cfg.APIKey,
		OrderID: orderID,
		Amount:  amount,
	}

	_, err := c.doPost(ctx, baseAPIURL+"/api/paymentsimulation", body)
	if err != nil {
		return PaymentPayload{}, err
	}

	payload, err := c.DetailPayment(ctx, orderID, amount)
	if err != nil {
		return PaymentPayload{}, err
	}

	payload.Status = "completed"
	return payload, nil
}


