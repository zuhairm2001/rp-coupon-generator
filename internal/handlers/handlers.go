package handlers

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/zuhairm2001/rp-coupon-generator/internal/models"
	"github.com/zuhairm2001/rp-coupon-generator/internal/woocommerce"
)

var formTmpl = template.Must(template.ParseFiles("templates/form.html"))
var resultTmpl = template.Must(template.ParseFiles("templates/result.html"))

// Handler to serve the form
func FormHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Amount        int
		MinimumAmount float64
	}{
		Amount:        10,
		MinimumAmount: 100.0,
	}

	if err := formTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ResultData represents the data passed to the result template
type ResultData struct {
	Success       bool
	CouponCode    string
	CouponDetails *models.CouponResponse
	ErrorMessage  string
}

// Handler to process form submission
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		renderError(w, "BASE_URL environment variable not set")
		return
	}
	consumerKey := os.Getenv("WOOCOMMERCE_API_KEY")
	if consumerKey == "" {
		renderError(w, "WOOCOMMERCE_API_KEY environment variable not set")
		return
	}
	consumerSecret := os.Getenv("WOOCOMMERCE_API_SECRET")
	if consumerSecret == "" {
		renderError(w, "WOOCOMMERCE_API_SECRET environment variable not set")
		return
	}

	if err := r.ParseForm(); err != nil {
		renderError(w, "Failed to parse form data")
		return
	}

	amount := r.FormValue("amount")
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		renderError(w, "Invalid amount - please enter a valid number")
		return
	}

	minimumAmount := r.FormValue("minimum_amount")
	minimumAmountFloat, err := strconv.ParseFloat(minimumAmount, 64)
	if err != nil {
		renderError(w, "Invalid minimum amount - please enter a valid number")
		return
	}

	log.Printf("Received update: Amount=%s, Min=%s", amount, minimumAmount)

	client := woocommerce.NewClient(baseURL, consumerKey, consumerSecret)
	ctx := context.Background()
	coupon, err := client.CreateCoupon(ctx, amountInt, minimumAmountFloat)
	if err != nil {
		log.Printf("Error creating coupon: %v", err)
		renderError(w, "Failed to generate coupon: "+err.Error())
		return
	}

	// Render success template
	data := ResultData{
		Success:       true,
		CouponCode:    coupon.Code,
		CouponDetails: coupon,
	}

	if err := resultTmpl.Execute(w, data); err != nil {
		log.Printf("Error executing result template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// renderError renders the error template with the given error message
func renderError(w http.ResponseWriter, errorMsg string) {
	data := ResultData{
		Success:      false,
		ErrorMessage: errorMsg,
	}

	if err := resultTmpl.Execute(w, data); err != nil {
		log.Printf("Error executing error template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
