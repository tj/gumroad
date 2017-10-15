// Package gumroad provides a Gumroad API client.
package gumroad

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// form mime type.
const formMime = "application/x-www-form-urlencoded"

// New client.
func New() *Client {
	return &Client{}
}

// Client is the api client.
type Client struct {
	Licenses
}

// Error is an api error.
type Error struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Error implementation.
func (e Error) Error() string {
	return e.Message
}

// Licenses is the license api client.
type Licenses struct{}

// License verification.
type License struct {
	Success  bool `json:"success"`
	Uses     int  `json:"uses"`
	Purchase struct {
		ID                      string    `json:"id"`
		CreatedAt               time.Time `json:"created_at"`
		CustomFields            []string  `json:"custom_fields"`
		Email                   string    `json:"email"`
		FullName                string    `json:"full_name"`
		ProductName             string    `json:"product_name"`
		SubscriptionCancelledAt time.Time `json:"subscription_cancelled_at"`
		SubscriptionFailedAt    time.Time `json:"subscription_failed_at"`
		Variants                string    `json:"variants"`
		Chargebacked            bool      `json:"chargebacked"`
		Refunded                bool      `json:"refunded"`
	} `json:"purchase"`
}

// Cancelled returns true if the subscription has been cancelled.
func (l *License) Cancelled() bool {
	return !l.Purchase.SubscriptionCancelledAt.IsZero()
}

// Failed returns true if the subscription has failed.
func (l *License) Failed() bool {
	return !l.Purchase.SubscriptionFailedAt.IsZero()
}

// Verify a product license key.
func (l *Licenses) Verify(product, key string) (v *License, err error) {
	body := url.Values{}
	body.Set("product_permalink", product)
	body.Set("license_key", key)

	res, err := http.Post("https://api.gumroad.com/v2/licenses/verify", formMime, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		err := Error{
			Status: res.StatusCode,
		}

		if err := json.NewDecoder(res.Body).Decode(&err); err != nil {
			return nil, errors.Wrap(err, "decoding error")
		}

		return nil, err
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrap(err, "decoding")
	}

	return
}
