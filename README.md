# Pakasir Payment Gateway Go SDK

Pakasir Go SDK is a lightweight Go client library for the [Pakasir](https://pakasir.com) payment gateway.
Supports QRIS and multi-bank Virtual Account.

[Installation](#installation) · [Quick Start](#quick-start) · [Configuration](#configuration) · [Payment Methods](#payment-methods) · [API Reference](#api-reference)

---

## Installation

```bash
go get github.com/dlcuy22/pakasir-go-sdk
```

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/dlcuy22/pakasir-go-sdk"
)

func main() {
    client, err := pakasir.NewClient(pakasir.Config{
        Slug:   "your-slug",
        APIKey: "your-api-key",
    })
    if err != nil {
        log.Fatal(err)
    }

    payment, err := client.CreatePayment(
        context.Background(),
        pakasir.PaymentMethodQRIS,
        "ORDER-12345",
        10000,
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", payment)
}
```

---

## Configuration

`NewClient` accepts a `Config` struct:

| Field    | Type     | Description                                                 |
| :------- | :------- | :---------------------------------------------------------- |
| `Slug`   | `string` | Required. [Project](https://app.pakasir.com/projects) slug. |
| `APIKey` | `string` | Required. [API key](https://app.pakasir.com/projects).      |

You can optionally set a custom HTTP client:

```go
client.SetHTTPClient(&http.Client{Timeout: 10 * time.Second})
```

---

## Payment Methods

For more information about payment methods, please visit [Pakasir Biaya](https://pakasir.com/p/pricing).

| Method           | Code                | Fee              | Go Constant                    |
| ---------------- | ------------------- | ---------------- | ------------------------------ |
| All Methods      | `all`               | Varies           | `PaymentMethodAll`             |
| QRIS             | `qris`              | 0.7% - 1%        | `PaymentMethodQRIS`            |
| BNI VA           | `bni_va`            | Rp3.500          | `PaymentMethodBNIVA`           |
| BRI VA           | `bri_va`            | Rp3.500          | `PaymentMethodBRIVA`           |
| CIMB Niaga VA    | `cimb_niaga_va`     | Rp3.500          | `PaymentMethodCIMBNiagaVA`     |
| Maybank VA       | `maybank_va`        | Rp3.500          | `PaymentMethodMaybankVA`       |
| Permata VA       | `permata_va`        | Rp3.500          | `PaymentMethodPermataVA`       |
| BNC VA           | `bnc_va`            | Rp3.500          | `PaymentMethodBNCVA`           |
| ATM Bersama VA   | `atm_bersama_va`    | Rp3.500          | `PaymentMethodATMBersamaVA`    |
| Sampoerna VA     | `sampoerna_va`      | Rp2.000          | `PaymentMethodSampoernaVA`     |
| Artha Graha VA   | `artha_graha_va`    | Rp2.000          | `PaymentMethodArthaGrahaVA`    |

---

## API Reference

> **Quick Jump:** [CreatePayment](#createpayment) · [GetPaymentURL](#getpaymenturl) · [DetailPayment](#detailpayment) · [CancelPayment](#cancelpayment) · [SimulationPayment](#simulationpayment) · [WatchPayment](#watchpayment)

### CreatePayment

Create a new payment transaction via API. All HTTP methods accept a `context.Context` as the first argument.

```go
payment, err := client.CreatePayment(ctx, pakasir.PaymentMethodQRIS, "ORDER-12345", 100000)
```

| Parameter | Type            | Description                                            |
| :-------- | :-------------- | :----------------------------------------------------- |
| `ctx`     | `context.Context` | Request context                                      |
| `method`  | `PaymentMethod` | Payment method constant                                |
| `orderID` | `string`        | Unique order ID (min 5 characters, URL-safe chars only) |
| `amount`  | `int`           | Amount in Rupiah (min Rp500)                           |

---

### GetPaymentURL

Generate payment URL synchronously without an API call. Useful for client-side redirects. Panics on invalid input.

```go
payment := client.GetPaymentURL(pakasir.PaymentMethodQRIS, "ORDER-12345", 100000)

fmt.Println(*payment.PaymentURL)
```

| Parameter   | Type            | Description                 |
| :---------- | :-------------- | :---------------------------- |
| `method`    | `PaymentMethod` | Payment method constant       |
| `orderID`   | `string`        | Unique order ID               |
| `amount`    | `int`           | Amount in Rupiah              |
| `redirectURL` | `...string`   | Optional redirect URL (variadic) |

---

### DetailPayment

Retrieve current status of a payment.

```go
detail, err := client.DetailPayment(ctx, "ORDER-12345", 100000)
```

| Parameter | Type              | Description      |
| :-------- | :---------------- | :--------------- |
| `ctx`     | `context.Context` | Request context  |
| `orderID` | `string`         | Order ID         |
| `amount`  | `int`            | Amount in Rupiah |

---

### CancelPayment

Cancel an existing pending payment.

```go
canceled, err := client.CancelPayment(ctx, "ORDER-12345", 100000)
```

---

### SimulationPayment

Simulate a successful payment for testing in Sandbox mode.

```go
simulated, err := client.SimulationPayment(ctx, "ORDER-12345", 100000)
```

---

### WatchPayment

Monitor payment status changes in real-time with polling. Returns a `*WatchHandle` whose `.Stop()` method stops the watcher and blocks until the goroutine exits.

```go
handle, err := client.WatchPayment(ctx, "ORDER-12345", 100000, pakasir.WatchOptions{
    Interval: 3 * time.Second,
    Timeout:  10 * time.Minute,
    OnStatusChange: func(p pakasir.PaymentPayload) {
        fmt.Println("Status:", p.Status)
    },
    OnError: func(err error) {
        log.Println("Watch error:", err)
    },
})
if err != nil {
    log.Fatal(err)
}
defer handle.Stop()
```

| Field            | Type                                    | Default      | Description                     |
| :--------------- | :-------------------------------------- | :----------- | :------------------------------ |
| `Interval`       | `time.Duration`                         | `3s`         | Polling interval                |
| `Timeout`        | `time.Duration`                         | `10m`        | Auto-stop timeout               |
| `OnStatusChange` | `func(PaymentPayload)`                  | -            | Callback on status change       |
| `OnError`        | `func(error)`                           | -            | Callback on error               |

---

### PaymentPayload

Struct returned by all payment methods:

```go
type PaymentPayload struct {
    Project       string
    OrderID       string
    Amount        int
    Fee           int
    Status        string       // "pending" | "canceled" | "completed"
    TotalPayment  int
    PaymentMethod string
    PaymentNumber *string      // nil when not available
    PaymentURL    *string      // nil when not available
    RedirectURL   *string      // nil when not available
    ExpiredAt     *time.Time   // nil when not available
    CompletedAt   *time.Time   // nil when not available
}
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature/my-feature`.
3. Commit your changes: `git commit -m 'Add some feature'`.
4. Push to the branch: `git push origin feature/my-feature`.
5. Open a Pull Request.

## Issues & Feedback

If you encounter any problems or have feature requests, please open an [issue](https://github.com/dlcuy22/pakasir-go-sdk/issues).

## License

Distributed under the **MIT License**. See `LICENSE` for details.
