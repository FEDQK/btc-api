package sendEmails

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/FEDQK/btc-api/constants"
	"github.com/FEDQK/btc-api/models"
)

var (
	smtpHost string
	smtpPort string
	smtpUsername string
	smtpPassword string
)

func init() {
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	smtpUsername = os.Getenv("SMTP_USERNAME")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" {
		log.Fatal("SMTP configuration values are missing")
	}
}

func sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	err := smtp.SendMail(smtpHost + ":" + smtpPort, auth, smtpUsername, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

func Post(w http.ResponseWriter, r *http.Request, subscribers *[]models.Subscriber) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	resp, err := http.Get(constants.PRICE_API_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var currencyResponse models.CurrencyResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&currencyResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rate := currencyResponse.Bitcoin.UAH
	for _, subscriber := range *subscribers {
		fmt.Printf("Sending BTC to UAH rate (%.2f) to %s\n", rate, subscriber.Email)
		subject := "BTC to UAH Rate Update"
		body := fmt.Sprintf("The current BTC to UAH rate is: %.2f", rate)
		err := sendEmail(subscriber.Email, subject, body)
		if err != nil {
			fmt.Printf("Failed to send email to %s: %v\n", subscriber.Email, err)
			continue
		}
	}

	fmt.Fprintf(w, "Notifications sent")
}