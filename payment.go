package pakasir

import (
	"fmt"
	"math"
	"net/url"
	"time"
)

const baseAPIURL = "https://app.pakasir.com"

func (c *Client) GetPaymentURL(method PaymentMethod, orderID string, amount int, redirectURL ...string) PaymentPayload {
	orderID = sanitizeUrlSafe(orderID)

	if len(orderID) < 5 {
		panic("order ID must be at least 5 characters long")
	}
	if amount < 500 {
		panic("amount must be at least Rp500")
	}

	redirect := ""
	if len(redirectURL) > 0 {
		redirect = redirectURL[0]
	}

	expiredAt := time.Now().Add(24 * time.Hour)
	var fee int
	var paymentURL string

	switch method {
	case PaymentMethodAll:
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s", baseAPIURL, c.cfg.Slug, amount, orderID)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodQRIS:
		if amount > 105000 {
			fee = int(math.Round(0.01 * float64(amount)))
		} else {
			fee = int(math.Round(0.007*float64(amount) + 310))
		}
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s", baseAPIURL, c.cfg.Slug, amount, orderID)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
		paymentURL += "&qris_only=1"
	case PaymentMethodPayPal:
		if amount < 10000 {
			panic("amount must be at least Rp10,000")
		}
		fee = int(math.Max(math.Round(0.01*float64(amount)), 3000))
		paymentURL = fmt.Sprintf("%s/paypal/%s/%d?order_id=%s", baseAPIURL, c.cfg.Slug, amount, orderID)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodCIMBNiagaVA:
		fee = 3500
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodBNIVA:
		fee = 3500
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodSampoernaVA:
		fee = 2000
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodBNCVA:
		fee = 3500
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodMaybankVA:
		fee = 3500
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodPermataVA:
		fee = 3500
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodATMBersamaVA:
		fee = 3500
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodArthaGrahaVA:
		fee = 2000
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	case PaymentMethodBRIVA:
		fee = 3500
		paymentURL = fmt.Sprintf("%s/pay/%s/%d?order_id=%s&payment_method=%s", baseAPIURL, c.cfg.Slug, amount, orderID, method)
		if redirect != "" {
			paymentURL += "&redirect=" + url.QueryEscape(redirect)
		}
	default:
		panic("invalid payment method")
	}

	return PaymentPayload{
		Project:       c.cfg.Slug,
		OrderID:       orderID,
		Amount:        amount,
		Fee:           fee,
		Status:        "pending",
		TotalPayment:  amount + fee,
		PaymentMethod: string(method),
		PaymentURL:    &paymentURL,
		RedirectURL:   nil,
		ExpiredAt:     &expiredAt,
		CompletedAt:   nil,
	}
}
