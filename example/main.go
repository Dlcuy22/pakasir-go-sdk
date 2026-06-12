package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dlcuy22/pakasir-go-sdk"
)

func main() {
	slug := os.Getenv("PAKASIR_SLUG")
	apiKey := os.Getenv("PAKASIR_API")

	if slug == "" || apiKey == "" {
		log.Fatal("PAKASIR_SLUG and PAKASIR_API must be set")
	}

	client, err := pakasir.NewClient(pakasir.Config{
		Slug:   slug,
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	payment, err := client.CreatePayment(
		ctx,
		pakasir.PaymentMethodQRIS,
		"EXAMPLE-002",
		1000,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order ID:    %s\n", payment.OrderID)
	fmt.Printf("Amount:      Rp%d\n", payment.Amount)
	fmt.Printf("Fee:         Rp%d\n", payment.Fee)
	fmt.Printf("Total:       Rp%d\n", payment.TotalPayment)
	fmt.Printf("Status:      %s\n", payment.Status)
	fmt.Printf("Method:      %s\n", payment.PaymentMethod)
	if payment.PaymentNumber != nil {
		fmt.Printf("Payment No:  %s\n", *payment.PaymentNumber)
	}
	if payment.PaymentURL != nil {
		fmt.Printf("Payment URL: %s\n", *payment.PaymentURL)
	}
}
