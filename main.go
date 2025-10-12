package main

//html template that will take in request with information on what coupon to generate

//there needs to be a request handler

//we will need to make a api request to woocommerce to generate the coupon

//on generation finish, we will need to respond to the original request with the coupon name

import (
	"log"
	"net/http"

	"github.com/zuhairm2001/rp-coupon-generator/internal/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.FormHandler)
	http.HandleFunc("/submit", handlers.SubmitHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
