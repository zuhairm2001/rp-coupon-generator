package woocommerce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zuhairm2001/rp-coupon-generator/internal/models"
)

func (c *Client) CreateCoupon(ctx context.Context, formData models.FormData) (*models.CouponResponse, error) {

	couponReq := models.CouponRequest{
		Code:             formData.CouponCode,
		DiscountType:     models.DiscountType(formData.DiscountType),
		UsageLimit:       1,
		Amount:           formData.Amount,
		IndividualUse:    true,
		ExcludeSaleItems: true,
		MinimumAmount:    fmt.Sprintf("%.2f", formData.MinimumAmount),
	}

	jsonData, err := json.Marshal(couponReq)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal coupon request: %w", err)
	}

	url := fmt.Sprintf("%swp-json/wc/v3/coupons", c.BaseURL)

	log.Printf("Creating coupon with URL: %s", url)
	log.Printf("Coupon request: %s", string(jsonData))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(c.ConsumerKey, c.ConsumerSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	log.Printf("Coupon response status: %s", resp.Status)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create coupon: %s", resp.Status)
	}

	var CouponResponse models.CouponResponse
	err = json.NewDecoder(resp.Body).Decode(&CouponResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Printf("Coupon response: %v", CouponResponse)

	return &CouponResponse, nil
}
