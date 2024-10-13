package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// GenerateOTP generates a 6-digit random OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())  // Seed the random number generator
	otp := rand.Intn(900000) + 100000 // Ensure a 6-digit number
	return fmt.Sprintf("%06d", otp)   // Format as 6-digit string
}

// SendOTP sends the OTP to the user's phone number via an SMS API
func SendSMSOTP(phoneNumber string, otp string) error {
	// Replace with actual SMS API details (e.g., Twilio, Razorpay, or other service)
	apiURL := "https://sms-provider.com/send"
	message := fmt.Sprintf("Your OTP is: %s", otp)

	// Request payload (adjust according to your SMS API)
	payload := fmt.Sprintf("phone=%s&message=%s", phoneNumber, message)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for successful response (status code 200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send OTP, status code: %d", resp.StatusCode)
	}

	return nil
}
