package subscribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/FEDQK/btc-api/constants"
	"github.com/FEDQK/btc-api/models"
)

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isEmailExist(email string, subscribers *[]models.Subscriber) bool {
	for _, subscriber := range *subscribers {
		if subscriber.Email == email {
			return true
		}
	}
	return false
}

func subscribe(email string, subscribers *[]models.Subscriber) error {
	subscriber := models.Subscriber{Email: email}
	*subscribers = append(*subscribers, subscriber)

	data, err := json.Marshal(subscribers)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(constants.SUBSCRIBERS_FILE_NAME, data, 0644)
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

	email := r.FormValue("email")
	if !isEmailValid(email) {
		http.Error(w, "Email is not valid", http.StatusBadRequest)
		return
	}

	if isEmailExist(email, subscribers) {
		http.Error(w, "Email already subscribed", http.StatusBadRequest)
		return
	}

	subscribeStatus := subscribe(email, subscribers)

	if subscribeStatus != nil {
		http.Error(w, subscribeStatus.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Subscribed successfully")
}