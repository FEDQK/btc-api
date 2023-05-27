package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Subscriber struct {
	Email string `json:"email"`
}

type CurrencyResponse struct {
	Bpi struct {
		UAH struct {
			Rate string `json:"rate"`
		} `json:"UAH"`
	} `json:"bpi"`
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

func main() {
	http.HandleFunc("/btc-to-uah", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice/BTC.json")
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

		fmt.Fprintf(w, "BTC to UAH rate: %s", currencyResponse.Bpi.UAH.Rate)
	})

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

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice/BTC.json")
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

		for _, subscriber := range subscribers {
			fmt.Printf("Sending BTC to UAH rate (%s) to %s\n", currencyResponse.Bpi.UAH.Rate, subscriber.Email)
		}

		fmt.Fprintf(w, "Notifications sent")
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
