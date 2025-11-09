package utils

import (
	"log"
	"os"
)

func Login(password string) bool {
	if os.Getenv("PASSWORD") == "" {
		log.Default().Print("No password set")
		return false
	}

	if password == os.Getenv("PASSWORD") {
		log.Default().Println("Login successful")
		return true
	}
	log.Default().Println("Login failed")
	return false
}
