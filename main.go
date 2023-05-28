package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
)

type Subscriber struct {
	Email string `json:"email"`
}

type CurrencyResponse struct {
	Bitcoin struct {
		UAH float64 `json:"uah"`
	} `json:"bitcoin"`
}

var subscribers []Subscriber

func init() {
	data, err := ioutil.ReadFile("subscribers.json")
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		json.Unmarshal(data, &subscribers)
	}
}

func sendEmail(to, subject, body string) error {
	// Налаштування підключення до SMTP-сервера
	smtpHost := "smtp.elasticemail.com"    // Адреса SMTP-сервера
	smtpPort := 2525                    // Порт SMTP-сервера
	smtpUsername := "vovanchikd1996@gmail.com"    // Ім'я користувача для аутентифікації
	smtpPassword := "D313F79FA8B78B33E229D16B54D5FC618DDD"    // Пароль для аутентифікації

	// Створення аутентифікаційних даних
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Формування повідомлення
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Відправка повідомлення
	err := smtp.SendMail(smtpHost+":"+fmt.Sprintf("%d", smtpPort), auth, smtpUsername, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}


func main() {
	http.HandleFunc("/rate", rate.get)

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
	
		email := r.FormValue("email")
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}
	
		for _, subscriber := range subscribers {
			if subscriber.Email == email {
				http.Error(w, "Already subscribed", http.StatusBadRequest)
				return
			}
		}
	
		subscriber := Subscriber{Email: email}
		subscribers = append(subscribers, subscriber)
	
		data, err := json.Marshal(subscribers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		err = ioutil.WriteFile("subscribers.json", data, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		fmt.Fprintf(w, "Subscribed successfully")
	})

	http.HandleFunc("/sendEmails", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
	
		resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=UAH")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
	
		var currencyResponse CurrencyResponse
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&currencyResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		data, err := ioutil.ReadFile("subscribers.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		err = json.Unmarshal(data, &subscribers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		rate := currencyResponse.Bitcoin.UAH
		for _, subscriber := range subscribers {
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
	})
	

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
