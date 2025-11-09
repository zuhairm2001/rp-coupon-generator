package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

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
	}{
		Amount:        10,
		MinimumAmount: 100.0,
		CouponCode:    utils.GenerateCode(8),
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

	//no error checking for discount type as it is a drop down
	discountType := r.FormValue("discount_type")

	couponCode := r.FormValue("coupon_code")
	if couponCode == "" {
		renderError(w, "Coupon code is required")
	}

	log.Printf("Received update: Amount=%s, Min=%s", amount, minimumAmount)

	data := models.FormData{
		Amount:        amountInt,
		MinimumAmount: minimumAmountFloat,
		DiscountType:  discountType,
		CouponCode:    couponCode,
	}
	return data

}
