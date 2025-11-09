package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zuhairm2001/rp-coupon-generator/internal/models"
	"github.com/zuhairm2001/rp-coupon-generator/internal/utils"
	"github.com/zuhairm2001/rp-coupon-generator/internal/woocommerce"
)

var formTmpl = template.Must(template.ParseFiles("templates/form.html"))
var resultTmpl = template.Must(template.ParseFiles("templates/result.html"))

// Handler to serve the form
func FormHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Amount        int
		MinimumAmount float64
		CouponCode    string
		Expires       bool
		ExpiryDate    string
		Today         string
		UsageLimit    int
	}{
		Amount:        10,
		MinimumAmount: 100.0,
		CouponCode:    utils.GenerateCode(8),
		Expires:       false,
		ExpiryDate:    time.Now().Format("2006-01-02"),
		Today:         time.Now().Format("2006-01-02"),
		UsageLimit:    0,
	}

	if err := formTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler to process form submission
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	envData := getEnv(w)
	formData := parseFormData(w, r)

	client := woocommerce.NewClient(envData.BaseURL, envData.APIKey, envData.APISecret)
	ctx := context.Background()
	coupon, err := client.CreateCoupon(ctx, formData)
	if err != nil {
		log.Printf("Error creating coupon: %v", err)
		renderError(w, "Failed to generate coupon: "+err.Error())
		return
	}

	// Render success template
	data := models.ResultData{
		Success:       true,
		CouponCode:    coupon.Code,
		CouponDetails: *coupon,
	}

	if err := resultTmpl.Execute(w, data); err != nil {
		log.Printf("Error executing result template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// renderError renders the error template with the given error message
func renderError(w http.ResponseWriter, errorMsg string) {
	data := models.ResultData{
		Success:      false,
		ErrorMessage: errorMsg,
	}

	if err := resultTmpl.Execute(w, data); err != nil {
		log.Printf("Error executing error template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getEnv(w http.ResponseWriter) models.EnvData {

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		renderError(w, "BASE_URL environment variable not set")
	}
	consumerKey := os.Getenv("WOOCOMMERCE_API_KEY")
	if consumerKey == "" {
		renderError(w, "WOOCOMMERCE_API_KEY environment variable not set")
	}
	consumerSecret := os.Getenv("WOOCOMMERCE_API_SECRET")
	if consumerSecret == "" {
		renderError(w, "WOOCOMMERCE_API_SECRET environment variable not set")
	}

	return models.EnvData{
		BaseURL:   baseURL,
		APIKey:    consumerKey,
		APISecret: consumerSecret,
	}

}

func parseFormData(w http.ResponseWriter, r *http.Request) models.FormData {

	amount := r.FormValue("amount")
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		renderError(w, "Invalid amount - please enter a valid number")
	}

	minimumAmount := r.FormValue("minimum_amount")
	minimumAmountFloat, err := strconv.ParseFloat(minimumAmount, 64)
	if err != nil {
		renderError(w, "Invalid minimum amount - please enter a valid number")
	}

	discountType := r.FormValue("discount_type")

	couponCode := r.FormValue("coupon_code")
	if couponCode == "" {
		renderError(w, "Coupon code is required")
	}

	expiryOption := r.FormValue("expiry_option")
	expiryDate := parseExpiryDate(expiryOption, w, r)

	usageLimitStr := r.FormValue("usage_limit")
	usageLimit := parseUsageLimit(usageLimitStr, w)

	log.Printf("Received update: Amount=%s, Min=%s, Expiry=%s, UsageLimit=%d", amount, minimumAmount, expiryDate, usageLimit)

	data := models.FormData{
		Amount:        amountInt,
		MinimumAmount: minimumAmountFloat,
		DiscountType:  discountType,
		CouponCode:    couponCode,
		ExpiryDate:    expiryDate,
		UsageLimit:    usageLimit,
	}
	return data

}

func parseExpiryDate(expiryOption string, w http.ResponseWriter, r *http.Request) string {

	var expiryDate string
	if expiryOption == "custom_date" {
		expiryDate = r.FormValue("expiry_date")
		if expiryDate == "" {
			renderError(w, "Expiry date is required when custom expiry is selected")
		}

		// Validate that expiry date is after today
		parsedExpiryDate, err := time.Parse("2006-01-02", expiryDate)
		if err != nil {
			renderError(w, "Invalid expiry date format")
		}

		today := time.Now().Truncate(24 * time.Hour)
		if parsedExpiryDate.Before(today) || parsedExpiryDate.Equal(today) {
			renderError(w, "Expiry date must be after today")
		}
	} else {
		expiryDate = "" // Empty string indicates no expiry
	}
	return expiryDate
}

func parseUsageLimit(usageLimitStr string, w http.ResponseWriter) int {
	if usageLimitStr != "" {
		usageLimit, err := strconv.Atoi(usageLimitStr)
		if err != nil {
			renderError(w, "Invalid usage limit - please enter a valid number")
		}
		if usageLimit < 0 {
			renderError(w, "Usage limit cannot be negative")
		}
		return usageLimit
	}
	return 0
}
