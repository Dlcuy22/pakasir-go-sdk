package pakasir

import "time"

type Config struct {
	Slug   string
	APIKey string
}

type PaymentMethod string

const (
	PaymentMethodAll           PaymentMethod = "all"
	PaymentMethodQRIS          PaymentMethod = "qris"
	PaymentMethodPayPal        PaymentMethod = "paypal"
	PaymentMethodCIMBNiagaVA   PaymentMethod = "cimb_niaga_va"
	PaymentMethodBNIVA         PaymentMethod = "bni_va"
	PaymentMethodSampoernaVA   PaymentMethod = "sampoerna_va"
	PaymentMethodBNCVA         PaymentMethod = "bnc_va"
	PaymentMethodMaybankVA     PaymentMethod = "maybank_va"
	PaymentMethodPermataVA     PaymentMethod = "permata_va"
	PaymentMethodATMBersamaVA  PaymentMethod = "atm_bersama_va"
	PaymentMethodArthaGrahaVA  PaymentMethod = "artha_graha_va"
	PaymentMethodBRIVA         PaymentMethod = "bri_va"
)

type PaymentPayload struct {
	Project       string     `json:"project"`
	OrderID       string     `json:"order_id"`
	Amount        int        `json:"amount"`
	Fee           int        `json:"fee"`
	Status        string     `json:"status"`
	TotalPayment  int        `json:"total_payment"`
	PaymentMethod string     `json:"payment_method"`
	PaymentNumber *string    `json:"payment_number"`
	PaymentURL    *string    `json:"payment_url"`
	RedirectURL   *string    `json:"redirect_url"`
	ExpiredAt     *time.Time `json:"expired_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

type WatchOptions struct {
	Interval      time.Duration
	Timeout       time.Duration
	OnStatusChange func(payment PaymentPayload)
	OnError       func(err error)
}
