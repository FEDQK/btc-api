package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/FEDQK/btc-api/api/rate"
	"github.com/FEDQK/btc-api/api/sendEmails"
	"github.com/FEDQK/btc-api/api/subscribe"
	"github.com/FEDQK/btc-api/constants"
	"github.com/FEDQK/btc-api/models"
)

var subscribers []models.Subscriber

func init() {
	data, err := ioutil.ReadFile(constants.SUBSCRIBERS_FILE_NAME)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		json.Unmarshal(data, &subscribers)
	}
}

func main() {
	http.HandleFunc("/rate", rate.Get)

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		subscribe.Post(w, r, &subscribers)
	})

	http.HandleFunc("/sendEmails", func(w http.ResponseWriter, r *http.Request) {
		sendEmails.Post(w, r, &subscribers)
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
