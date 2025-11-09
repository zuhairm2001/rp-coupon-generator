package main

import (
	"log"
	"net/http"

	"github.com/zuhairm2001/rp-coupon-generator/internal/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.LoginPageHandler)
	http.HandleFunc("/submit", handlers.SubmitHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
