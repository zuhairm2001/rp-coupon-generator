package woocommerce

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL        string
	ConsumerKey    string
	ConsumerSecret string
	HTTPClient     *http.Client
}

func NewClient(baseURL, consumerKey, consumerSecret string) *Client {
	return &Client{
		BaseURL:        baseURL,
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
