package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/zuhairm2001/rp-coupon-generator/internal/utils"
)

var loginTmpl = template.Must(template.ParseFiles("templates/login.html"))

type LoginData struct {
	ErrorMessage string
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		LoginSubmitHandler(w, r)
		return
	}

	data := LoginData{
		ErrorMessage: "",
	}

	if err := loginTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Login handles the login request
func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	password := r.FormValue("password")
	if !utils.Login(password) {
		data := LoginData{
			ErrorMessage: "Invalid password",
		}
		if err := loginTmpl.Execute(w, data); err != nil {
			log.Printf("Error executing login template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		renderError(w, "Invalid password")
		return
	}

	FormHandler(w, r)

}
